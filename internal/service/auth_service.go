// Package service contains the business logic for authentication, session
// management, and the workflow slices. It depends on repository interfaces,
// not concrete database adapters, so it can be tested without a real
// database.
package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/auth"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/observability"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/repository"
	"github.com/google/uuid"
)

// Errors are deliberately neutral — they don't reveal whether an account
// exists.
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrOTPInvalid         = errors.New("verification code is invalid or expired")
	ErrSessionInvalid     = errors.New("session is invalid or expired")
	ErrAccountInactive    = errors.New("account is inactive")
	ErrInvalidRolePreview = errors.New("invalid role preview")
)

// AuthService handles login, email verification, logout, session validation,
// and role preview. It coordinates the user, session, email-challenge, and
// audit repositories.
type AuthService struct {
	users      repository.UserRepository
	sessions   repository.SessionRepository
	challenges repository.EmailChallengeRepository
	audit      repository.AuditEventRepository
	otpCfg     auth.OTPConfig
	sessionTTL time.Duration
	idleTTL    time.Duration
	devCode    bool
	log        *slog.Logger
}

// NewAuthService constructs the auth service.
func NewAuthService(
	users repository.UserRepository,
	sessions repository.SessionRepository,
	challenges repository.EmailChallengeRepository,
	audit repository.AuditEventRepository,
	otpCfg auth.OTPConfig,
	sessionTTL, idleTTL time.Duration,
	devCodeVisible bool,
	log *slog.Logger,
) *AuthService {
	return &AuthService{
		users:      users,
		sessions:   sessions,
		challenges: challenges,
		audit:      audit,
		otpCfg:     otpCfg,
		sessionTTL: sessionTTL,
		idleTTL:    idleTTL,
		devCode:    devCodeVisible,
		log:        log,
	}
}

// LoginResult holds the outcome of a login attempt.
type LoginResult struct {
	SessionToken string
	NeedsOTP     bool
	DisplayName  string
}

// Login validates credentials. If the user exists and the password matches,
// it creates a session and returns the session token. If email is not yet
// verified, it returns NeedsOTP=true instead of a session token. On any
// failure (user not found, wrong password, inactive), it returns the same
// neutral error so the caller cannot distinguish the cause.
func (s *AuthService) Login(ctx context.Context, identifier, password, ip, userAgent string) (*LoginResult, error) {
	corrID := observability.CorrelationFrom(ctx)

	user, err := s.users.GetByUsernameOrEmail(ctx, identifier)
	if err != nil {
		// User not found — verify against a dummy hash to equalize timing,
		// then record a failed attempt with a neutral message.
		_, _ = auth.VerifyPassword(password, auth.DummyHash())
		s.auditLogin(ctx, nil, domain.AuditActionLoginFailed, corrID, ip, userAgent)
		return nil, ErrInvalidCredentials
	}

	ok, err := auth.VerifyPassword(password, user.PasswordHash)
	if err != nil || !ok {
		s.auditLogin(ctx, &user.ID, domain.AuditActionLoginFailed, corrID, ip, userAgent)
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		s.auditLogin(ctx, &user.ID, domain.AuditActionLoginFailed, corrID, ip, userAgent)
		return nil, ErrInvalidCredentials
	}

	if !user.EmailVerified {
		// Generate a new OTP challenge for email verification.
		code, _, err := s.createEmailChallenge(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("create email challenge: %w", err)
		}
		if s.devCode {
			s.log.InfoContext(ctx, "dev otp code generated (dev mode only)",
				"user_id", string(user.ID),
				"otp_code", code,
			)
		}
		return &LoginResult{NeedsOTP: true, DisplayName: user.DisplayName}, nil
	}

	// Email verified — create a session.
	token, hash, err := auth.GenerateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("generate session: %w", err)
	}
	now := time.Now().UTC()
	session := &domain.Session{
		ID:          domain.SessionID(uuid.NewString()),
		UserID:      user.ID,
		TokenHash:   hash,
		RolePreview: domain.RolePreviewNone,
		IPAddress:   ip,
		UserAgent:   userAgent,
		CreatedAt:   now,
		LastSeenAt:  now,
		ExpiresAt:   now.Add(s.sessionTTL),
	}
	if err := s.sessions.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	s.auditLogin(ctx, &user.ID, domain.AuditActionLoginSuccess, corrID, ip, userAgent)

	return &LoginResult{
		SessionToken: token,
		DisplayName:  user.DisplayName,
	}, nil
}

// VerifyEmail consumes an OTP code and marks the user's email as verified.
// If successful, it also creates a full session.
func (s *AuthService) VerifyEmail(ctx context.Context, identifier, code, ip, userAgent string) (*LoginResult, error) {
	corrID := observability.CorrelationFrom(ctx)

	user, err := s.users.GetByUsernameOrEmail(ctx, identifier)
	if err != nil {
		s.auditLog(ctx, nil, domain.AuditActionEmailVerifyFailed, "email_challenge", "", corrID)
		return nil, ErrOTPInvalid
	}

	challenge, err := s.challenges.GetLatestByUserID(ctx, user.ID, "email_verification")
	if err != nil || challenge == nil {
		s.auditLog(ctx, &user.ID, domain.AuditActionEmailVerifyFailed, "email_challenge", "", corrID)
		return nil, ErrOTPInvalid
	}

	now := time.Now().UTC()
	if challenge.IsExpired(now) || challenge.IsConsumed() {
		s.auditLog(ctx, &user.ID, domain.AuditActionEmailVerifyFailed, "email_challenge", string(challenge.ID), corrID)
		return nil, ErrOTPInvalid
	}

	if !auth.VerifyOTPHash(code, challenge.CodeHash) {
		s.auditLog(ctx, &user.ID, domain.AuditActionEmailVerifyFailed, "email_challenge", string(challenge.ID), corrID)
		return nil, ErrOTPInvalid
	}

	// Consume the challenge.
	if err := s.challenges.MarkConsumed(ctx, challenge.ID, now); err != nil {
		return nil, fmt.Errorf("consume challenge: %w", err)
	}

	// Mark email verified.
	if err := s.users.UpdateEmailVerified(ctx, user.ID, true); err != nil {
		return nil, fmt.Errorf("update email verified: %w", err)
	}

	// Create session.
	token, hash, err := auth.GenerateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("generate session: %w", err)
	}
	session := &domain.Session{
		ID:          domain.SessionID(uuid.NewString()),
		UserID:      user.ID,
		TokenHash:   hash,
		RolePreview: domain.RolePreviewNone,
		IPAddress:   ip,
		UserAgent:   userAgent,
		CreatedAt:   now,
		LastSeenAt:  now,
		ExpiresAt:   now.Add(s.sessionTTL),
	}
	if err := s.sessions.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	s.auditLog(ctx, &user.ID, domain.AuditActionEmailVerified, "user", string(user.ID), corrID)

	return &LoginResult{
		SessionToken: token,
		DisplayName:  user.DisplayName,
	}, nil
}

// Logout revokes the session identified by token hash and records an audit
// event.
func (s *AuthService) Logout(ctx context.Context, tokenHash string) error {
	sess, err := s.sessions.GetByTokenHash(ctx, tokenHash)
	if err != nil || sess == nil {
		return ErrSessionInvalid
	}
	now := time.Now().UTC()
	if err := s.sessions.Revoke(ctx, sess.ID, now); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	s.auditLog(ctx, &sess.UserID, domain.AuditActionLogout, "session", string(sess.ID), observability.CorrelationFrom(ctx))
	return nil
}

// ValidateSession checks the session by token hash, returning the session
// and user if valid. It updates last_seen_at on success.
func (s *AuthService) ValidateSession(ctx context.Context, tokenHash string) (*domain.Session, *domain.User, error) {
	sess, err := s.sessions.GetByTokenHash(ctx, tokenHash)
	if err != nil || sess == nil {
		return nil, nil, ErrSessionInvalid
	}
	now := time.Now().UTC()
	if !sess.IsValid(now, s.idleTTL) {
		return nil, nil, ErrSessionInvalid
	}
	user, err := s.users.GetByID(ctx, sess.UserID)
	if err != nil || user == nil || !user.IsActive {
		return nil, nil, ErrSessionInvalid
	}
	if err := s.sessions.UpdateLastSeen(ctx, sess.ID, now); err != nil {
		s.log.WarnContext(ctx, "failed to update last_seen", "error", err)
	}
	return sess, user, nil
}

// SetRolePreview updates the role preview on the session. Per the spec, this
// changes the visible navigation and allowed workflow actions but NOT the
// authenticated identity, audit actor, or server-side access decision.
func (s *AuthService) SetRolePreview(ctx context.Context, tokenHash string, preview domain.RolePreview) error {
	if !domain.IsValidRolePreview(preview) {
		return ErrInvalidRolePreview
	}
	sess, err := s.sessions.GetByTokenHash(ctx, tokenHash)
	if err != nil || sess == nil {
		return ErrSessionInvalid
	}
	now := time.Now().UTC()
	if !sess.IsValid(now, s.idleTTL) {
		return ErrSessionInvalid
	}
	if err := s.sessions.UpdateRolePreview(ctx, sess.ID, preview); err != nil {
		return fmt.Errorf("update role preview: %w", err)
	}
	action := domain.AuditActionRolePreviewSet
	if preview == domain.RolePreviewNone {
		action = domain.AuditActionRolePreviewCleared
	}
	s.auditLog(ctx, &sess.UserID, action, "session", string(sess.ID), observability.CorrelationFrom(ctx))
	return nil
}

// RevokeAllSessions revokes all sessions for a user (e.g., on password change).
func (s *AuthService) RevokeAllSessions(ctx context.Context, userID domain.UserID) error {
	now := time.Now().UTC()
	if err := s.sessions.RevokeAllForUser(ctx, userID, now); err != nil {
		return fmt.Errorf("revoke all sessions: %w", err)
	}
	s.auditLog(ctx, &userID, domain.AuditActionSessionRevoked, "user", string(userID), observability.CorrelationFrom(ctx))
	return nil
}

func (s *AuthService) createEmailChallenge(ctx context.Context, userID domain.UserID) (plaintext, hash string, err error) {
	now := time.Now().UTC()
	code, err := auth.GenerateOTP(s.otpCfg, now)
	if err != nil {
		return "", "", err
	}
	hash = auth.HashOTPCode(code)
	challenge := &domain.EmailChallenge{
		ID:        uuid.NewString(),
		UserID:    userID,
		CodeHash:  hash,
		Purpose:   "email_verification",
		CreatedAt: now,
		ExpiresAt: now.Add(10 * time.Minute),
	}
	if err := s.challenges.Create(ctx, challenge); err != nil {
		return "", "", err
	}
	return code, hash, nil
}

func (s *AuthService) auditLogin(ctx context.Context, userID *domain.UserID, action, corrID, ip, userAgent string) {
	details := map[string]any{
		"ip":         ip,
		"user_agent": userAgent,
	}
	if userID != nil {
		_ = s.audit.Create(ctx, &domain.AuditEvent{
			CorrelationID: corrID,
			ActorUserID:   userID,
			ActorRole:     "administrator",
			Action:        action,
			EntityType:    "session",
			Details:       details,
		})
	} else {
		_ = s.audit.Create(ctx, &domain.AuditEvent{
			CorrelationID: corrID,
			ActorRole:     "unknown",
			Action:        action,
			EntityType:    "session",
			Details:       details,
		})
	}
}

func (s *AuthService) auditLog(ctx context.Context, userID *domain.UserID, action, entityType, entityID, corrID string) {
	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID,
		ActorUserID:   userID,
		ActorRole:     "administrator",
		Action:        action,
		EntityType:    entityType,
		EntityID:      entityID,
		Details:       map[string]any{},
	})
}

// Compile-time assertion that the package compiles.
var _ = time.UTC
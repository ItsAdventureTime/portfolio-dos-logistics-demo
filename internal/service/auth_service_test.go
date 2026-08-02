package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/auth"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/observability"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/store"
	"log/slog"
)

func testOTPConfig() auth.OTPConfig {
	return auth.DefaultOTPConfig([]byte("0123456789abcdef0123456789abcdef"))
}

func newTestService(t *testing.T, devCodeVisible bool) (*AuthService, *store.MemUserStore, *store.MemSessionStore, *store.MemEmailChallengeStore, *store.MemAuditStore) {
	t.Helper()
	users := store.NewMemUserStore()
	sessions := store.NewMemSessionStore()
	challenges := store.NewMemEmailChallengeStore()
	audit := store.NewMemAuditStore()
	svc := NewAuthService(
		users, sessions, challenges, audit,
		testOTPConfig(),
		24*time.Hour, 1*time.Hour,
		devCodeVisible,
		slog.Default(),
	)
	return svc, users, sessions, challenges, audit
}

func seedVerifiedUser(t *testing.T, users *store.MemUserStore, username, email string) *domain.User {
	t.Helper()
	hash, err := auth.HashPassword("Password123!")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	u := &domain.User{
		ID:            domain.UserID("user-" + username),
		Username:      username,
		Email:         email,
		PasswordHash:  hash,
		DisplayName:   "Test User",
		EmailVerified: true,
		IsActive:      true,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	users.Seed(u)
	return u
}

func seedUnverifiedUser(t *testing.T, users *store.MemUserStore, username, email string) *domain.User {
	u := seedVerifiedUser(t, users, username, email)
	u.EmailVerified = false
	return u
}

// Acceptance check: "Invalid credentials produce a neutral, recoverable message."
// Both unknown-user and wrong-password must return the same error.
func TestLogin_UnknownUser_NeutralError(t *testing.T) {
	svc, _, _, _, audit := newTestService(t, false)
	ctx := observability.WithCorrelation(context.Background(), "test-corr")
	_, err := svc.Login(ctx, "nobody@example.com", "anything", "1.2.3.4", "test")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials for unknown user, got %v", err)
	}
	// Audit should record a failed login without an actor (unknown user).
	events := audit.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}
	if events[0].Action != domain.AuditActionLoginFailed {
		t.Errorf("action = %s, want %s", events[0].Action, domain.AuditActionLoginFailed)
	}
}

func TestLogin_WrongPassword_NeutralError(t *testing.T) {
	svc, users, _, _, audit := newTestService(t, false)
	seedVerifiedUser(t, users, "testuser", "test@example.com")
	ctx := observability.WithCorrelation(context.Background(), "test-corr")
	_, err := svc.Login(ctx, "test@example.com", "wrong-password", "1.2.3.4", "test")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials for wrong password, got %v", err)
	}
	// Both unknown-user and wrong-password return the same error: neutral.
	events := audit.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}
	if events[0].Action != domain.AuditActionLoginFailed {
		t.Errorf("action = %s, want %s", events[0].Action, domain.AuditActionLoginFailed)
	}
}

// Acceptance check: "Administrator can sign in with password and email verification code."
func TestLogin_VerifiedUser_CreatesSession(t *testing.T) {
	svc, users, sessions, _, audit := newTestService(t, false)
	seedVerifiedUser(t, users, "admin", "admin@example.com")
	ctx := observability.WithCorrelation(context.Background(), "test-corr")
	result, err := svc.Login(ctx, "admin", "Password123!", "1.2.3.4", "test")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if result.NeedsOTP {
		t.Error("verified user should not need OTP")
	}
	if result.SessionToken == "" {
		t.Error("session token should be returned")
	}
	// Session should be in the store.
	hash := auth.HashToken(result.SessionToken)
	sess, err := sessions.GetByTokenHash(ctx, hash)
	if err != nil || sess == nil {
		t.Errorf("session not found in store: %v", err)
	}
	// Audit should record success.
	events := audit.Events()
	found := false
	for _, e := range events {
		if e.Action == domain.AuditActionLoginSuccess {
			found = true
		}
	}
	if !found {
		t.Error("audit log missing login_success event")
	}
}

// Acceptance check: unverified user gets OTP challenge.
func TestLogin_UnverifiedUser_NeedsOTP(t *testing.T) {
	svc, users, _, challenges, _ := newTestService(t, false)
	seedUnverifiedUser(t, users, "unverified", "unverified@example.com")
	ctx := context.Background()
	result, err := svc.Login(ctx, "unverified", "Password123!", "1.2.3.4", "test")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if !result.NeedsOTP {
		t.Error("unverified user should need OTP")
	}
	// A challenge should have been created.
	ch, err := challenges.GetLatestByUserID(ctx, "user-unverified", "email_verification")
	if err != nil || ch == nil {
		t.Errorf("email challenge not created: %v", err)
	}
}

// Acceptance check: email verification with valid code succeeds and creates session.
func TestVerifyEmail_ValidCode_CreatesSession(t *testing.T) {
	svc, users, sessions, _, _ := newTestService(t, true)
	seedUnverifiedUser(t, users, "verifyme", "verify@example.com")
	ctx := context.Background()
	// Login to generate the OTP challenge.
	_, err := svc.Login(ctx, "verifyme", "Password123!", "1.2.3.4", "test")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	// With dev code visible + deterministic TOTP, regenerate the code.
	code, err := auth.GenerateOTP(testOTPConfig(), time.Now().UTC())
	if err != nil {
		t.Fatalf("GenerateOTP: %v", err)
	}
	result, err := svc.VerifyEmail(ctx, "verifyme", code, "1.2.3.4", "test")
	if err != nil {
		t.Fatalf("VerifyEmail: %v", err)
	}
	if result.SessionToken == "" {
		t.Error("session token should be returned after email verification")
	}
	hash := auth.HashToken(result.SessionToken)
	sess, err := sessions.GetByTokenHash(ctx, hash)
	if err != nil || sess == nil {
		t.Errorf("session not found in store: %v", err)
	}
}

// Acceptance check: invalid OTP code is rejected.
func TestVerifyEmail_InvalidCode_Rejected(t *testing.T) {
	svc, users, _, _, _ := newTestService(t, false)
	seedUnverifiedUser(t, users, "rejectme", "reject@example.com")
	ctx := context.Background()
	_, _ = svc.Login(ctx, "reject", "Password123!", "1.2.3.4", "test")
	_, err := svc.VerifyEmail(ctx, "reject", "999999", "1.2.3.4", "test")
	if !errors.Is(err, ErrOTPInvalid) {
		t.Errorf("expected ErrOTPInvalid, got %v", err)
	}
}

// Acceptance check: OTP replay (consumed code) is rejected.
func TestVerifyEmail_ConsumedCode_Rejected(t *testing.T) {
	// This is covered by the challenge.ConsumedAt check. A second attempt
	// with the same code should fail because the challenge is consumed.
	svc, users, _, challenges, _ := newTestService(t, true)
	seedUnverifiedUser(t, users, "replay", "replay@example.com")
	ctx := context.Background()
	_, _ = svc.Login(ctx, "replay", "Password123!", "1.2.3.4", "test")
	code, _ := auth.GenerateOTP(testOTPConfig(), time.Now().UTC())
	_, err := svc.VerifyEmail(ctx, "replay", code, "1.2.3.4", "test")
	if err != nil {
		t.Fatalf("first VerifyEmail: %v", err)
	}
	// Second attempt should fail.
	_, err = svc.VerifyEmail(ctx, "replay", code, "1.2.3.4", "test")
	if !errors.Is(err, ErrOTPInvalid) {
		t.Errorf("expected ErrOTPInvalid on replay, got %v", err)
	}
	_ = challenges
}

// Acceptance check: "Logout and session expiration remove protected access."
func TestLogout_RevokesSession(t *testing.T) {
	svc, users, sessions, _, _ := newTestService(t, false)
	seedVerifiedUser(t, users, "logout", "logout@example.com")
	ctx := context.Background()
	result, _ := svc.Login(ctx, "logout", "Password123!", "1.2.3.4", "test")
	tokenHash := auth.HashToken(result.SessionToken)
	// Session should be valid.
	sess, _, err := svc.ValidateSession(ctx, tokenHash)
	if err != nil || sess == nil {
		t.Fatalf("ValidateSession before logout: %v", err)
	}
	// Logout.
	if err := svc.Logout(ctx, tokenHash); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	// Session should now be invalid.
	_, _, err = svc.ValidateSession(ctx, tokenHash)
	if !errors.Is(err, ErrSessionInvalid) {
		t.Errorf("expected ErrSessionInvalid after logout, got %v", err)
	}
	_ = sessions
}

// Acceptance check: session expiration.
func TestValidateSession_ExpiredSession_Rejected(t *testing.T) {
	svc, users, sessions, _, _ := newTestService(t, false)
	seedVerifiedUser(t, users, "expired", "expired@example.com")
	// Create a session with a very short TTL.
	svc.sessionTTL = 1 * time.Millisecond
	ctx := context.Background()
	result, _ := svc.Login(ctx, "expired", "Password123!", "1.2.3.4", "test")
	tokenHash := auth.HashToken(result.SessionToken)
	time.Sleep(10 * time.Millisecond)
	_, _, err := svc.ValidateSession(ctx, tokenHash)
	if !errors.Is(err, ErrSessionInvalid) {
		t.Errorf("expected ErrSessionInvalid for expired session, got %v", err)
	}
	_ = sessions
}

// Acceptance check: role preview does not change the audit actor.
func TestSetRolePreview_DoesNotChangeAuditActor(t *testing.T) {
	svc, users, _, _, audit := newTestService(t, false)
	seedVerifiedUser(t, users, "preview", "preview@example.com")
	ctx := observability.WithCorrelation(context.Background(), "preview-corr")
	result, _ := svc.Login(ctx, "preview", "Password123!", "1.2.3.4", "test")
	tokenHash := auth.HashToken(result.SessionToken)
	// Set role preview.
	if err := svc.SetRolePreview(ctx, tokenHash, domain.RolePreviewLogisticsCoordinator); err != nil {
		t.Fatalf("SetRolePreview: %v", err)
	}
	// The audit actor should still be the same user, not "logistics_coordinator".
	events := audit.Events()
	var previewEvent *domain.AuditEvent
	for _, e := range events {
		if e.Action == domain.AuditActionRolePreviewSet {
			previewEvent = e
			break
		}
	}
	if previewEvent == nil {
		t.Fatal("missing role_preview_set audit event")
	}
	if previewEvent.ActorRole != "administrator" {
		t.Errorf("actor_role = %s, want administrator (preview must not change audit actor)",
			previewEvent.ActorRole)
	}
}

// Acceptance check: invalid role preview is rejected.
func TestSetRolePreview_InvalidRole_Rejected(t *testing.T) {
	svc, users, _, _, _ := newTestService(t, false)
	seedVerifiedUser(t, users, "invalid", "invalid@example.com")
	ctx := context.Background()
	result, _ := svc.Login(ctx, "invalid", "Password123!", "1.2.3.4", "test")
	tokenHash := auth.HashToken(result.SessionToken)
	err := svc.SetRolePreview(ctx, tokenHash, domain.RolePreview("superadmin"))
	if !errors.Is(err, ErrInvalidRolePreview) {
		t.Errorf("expected ErrInvalidRolePreview, got %v", err)
	}
}

// Acceptance check: inactive account cannot login.
func TestLogin_InactiveAccount_Rejected(t *testing.T) {
	svc, users, _, _, _ := newTestService(t, false)
	u := seedVerifiedUser(t, users, "inactive", "inactive@example.com")
	u.IsActive = false
	ctx := context.Background()
	_, err := svc.Login(ctx, "inactive", "Password123!", "1.2.3.4", "test")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials for inactive account, got %v", err)
	}
}
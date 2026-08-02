// Package domain defines the core business types for identity, audit, and
// workflow records. These types are stack-agnostic. They depend on no
// transport, database driver, or framework.
package domain

import (
	"time"
)

// UserID uniquely identifies a user.
type UserID string

// SessionID uniquely identifies a session.
type SessionID string

// RolePreview represents the operating view the Administrator may preview.
// Role preview changes the visible navigation and allowed workflow actions.
// It does not change the authenticated identity, the audit actor, or the
// server-side access decision.
type RolePreview string

const (
	// RolePreviewNone is the default. The Administrator sees the full set.
	RolePreviewNone RolePreview = ""
	// RolePreviewLogisticsCoordinator simulates the logistics coordinator view.
	RolePreviewLogisticsCoordinator RolePreview = "logistics_coordinator"
	// RolePreviewFreightOpsApprover simulates the freight operations approver view.
	RolePreviewFreightOpsApprover RolePreview = "freight_ops_approver"
	// RolePreviewDisbursementController simulates the disbursement controller view.
	RolePreviewDisbursementController RolePreview = "disbursement_controller"
	// RolePreviewFinanceOpsLead simulates the finance operations lead view.
	RolePreviewFinanceOpsLead RolePreview = "finance_ops_lead"
)

// AllRolePreviews returns every valid preview role for validation.
func AllRolePreviews() []RolePreview {
	return []RolePreview{
		RolePreviewNone,
		RolePreviewLogisticsCoordinator,
		RolePreviewFreightOpsApprover,
		RolePreviewDisbursementController,
		RolePreviewFinanceOpsLead,
	}
}

// IsValidRolePreview reports whether r is a recognized preview role.
func IsValidRolePreview(r RolePreview) bool {
	for _, v := range AllRolePreviews() {
		if v == r {
			return true
		}
	}
	return false
}

// User is the single authenticated account type (Administrator).
type User struct {
	ID            UserID
	Username      string
	Email         string
	PasswordHash  string
	DisplayName   string
	EmailVerified bool
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Session is a server-managed session. The plaintext token is never stored;
// only the hash is persisted. Sessions support idle and absolute timeouts,
// rotation, and revocation.
type Session struct {
	ID          SessionID
	UserID      UserID
	TokenHash   string
	RolePreview RolePreview
	IPAddress   string
	UserAgent   string
	CreatedAt   time.Time
	LastSeenAt  time.Time
	ExpiresAt   time.Time
	RevokedAt   *time.Time
}

// IsRevoked reports whether the session has been revoked.
func (s Session) IsRevoked() bool { return s.RevokedAt != nil }

// IsExpired reports whether the session has passed its absolute expiry.
func (s Session) IsExpired(now time.Time) bool { return now.After(s.ExpiresAt) }

// IsIdleExpired reports whether the session has been idle too long.
func (s Session) IsIdleExpired(now time.Time, idleTimeout time.Duration) bool {
	return now.Sub(s.LastSeenAt) > idleTimeout
}

// IsValid reports whether the session is usable (not revoked, not expired,
// not idle-expired).
func (s Session) IsValid(now time.Time, idleTimeout time.Duration) bool {
	return !s.IsRevoked() && !s.IsExpired(now) && !s.IsIdleExpired(now, idleTimeout)
}

// EmailChallenge is a one-time OTP code for email verification.
type EmailChallenge struct {
	ID        string
	UserID    UserID
	CodeHash  string
	Purpose   string
	CreatedAt time.Time
	ExpiresAt time.Time
	ConsumedAt *time.Time
}

// IsExpired reports whether the challenge has expired.
func (c EmailChallenge) IsExpired(now time.Time) bool { return now.After(c.ExpiresAt) }

// IsConsumed reports whether the challenge has already been used.
func (c EmailChallenge) IsConsumed() bool { return c.ConsumedAt != nil }

// AuditEvent records a security-relevant action. Per docs/spec/03, every
// sensitive action must create an audit event identifying actor, action,
// entity, time, and correlation id. Audit events are append-only.
type AuditEvent struct {
	ID            int64
	CorrelationID string
	ActorUserID   *UserID
	ActorRole     string
	Action        string
	EntityType    string
	EntityID      string
	Details       map[string]any
	CreatedAt     time.Time
}

// Audit action constants. These are the names recorded in audit_events.action.
const (
	AuditActionLoginAttempt      = "login_attempt"
	AuditActionLoginSuccess      = "login_success"
	AuditActionLoginFailed       = "login_failed"
	AuditActionEmailVerified     = "email_verified"
	AuditActionEmailVerifyFailed = "email_verify_failed"
	AuditActionLogout            = "logout"
	AuditActionSessionRevoked    = "session_revoked"
	AuditActionRolePreviewSet    = "role_preview_set"
	AuditActionRolePreviewCleared = "role_preview_cleared"
)
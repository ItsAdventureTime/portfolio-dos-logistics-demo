package domain

import "time"

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
	RolePreviewNone                  RolePreview = ""
	RolePreviewLogisticsCoordinator  RolePreview = "logistics_coordinator"
	RolePreviewFreightOpsApprover    RolePreview = "freight_ops_approver"
	RolePreviewDisbursementController RolePreview = "disbursement_controller"
	RolePreviewFinanceOpsLead        RolePreview = "finance_ops_lead"
)

func AllRolePreviews() []RolePreview {
	return []RolePreview{
		RolePreviewNone,
		RolePreviewLogisticsCoordinator,
		RolePreviewFreightOpsApprover,
		RolePreviewDisbursementController,
		RolePreviewFinanceOpsLead,
	}
}

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
	ID            UserID    `json:"user_id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	DisplayName   string    `json:"display_name"`
	EmailVerified bool      `json:"email_verified"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Session is a server-managed session. The plaintext token is never stored;
// only the hash is persisted.
type Session struct {
	ID          SessionID    `json:"id"`
	UserID      UserID       `json:"user_id"`
	TokenHash   string      `json:"-"`
	RolePreview RolePreview  `json:"role_preview"`
	IPAddress   string       `json:"ip_address"`
	UserAgent   string       `json:"user_agent"`
	CreatedAt   time.Time    `json:"created_at"`
	LastSeenAt  time.Time    `json:"last_seen_at"`
	ExpiresAt   time.Time    `json:"expires_at"`
	RevokedAt   *time.Time   `json:"revoked_at,omitempty"`
}

func (s Session) IsRevoked() bool { return s.RevokedAt != nil }
func (s Session) IsExpired(now time.Time) bool { return now.After(s.ExpiresAt) }
func (s Session) IsIdleExpired(now time.Time, idleTimeout time.Duration) bool {
	return now.Sub(s.LastSeenAt) > idleTimeout
}
func (s Session) IsValid(now time.Time, idleTimeout time.Duration) bool {
	return !s.IsRevoked() && !s.IsExpired(now) && !s.IsIdleExpired(now, idleTimeout)
}

// EmailChallenge is a one-time OTP code for email verification.
type EmailChallenge struct {
	ID        string     `json:"id"`
	UserID    UserID     `json:"user_id"`
	CodeHash  string     `json:"-"`
	Purpose   string     `json:"purpose"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	ConsumedAt *time.Time `json:"consumed_at,omitempty"`
}

func (c EmailChallenge) IsExpired(now time.Time) bool { return now.After(c.ExpiresAt) }
func (c EmailChallenge) IsConsumed() bool { return c.ConsumedAt != nil }

// AuditEvent records a security-relevant action. Every sensitive action
// creates an audit event identifying actor, action, entity, time, and
// correlation id. Audit events are append-only.
type AuditEvent struct {
	ID            int64          `json:"id"`
	CorrelationID string         `json:"correlation_id"`
	ActorUserID   *UserID        `json:"actor_user_id,omitempty"`
	ActorRole     string         `json:"actor_role"`
	Action        string         `json:"action"`
	EntityType    string         `json:"entity_type"`
	EntityID      string         `json:"entity_id"`
	Details       map[string]any `json:"details"`
	CreatedAt     time.Time      `json:"created_at"`
}

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
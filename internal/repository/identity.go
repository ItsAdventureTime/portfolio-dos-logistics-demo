// Package repository defines the storage interfaces (consumer-owned). The
// service layer depends on these interfaces; the concrete pgx adapters live
// in internal/store. This keeps business logic independent of the database.
package repository

import (
	"context"
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
)

// UserRepository persists and retrieves users.
type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id domain.UserID) (*domain.User, error)
	GetByUsernameOrEmail(ctx context.Context, identifier string) (*domain.User, error)
	UpdateEmailVerified(ctx context.Context, id domain.UserID, verified bool) error
}

// SessionRepository persists and retrieves sessions. Only the token hash is
// stored; the plaintext token never enters the database.
type SessionRepository interface {
	Create(ctx context.Context, s *domain.Session) error
	GetByTokenHash(ctx context.Context, hash string) (*domain.Session, error)
	UpdateLastSeen(ctx context.Context, id domain.SessionID, at time.Time) error
	UpdateRolePreview(ctx context.Context, id domain.SessionID, preview domain.RolePreview) error
	Revoke(ctx context.Context, id domain.SessionID, at time.Time) error
	RevokeAllForUser(ctx context.Context, userID domain.UserID, at time.Time) error
}

// EmailChallengeRepository persists one-time OTP challenges.
type EmailChallengeRepository interface {
	Create(ctx context.Context, c *domain.EmailChallenge) error
	GetLatestByUserID(ctx context.Context, userID domain.UserID, purpose string) (*domain.EmailChallenge, error)
	MarkConsumed(ctx context.Context, id string, at time.Time) error
}

// AuditEventRepository writes append-only audit events. There is no Update
// or Delete method — by design.
type AuditEventRepository interface {
	Create(ctx context.Context, e *domain.AuditEvent) error
}
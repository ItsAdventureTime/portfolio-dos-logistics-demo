// Package store provides in-memory implementations of the repository
// interfaces for testing. The pgx-backed adapters will live in a separate
// file once a PostgreSQL connection is available in the test environment.
package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
)

// MemUserStore is an in-memory UserRepository for testing.
type MemUserStore struct {
	mu    sync.RWMutex
	users map[string]*domain.User // keyed by ID
	byID  map[domain.UserID]*domain.User
}

// NewMemUserStore creates an empty in-memory user store.
func NewMemUserStore() *MemUserStore {
	return &MemUserStore{
		users: make(map[string]*domain.User),
		byID:  make(map[domain.UserID]*domain.User),
	}
}

// Seed adds a user to the store for test setup.
func (s *MemUserStore) Seed(u *domain.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.Username] = u
	s.users[u.Email] = u
	s.byID[u.ID] = u
}

func (s *MemUserStore) Create(ctx context.Context, u *domain.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.Username] = u
	s.users[u.Email] = u
	s.byID[u.ID] = u
	return nil
}

func (s *MemUserStore) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.byID[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *MemUserStore) GetByUsernameOrEmail(ctx context.Context, identifier string) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[identifier]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *MemUserStore) UpdateEmailVerified(ctx context.Context, id domain.UserID, verified bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.byID[id]
	if !ok {
		return fmt.Errorf("user not found")
	}
	u.EmailVerified = verified
	return nil
}

// MemSessionStore is an in-memory SessionRepository for testing.
type MemSessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*domain.Session // keyed by TokenHash
	byID     map[domain.SessionID]*domain.Session
}

func NewMemSessionStore() *MemSessionStore {
	return &MemSessionStore{
		sessions: make(map[string]*domain.Session),
		byID:     make(map[domain.SessionID]*domain.Session),
	}
}

func (s *MemSessionStore) Create(ctx context.Context, sess *domain.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sess.TokenHash] = sess
	s.byID[sess.ID] = sess
	return nil
}

func (s *MemSessionStore) GetByTokenHash(ctx context.Context, hash string) (*domain.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[hash]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return sess, nil
}

func (s *MemSessionStore) UpdateLastSeen(ctx context.Context, id domain.SessionID, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.byID[id]
	if !ok {
		return fmt.Errorf("session not found")
	}
	sess.LastSeenAt = at
	return nil
}

func (s *MemSessionStore) UpdateRolePreview(ctx context.Context, id domain.SessionID, preview domain.RolePreview) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.byID[id]
	if !ok {
		return fmt.Errorf("session not found")
	}
	sess.RolePreview = preview
	return nil
}

func (s *MemSessionStore) Revoke(ctx context.Context, id domain.SessionID, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.byID[id]
	if !ok {
		return fmt.Errorf("session not found")
	}
	sess.RevokedAt = &at
	return nil
}

func (s *MemSessionStore) RevokeAllForUser(ctx context.Context, userID domain.UserID, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sess := range s.byID {
		if sess.UserID == userID {
			sess.RevokedAt = &at
		}
	}
	return nil
}

// MemEmailChallengeStore is an in-memory EmailChallengeRepository.
type MemEmailChallengeStore struct {
	mu         sync.RWMutex
	challenges map[domain.UserID][]*domain.EmailChallenge
}

func NewMemEmailChallengeStore() *MemEmailChallengeStore {
	return &MemEmailChallengeStore{
		challenges: make(map[domain.UserID][]*domain.EmailChallenge),
	}
}

func (s *MemEmailChallengeStore) Create(ctx context.Context, c *domain.EmailChallenge) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.challenges[c.UserID] = append(s.challenges[c.UserID], c)
	return nil
}

func (s *MemEmailChallengeStore) GetLatestByUserID(ctx context.Context, userID domain.UserID, purpose string) (*domain.EmailChallenge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := s.challenges[userID]
	for i := len(list) - 1; i >= 0; i-- {
		if list[i].Purpose == purpose {
			return list[i], nil
		}
	}
	return nil, fmt.Errorf("challenge not found")
}

func (s *MemEmailChallengeStore) MarkConsumed(ctx context.Context, id string, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, list := range s.challenges {
		for _, c := range list {
			if c.ID == id {
				c.ConsumedAt = &at
				return nil
			}
		}
	}
	return fmt.Errorf("challenge not found")
}

// MemAuditStore is an in-memory AuditEventRepository for testing.
type MemAuditStore struct {
	mu     sync.Mutex
	events []*domain.AuditEvent
}

func NewMemAuditStore() *MemAuditStore {
	return &MemAuditStore{events: make([]*domain.AuditEvent, 0)}
}

func (s *MemAuditStore) Create(ctx context.Context, e *domain.AuditEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, e)
	return nil
}

// Events returns a copy of the stored audit events for test assertions.
func (s *MemAuditStore) Events() []*domain.AuditEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*domain.AuditEvent, len(s.events))
	copy(out, s.events)
	return out
}
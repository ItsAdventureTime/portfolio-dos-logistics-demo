package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
)

// MemClientStore is an in-memory ClientRepository.
type MemClientStore struct {
	mu      sync.RWMutex
	clients map[domain.ClientID]*domain.Client
	byCode  map[string]*domain.Client
}

func NewMemClientStore() *MemClientStore {
	return &MemClientStore{
		clients: make(map[domain.ClientID]*domain.Client),
		byCode:  make(map[string]*domain.Client),
	}
}

func (s *MemClientStore) Seed(c *domain.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[c.ID] = c
	s.byCode[c.Code] = c
}

func (s *MemClientStore) Create(ctx context.Context, c *domain.Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[c.ID] = c
	s.byCode[c.Code] = c
	return nil
}

func (s *MemClientStore) GetByID(ctx context.Context, id domain.ClientID) (*domain.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.clients[id]
	if !ok {
		return nil, fmt.Errorf("client not found")
	}
	return c, nil
}

func (s *MemClientStore) List(ctx context.Context) ([]*domain.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*domain.Client, 0, len(s.clients))
	for _, c := range s.clients {
		out = append(out, c)
	}
	return out, nil
}

func (s *MemClientStore) Update(ctx context.Context, c *domain.Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.clients[c.ID]; !ok {
		return fmt.Errorf("client not found")
	}
	s.clients[c.ID] = c
	return nil
}

// MemQuotationStore is an in-memory QuotationRepository.
type MemQuotationStore struct {
	mu         sync.RWMutex
	quotations map[domain.QuotationID]*domain.Quotation
}

func NewMemQuotationStore() *MemQuotationStore {
	return &MemQuotationStore{
		quotations: make(map[domain.QuotationID]*domain.Quotation),
	}
}

func (s *MemQuotationStore) Seed(q *domain.Quotation) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.quotations[q.ID] = q
}

func (s *MemQuotationStore) Create(ctx context.Context, q *domain.Quotation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.quotations[q.ID] = q
	return nil
}

func (s *MemQuotationStore) GetByID(ctx context.Context, id domain.QuotationID) (*domain.Quotation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.quotations[id]
	if !ok {
		return nil, fmt.Errorf("quotation not found")
	}
	return q, nil
}

func (s *MemQuotationStore) GetByIDWithLines(ctx context.Context, id domain.QuotationID) (*domain.Quotation, error) {
	return s.GetByID(ctx, id)
}

func (s *MemQuotationStore) List(ctx context.Context, clientID *domain.ClientID) ([]*domain.Quotation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*domain.Quotation, 0)
	for _, q := range s.quotations {
		if clientID == nil || q.ClientID == *clientID {
			out = append(out, q)
		}
	}
	return out, nil
}

func (s *MemQuotationStore) UpdateStatus(ctx context.Context, id domain.QuotationID, status domain.QuotationStatus, expectedVersion int) (*domain.Quotation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	q, ok := s.quotations[id]
	if !ok {
		return nil, fmt.Errorf("quotation not found")
	}
	if q.Version != expectedVersion {
		return nil, ErrVersionConflict
	}
	q.Status = status
	q.Version++
	return q, nil
}

func (s *MemQuotationStore) UpdateLines(ctx context.Context, id domain.QuotationID, lines []domain.QuotationLine) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	q, ok := s.quotations[id]
	if !ok {
		return fmt.Errorf("quotation not found")
	}
	q.Lines = lines
	return nil
}

func (s *MemQuotationStore) RecomputeTotals(ctx context.Context, id domain.QuotationID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	q, ok := s.quotations[id]
	if !ok {
		return fmt.Errorf("quotation not found")
	}
	var subtotal int64
	for _, l := range q.Lines {
		subtotal += l.LineTotal
	}
	q.Subtotal = subtotal
	q.Total = subtotal + q.TaxAmount
	return nil
}

// ErrVersionConflict is returned when an optimistic concurrency check fails.
var ErrVersionConflict = fmt.Errorf("version conflict: record was modified by another transaction")
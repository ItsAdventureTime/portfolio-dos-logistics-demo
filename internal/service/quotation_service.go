package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/observability"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrQuotationNotFound      = errors.New("quotation not found")
	ErrQuotationNoLines       = errors.New("quotation must have at least one charge")
	ErrQuotationInvalidCurrency = errors.New("invalid currency code")
	ErrInvalidTransition      = errors.New("invalid state transition")
	ErrVersionConflict        = errors.New("version conflict: the record was modified by another transaction")
	ErrQuotationNotAccepted   = errors.New("quotation must be accepted before creating a funding request")
)

// QuotationService handles the quotation workflow.
type QuotationService struct {
	clients    repository.ClientRepository
	quotations repository.QuotationRepository
	audit      repository.AuditEventRepository
	log        *slog.Logger
}

func NewQuotationService(
	clients repository.ClientRepository,
	quotations repository.QuotationRepository,
	audit repository.AuditEventRepository,
	log *slog.Logger,
) *QuotationService {
	return &QuotationService{
		clients:    clients,
		quotations: quotations,
		audit:      audit,
		log:        log,
	}
}

// CreateClient creates a new client.
func (s *QuotationService) CreateClient(ctx context.Context, name, code, email, phone, address string) (*domain.Client, error) {
	c := &domain.Client{
		ID:           domain.ClientID(uuid.NewString()),
		Name:         name,
		Code:         code,
		ContactEmail: email,
		ContactPhone: phone,
		Address:      address,
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Version:      1,
	}
	if err := s.clients.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	s.auditQuotation(ctx, nil, domain.AuditActionClientCreated, "client", string(c.ID))
	return c, nil
}

// CreateQuotationInput holds the data for creating a new quotation.
type CreateQuotationInput struct {
	ClientID     domain.ClientID
	CurrencyCode string
	Notes        string
	Lines        []QuotationLineInput
}

type QuotationLineInput struct {
	Description string
	Quantity    int64
	UnitPrice   int64
}

// CreateQuotation creates a new quotation in draft status. Validates that
// the client exists, the currency is valid, and at least one charge line exists.
func (s *QuotationService) CreateQuotation(ctx context.Context, createdBy domain.UserID, input CreateQuotationInput) (*domain.Quotation, error) {
	corrID := observability.CorrelationFrom(ctx)

	if len(input.Lines) == 0 {
		return nil, ErrQuotationNoLines
	}
	if input.CurrencyCode == "" {
		return nil, ErrQuotationInvalidCurrency
	}

	// Verify client exists.
	client, err := s.clients.GetByID(ctx, input.ClientID)
	if err != nil || client == nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}
	if !client.IsActive {
		return nil, fmt.Errorf("client is not active")
	}

	// Build quotation lines.
	lines := make([]domain.QuotationLine, len(input.Lines))
	var subtotal int64
	for i, li := range input.Lines {
		lineTotal := li.Quantity * li.UnitPrice
		lines[i] = domain.QuotationLine{
			ID:           uuid.NewString(),
			QuotationID:  "", // set after quotation is created
			Description:  li.Description,
			Quantity:     li.Quantity,
			UnitPrice:    li.UnitPrice,
			LineTotal:    lineTotal,
			SortOrder:    i,
			CreatedAt:     time.Now().UTC(),
		}
		subtotal += lineTotal
	}

	q := &domain.Quotation{
		ID:              domain.QuotationID(uuid.NewString()),
		ClientID:        input.ClientID,
		QuotationNumber: fmt.Sprintf("Q-%s-%d", client.Code, time.Now().Unix()),
		Status:          domain.QuotationStatusDraft,
		CurrencyCode:    input.CurrencyCode,
		Subtotal:        subtotal,
		TaxAmount:       0,
		Total:           subtotal,
		Notes:           input.Notes,
		CreatedBy:       &createdBy,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		Version:         1,
		Lines:           lines,
	}

	if err := s.quotations.Create(ctx, q); err != nil {
		return nil, fmt.Errorf("create quotation: %w", err)
	}

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID,
		ActorUserID:   &createdBy,
		ActorRole:     "administrator",
		Action:        domain.AuditActionQuotationCreated,
		EntityType:    "quotation",
		EntityID:      string(q.ID),
		Details:       map[string]any{"quotation_number": q.QuotationNumber},
	})

	return q, nil
}

// TransitionQuotation performs a state transition with optimistic concurrency.
// It validates the transition is legal, checks the version, and records an
// audit event. Returns the updated quotation.
func (s *QuotationService) TransitionQuotation(ctx context.Context, id domain.QuotationID, target domain.QuotationStatus, expectedVersion int, actor domain.UserID) (*domain.Quotation, error) {
	corrID := observability.CorrelationFrom(ctx)

	q, err := s.quotations.GetByID(ctx, id)
	if err != nil || q == nil {
		return nil, ErrQuotationNotFound
	}

	if !domain.CanTransitionQuotation(q.Status, target) {
		return nil, fmt.Errorf("%w: %s → %s", ErrInvalidTransition, q.Status, target)
	}

	// Terminal statuses cannot transition.
	if domain.IsTerminalQuotationStatus(q.Status) {
		return nil, fmt.Errorf("%w: quotation is in terminal state %s", ErrInvalidTransition, q.Status)
	}

	updated, err := s.quotations.UpdateStatus(ctx, id, target, expectedVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVersionConflict, err)
	}

	// Set timestamp fields.
	now := time.Now().UTC()
	switch target {
	case domain.QuotationStatusAccepted:
		updated.AcceptedAt = &now
	case domain.QuotationStatusVoid:
		updated.VoidedAt = &now
	}

	action := domain.AuditActionQuotationCreated // default
	switch target {
	case domain.QuotationStatusReview:
		action = domain.AuditActionQuotationSubmitted
	case domain.QuotationStatusApproved:
		action = domain.AuditActionQuotationApproved
	case domain.QuotationStatusAccepted:
		action = domain.AuditActionQuotationAccepted
	case domain.QuotationStatusRevised:
		action = domain.AuditActionQuotationRevised
	case domain.QuotationStatusVoid:
		action = domain.AuditActionQuotationVoided
	}

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID,
		ActorUserID:   &actor,
		ActorRole:     "administrator",
		Action:        action,
		EntityType:    "quotation",
		EntityID:      string(id),
		Details: map[string]any{
			"from_status":    string(q.Status),
			"to_status":      string(target),
			"new_version":    updated.Version,
			"expected_version": expectedVersion,
		},
	})

	return updated, nil
}

// GetQuotation returns a quotation with its lines.
func (s *QuotationService) GetQuotation(ctx context.Context, id domain.QuotationID) (*domain.Quotation, error) {
	q, err := s.quotations.GetByIDWithLines(ctx, id)
	if err != nil || q == nil {
		return nil, ErrQuotationNotFound
	}
	return q, nil
}

// ListQuotations lists quotations, optionally filtered by client.
func (s *QuotationService) ListQuotations(ctx context.Context, clientID *domain.ClientID) ([]*domain.Quotation, error) {
	return s.quotations.List(ctx, clientID)
}

// ListClients returns all clients.
func (s *QuotationService) ListClients(ctx context.Context) ([]*domain.Client, error) {
	return s.clients.List(ctx)
}

func (s *QuotationService) auditQuotation(ctx context.Context, userID *domain.UserID, action, entityType, entityID string) {
	corrID := observability.CorrelationFrom(ctx)
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
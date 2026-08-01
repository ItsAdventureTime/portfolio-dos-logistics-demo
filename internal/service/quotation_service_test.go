package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/observability"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/store"
	"github.com/google/uuid"
)

func newQuotationService(t *testing.T) (*QuotationService, *store.MemClientStore, *store.MemQuotationStore, *store.MemAuditStore) {
	t.Helper()
	clients := store.NewMemClientStore()
	quotations := store.NewMemQuotationStore()
	audit := store.NewMemAuditStore()
	svc := NewQuotationService(clients, quotations, audit, slog.Default())
	return svc, clients, quotations, audit
}

func seedClient(t *testing.T, clients *store.MemClientStore) *domain.Client {
	t.Helper()
	c := &domain.Client{
		ID:           domain.ClientID(uuid.NewString()),
		Name:         "Acme Logistics Inc.",
		Code:         "ACME",
		ContactEmail: "contact@acme.example",
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Version:      1,
	}
	clients.Seed(c)
	return c
}

func actor() domain.UserID {
	return domain.UserID(uuid.NewString())
}

// Acceptance: "A quotation must have at least one charge and a valid currency."
func TestCreateQuotation_NoLines_Rejected(t *testing.T) {
	svc, clients, _, _ := newQuotationService(t)
	client := seedClient(t, clients)
	ctx := context.Background()
	_, err := svc.CreateQuotation(ctx, actor(), CreateQuotationInput{
		ClientID:     client.ID,
		CurrencyCode: "USD",
		Lines:        []QuotationLineInput{},
	})
	if !errors.Is(err, ErrQuotationNoLines) {
		t.Errorf("expected ErrQuotationNoLines, got %v", err)
	}
}

func TestCreateQuotation_NoCurrency_Rejected(t *testing.T) {
	svc, clients, _, _ := newQuotationService(t)
	client := seedClient(t, clients)
	ctx := context.Background()
	_, err := svc.CreateQuotation(ctx, actor(), CreateQuotationInput{
		ClientID:     client.ID,
		CurrencyCode: "",
		Lines:        []QuotationLineInput{{Description: "Freight", Quantity: 1, UnitPrice: 10000}},
	})
	if !errors.Is(err, ErrQuotationInvalidCurrency) {
		t.Errorf("expected ErrQuotationInvalidCurrency, got %v", err)
	}
}

// Acceptance: valid quotation creates in draft status with correct totals.
func TestCreateQuotation_Valid_CreatesDraft(t *testing.T) {
	svc, clients, _, audit := newQuotationService(t)
	client := seedClient(t, clients)
	ctx := observability.WithCorrelation(context.Background(), "corr-q")
	a := actor()
	q, err := svc.CreateQuotation(ctx, a, CreateQuotationInput{
		ClientID:     client.ID,
		CurrencyCode: "USD",
		Lines: []QuotationLineInput{
			{Description: "Ocean freight", Quantity: 1, UnitPrice: 500000}, // $5000.00
			{Description: "Insurance", Quantity: 1, UnitPrice: 25000},      // $250.00
		},
	})
	if err != nil {
		t.Fatalf("CreateQuotation: %v", err)
	}
	if q.Status != domain.QuotationStatusDraft {
		t.Errorf("status = %s, want draft", q.Status)
	}
	if q.Subtotal != 525000 {
		t.Errorf("subtotal = %d, want 525000", q.Subtotal)
	}
	if q.Total != 525000 {
		t.Errorf("total = %d, want 525000", q.Total)
	}
	if q.Version != 1 {
		t.Errorf("version = %d, want 1", q.Version)
	}
	if len(q.Lines) != 2 {
		t.Errorf("lines = %d, want 2", len(q.Lines))
	}
	// Audit event.
	events := audit.Events()
	found := false
	for _, e := range events {
		if e.Action == domain.AuditActionQuotationCreated {
			found = true
			if e.CorrelationID != "corr-q" {
				t.Errorf("correlation_id = %s, want corr-q", e.CorrelationID)
			}
		}
	}
	if !found {
		t.Error("missing quotation_created audit event")
	}
}

// Acceptance: "Quotation-to-collection workflow enforces its state transitions."
func TestTransitionQuotation_LegalTransition_Succeeds(t *testing.T) {
	svc, clients, _, _ := newQuotationService(t)
	client := seedClient(t, clients)
	ctx := context.Background()
	q, _ := svc.CreateQuotation(ctx, actor(), CreateQuotationInput{
		ClientID:     client.ID,
		CurrencyCode: "USD",
		Lines:        []QuotationLineInput{{Description: "Freight", Quantity: 1, UnitPrice: 10000}},
	})
	// draft → review → approved → accepted
	q2, err := svc.TransitionQuotation(ctx, q.ID, domain.QuotationStatusReview, q.Version, actor())
	if err != nil {
		t.Fatalf("draft→review: %v", err)
	}
	q3, err := svc.TransitionQuotation(ctx, q.ID, domain.QuotationStatusApproved, q2.Version, actor())
	if err != nil {
		t.Fatalf("review→approved: %v", err)
	}
	_, err = svc.TransitionQuotation(ctx, q.ID, domain.QuotationStatusAccepted, q3.Version, actor())
	if err != nil {
		t.Fatalf("approved→accepted: %v", err)
	}
}

// Acceptance: illegal transitions are rejected.
func TestTransitionQuotation_IllegalTransition_Rejected(t *testing.T) {
	svc, clients, _, _ := newQuotationService(t)
	client := seedClient(t, clients)
	ctx := context.Background()
	q, _ := svc.CreateQuotation(ctx, actor(), CreateQuotationInput{
		ClientID:     client.ID,
		CurrencyCode: "USD",
		Lines:        []QuotationLineInput{{Description: "Freight", Quantity: 1, UnitPrice: 10000}},
	})
	// draft → accepted is illegal (must go through review → approved first)
	_, err := svc.TransitionQuotation(ctx, q.ID, domain.QuotationStatusAccepted, q.Version, actor())
	if !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("expected ErrInvalidTransition for draft→accepted, got %v", err)
	}
}

// Acceptance: terminal status cannot transition.
func TestTransitionQuotation_TerminalStatus_Rejected(t *testing.T) {
	svc, clients, _, _ := newQuotationService(t)
	client := seedClient(t, clients)
	ctx := context.Background()
	q, _ := svc.CreateQuotation(ctx, actor(), CreateQuotationInput{
		ClientID:     client.ID,
		CurrencyCode: "USD",
		Lines:        []QuotationLineInput{{Description: "Freight", Quantity: 1, UnitPrice: 10000}},
	})
	// draft → void (terminal)
	q2, _ := svc.TransitionQuotation(ctx, q.ID, domain.QuotationStatusVoid, q.Version, actor())
	// void → draft should fail (terminal)
	_, err := svc.TransitionQuotation(ctx, q.ID, domain.QuotationStatusDraft, q2.Version, actor())
	if !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("expected ErrInvalidTransition from terminal state, got %v", err)
	}
}

// Acceptance: "Every state-changing request must protect against stale record versions."
func TestTransitionQuotation_StaleVersion_Rejected(t *testing.T) {
	svc, clients, _, _ := newQuotationService(t)
	client := seedClient(t, clients)
	ctx := context.Background()
	q, _ := svc.CreateQuotation(ctx, actor(), CreateQuotationInput{
		ClientID:     client.ID,
		CurrencyCode: "USD",
		Lines:        []QuotationLineInput{{Description: "Freight", Quantity: 1, UnitPrice: 10000}},
	})
	originalVersion := q.Version // capture before any mutation
	// Transition with the correct version first (incrementing to 2).
	_, err := svc.TransitionQuotation(ctx, q.ID, domain.QuotationStatusReview, originalVersion, actor())
	if err != nil {
		t.Fatalf("first transition: %v", err)
	}
	// Now try with the stale version (1 instead of 2).
	_, err = svc.TransitionQuotation(ctx, q.ID, domain.QuotationStatusApproved, originalVersion, actor())
	if !errors.Is(err, ErrVersionConflict) {
		t.Errorf("expected ErrVersionConflict for stale version, got %v", err)
	}
}

// Acceptance: audit events identify actor, action, entity, time, correlation ID.
func TestTransitionQuotation_AuditEventComplete(t *testing.T) {
	svc, clients, _, audit := newQuotationService(t)
	client := seedClient(t, clients)
	ctx := observability.WithCorrelation(context.Background(), "corr-trans")
	a := actor()
	q, _ := svc.CreateQuotation(ctx, a, CreateQuotationInput{
		ClientID:     client.ID,
		CurrencyCode: "USD",
		Lines:        []QuotationLineInput{{Description: "Freight", Quantity: 1, UnitPrice: 10000}},
	})
	_, err := svc.TransitionQuotation(ctx, q.ID, domain.QuotationStatusReview, q.Version, a)
	if err != nil {
		t.Fatalf("transition: %v", err)
	}
	events := audit.Events()
	var found *domain.AuditEvent
	for _, e := range events {
		if e.Action == domain.AuditActionQuotationSubmitted {
			found = e
			break
		}
	}
	if found == nil {
		t.Fatal("missing quotation_submitted audit event")
	}
	if found.ActorUserID == nil || *found.ActorUserID != a {
		t.Error("audit actor mismatch")
	}
	if found.EntityType != "quotation" {
		t.Errorf("entity_type = %s, want quotation", found.EntityType)
	}
	if found.EntityID != string(q.ID) {
		t.Errorf("entity_id = %s, want %s", found.EntityID, q.ID)
	}
	if found.CorrelationID != "corr-trans" {
		t.Errorf("correlation_id = %s, want corr-trans", found.CorrelationID)
	}
}
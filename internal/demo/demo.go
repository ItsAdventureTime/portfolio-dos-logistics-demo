// Package demo provides fictional demonstration fixtures and a controlled
// reset that removes only synthetic demonstration records. All fixture
// emails use the .example TLD. All amounts are synthetic. No customer data,
// production credentials, or confidential business records appear anywhere
// in the fixtures.
package demo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
	"github.com/google/uuid"
)

// DemoFixtures contains fictional data for the demonstration environment.
// All emails use .example TLD. All amounts are synthetic.
type DemoFixtures struct {
	Users      []*domain.User
	Clients    []*domain.Client
	Quotations []*domain.Quotation
}

// Seed returns the fictional demo fixtures.
func Seed() *DemoFixtures {
	now := time.Now().UTC()

	// --- Users ---
	admin := &domain.User{
		ID:            domain.UserID("demo-admin-001"),
		Username:      "admin",
		Email:         "admin@dosfreightflow.example",
		PasswordHash:  "$argon2id$v=19$m=65536,t=3,p=2$AAAAAAAAAAAAAAAAAAAAAA$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		DisplayName:   "Demo Administrator",
		EmailVerified: true,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	coordinator := &domain.User{
		ID:            domain.UserID("demo-user-002"),
		Username:      "j.delacruz",
		Email:         "j.delacruz@dosfreightflow.example",
		PasswordHash:  "$argon2id$v=19$m=65536,t=3,p=2$AAAAAAAAAAAAAAAAAAAAAA$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		DisplayName:   "Juan Dela Cruz",
		EmailVerified: true,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// --- Clients ---
	clients := []*domain.Client{
		{
			ID:           domain.ClientID("demo-client-acme"),
			Name:         "Acme Freight Logistics Inc.",
			Code:         "ACME",
			ContactEmail:  "operations@acme.example",
			ContactPhone:  "+1-555-0100",
			Address:       "123 Harbor Drive, Wilmington, DE 19801",
			IsActive:      true,
			CreatedAt:     now,
			UpdatedAt:     now,
			Version:       1,
		},
		{
			ID:           domain.ClientID("demo-client-globex"),
			Name:         "Globex Shipping Co.",
			Code:         "GLOBX",
			ContactEmail:  "finance@globex.example",
			ContactPhone:  "+1-555-0200",
			Address:       "456 Port Authority Blvd, Long Beach, CA 90802",
			IsActive:      true,
			CreatedAt:     now,
			UpdatedAt:     now,
			Version:       1,
		},
		{
			ID:           domain.ClientID("demo-client-stark"),
			Name:         "Stark Industries Transport",
			Code:         "STARK",
			ContactEmail:  "ap@stark.example",
			ContactPhone:  "+1-555-0300",
			Address:       "789 Industrial Parkway, Houston, TX 77001",
			IsActive:      true,
			CreatedAt:     now,
			UpdatedAt:     now,
			Version:       1,
		},
	}

	// --- Quotations ---
	quotations := []*domain.Quotation{
		{
			ID:              domain.QuotationID("demo-quot-001"),
			ClientID:        clients[0].ID,
			QuotationNumber:  "Q-ACME-100001",
			Status:          domain.QuotationStatusAccepted,
			CurrencyCode:    "USD",
			Subtotal:        4250000, // $42,500.00
			TaxAmount:       0,
			Total:           4250000,
			Notes:           "Ocean freight - 40ft container, Manila to LA",
			CreatedBy:       &admin.ID,
			CreatedAt:       now,
			UpdatedAt:       now,
			AcceptedAt:      &now,
			Version:         3,
			Lines: []domain.QuotationLine{
				{ID: uuid.NewString(), Description: "Ocean freight 40ft", Quantity: 1, UnitPrice: 3500000, LineTotal: 3500000},
				{ID: uuid.NewString(), Description: "Origin handling", Quantity: 1, UnitPrice: 500000, LineTotal: 500000},
				{ID: uuid.NewString(), Description: "Insurance", Quantity: 1, UnitPrice: 250000, LineTotal: 250000},
			},
		},
		{
			ID:              domain.QuotationID("demo-quot-002"),
			ClientID:        clients[1].ID,
			QuotationNumber:  "Q-GLOBX-100002",
			Status:          domain.QuotationStatusDraft,
			CurrencyCode:    "USD",
			Subtotal:        1850000,
			TaxAmount:       0,
			Total:           1850000,
			Notes:           "Air freight - JFK to Heathrow",
			CreatedBy:       &coordinator.ID,
			CreatedAt:       now,
			UpdatedAt:       now,
			Version:         1,
			Lines: []domain.QuotationLine{
				{ID: uuid.NewString(), Description: "Air freight express", Quantity: 1, UnitPrice: 1500000, LineTotal: 1500000},
				{ID: uuid.NewString(), Description: "Customs clearance", Quantity: 1, UnitPrice: 350000, LineTotal: 350000},
			},
		},
	}

	return &DemoFixtures{
		Users:      []*domain.User{admin, coordinator},
		Clients:    clients,
		Quotations: quotations,
	}
}

// DemoPrefix is the prefix used for all demo record IDs.
const DemoPrefix = "demo-"

// IsDemoRecord reports whether a record ID belongs to the demo fixtures.
func IsDemoRecord(id string) bool {
	return len(id) >= len(DemoPrefix) && id[:len(DemoPrefix)] == DemoPrefix
}

// ResetService provides a controlled reset that removes only demo-prefixed
// records. Non-demo data is never touched.
type ResetService struct {
	log *slog.Logger
}

// NewResetService creates a reset service.
func NewResetService(log *slog.Logger) *ResetService {
	return &ResetService{log: log}
}

// ResetResult records what was removed.
type ResetResult struct {
	UsersRemoved      int
	ClientsRemoved    int
	QuotationsRemoved int
	DocumentsRemoved  int
}

// Reset removes only demo-prefixed records. The caller provides cleanup
// functions for each entity type, and only records whose IDs start with
// "demo-" are removed.
func (s *ResetService) Reset(
	ctx context.Context,
	actor domain.UserID,
	removeUsers func(ctx context.Context, predicate func(id string) bool) int,
	removeClients func(ctx context.Context, predicate func(id string) bool) int,
	removeQuotations func(ctx context.Context, predicate func(id string) bool) int,
	removeDocuments func(ctx context.Context, predicate func(id string) bool) int,
) (*ResetResult, error) {
	result := &ResetResult{
		UsersRemoved:      removeUsers(ctx, IsDemoRecord),
		ClientsRemoved:    removeClients(ctx, IsDemoRecord),
		QuotationsRemoved: removeQuotations(ctx, IsDemoRecord),
		DocumentsRemoved:  removeDocuments(ctx, IsDemoRecord),
	}

	s.log.InfoContext(ctx, "demo reset completed",
		"users_removed", result.UsersRemoved,
		"clients_removed", result.ClientsRemoved,
		"quotations_removed", result.QuotationsRemoved,
		"documents_removed", result.DocumentsRemoved,
		"actor", string(actor),
	)

	return result, nil
}

// EnvironmentMarker returns a marker string identifying the environment
// as a demonstration, for display in the UI.
func EnvironmentMarker(env string) string {
	if env == "demo" || env == "development" {
		return fmt.Sprintf("DEMONSTRATION ENVIRONMENT — %s", env)
	}
	return ""
}
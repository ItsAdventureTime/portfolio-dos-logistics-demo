package repository

import (
	"context"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
)

// ClientRepository persists and retrieves clients.
type ClientRepository interface {
	Create(ctx context.Context, c *domain.Client) error
	GetByID(ctx context.Context, id domain.ClientID) (*domain.Client, error)
	List(ctx context.Context) ([]*domain.Client, error)
	Update(ctx context.Context, c *domain.Client) error
}

// QuotationRepository persists quotations and their lines.
type QuotationRepository interface {
	Create(ctx context.Context, q *domain.Quotation) error
	GetByID(ctx context.Context, id domain.QuotationID) (*domain.Quotation, error)
	GetByIDWithLines(ctx context.Context, id domain.QuotationID) (*domain.Quotation, error)
	List(ctx context.Context, clientID *domain.ClientID) ([]*domain.Quotation, error)
	UpdateStatus(ctx context.Context, id domain.QuotationID, status domain.QuotationStatus, expectedVersion int) (*domain.Quotation, error)
	UpdateLines(ctx context.Context, id domain.QuotationID, lines []domain.QuotationLine) error
	RecomputeTotals(ctx context.Context, id domain.QuotationID) error
}
package repository

import (
	"context"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
)

// --- Funding Source ---

type FundingSourceRepository interface {
	Create(ctx context.Context, fs *domain.FundingSource) error
	GetByID(ctx context.Context, id domain.FundingSourceID) (*domain.FundingSource, error)
	List(ctx context.Context) ([]*domain.FundingSource, error)
}

// --- Budget Request ---

type BudgetRequestRepository interface {
	Create(ctx context.Context, br *domain.BudgetRequest) error
	GetByID(ctx context.Context, id domain.BudgetRequestID) (*domain.BudgetRequest, error)
	GetByIDWithLines(ctx context.Context, id domain.BudgetRequestID) (*domain.BudgetRequest, error)
	List(ctx context.Context, clientID *domain.ClientID) ([]*domain.BudgetRequest, error)
	UpdateStatus(ctx context.Context, id domain.BudgetRequestID, status domain.BudgetRequestStatus, expectedVersion int) (*domain.BudgetRequest, error)
}

// --- Approval Decision ---

type ApprovalDecisionRepository interface {
	Create(ctx context.Context, d *domain.ApprovalDecision) error
	ListByRequestID(ctx context.Context, reqID domain.BudgetRequestID) ([]*domain.ApprovalDecision, error)
}

// --- Disbursement ---

type DisbursementRepository interface {
	Create(ctx context.Context, d *domain.Disbursement) error
	GetByID(ctx context.Context, id domain.DisbursementID) (*domain.Disbursement, error)
	List(ctx context.Context, budgetRequestID *domain.BudgetRequestID) ([]*domain.Disbursement, error)
	UpdateStatus(ctx context.Context, id domain.DisbursementID, status domain.DisbursementStatus, expectedVersion int) (*domain.Disbursement, error)
}

// --- Payment Proof ---

type PaymentProofRepository interface {
	Create(ctx context.Context, p *domain.PaymentProof) error
	ListByDisbursementID(ctx context.Context, dID domain.DisbursementID) ([]*domain.PaymentProof, error)
}

// --- Liquidation ---

type LiquidationRepository interface {
	Create(ctx context.Context, l *domain.Liquidation) error
	GetByID(ctx context.Context, id domain.LiquidationID) (*domain.Liquidation, error)
	GetByDisbursementID(ctx context.Context, dID domain.DisbursementID) (*domain.Liquidation, error)
	UpdateStatus(ctx context.Context, id domain.LiquidationID, status domain.LiquidationStatus, actualAmount, variance int64, expectedVersion int) (*domain.Liquidation, error)
}

// --- Liquidation Evidence ---

type LiquidationEvidenceRepository interface {
	Create(ctx context.Context, e *domain.LiquidationEvidence) error
	ListByLiquidationID(ctx context.Context, lID domain.LiquidationID) ([]*domain.LiquidationEvidence, error)
}

// --- Billing Record ---

type BillingRecordRepository interface {
	Create(ctx context.Context, b *domain.BillingRecord) error
	GetByID(ctx context.Context, id domain.BillingRecordID) (*domain.BillingRecord, error)
	GetByIDWithLines(ctx context.Context, id domain.BillingRecordID) (*domain.BillingRecord, error)
	List(ctx context.Context, clientID *domain.ClientID) ([]*domain.BillingRecord, error)
	UpdateStatus(ctx context.Context, id domain.BillingRecordID, status domain.BillingStatus, expectedVersion int) (*domain.BillingRecord, error)
	UpdateTotals(ctx context.Context, id domain.BillingRecordID, subtotal, tax, total int64) error
	MarkReplaced(ctx context.Context, id domain.BillingRecordID, replacedBy domain.BillingRecordID) error
}

// --- Credit Memo ---

type CreditMemoRepository interface {
	Create(ctx context.Context, c *domain.CreditMemo) error
	ListByClientID(ctx context.Context, cID domain.ClientID) ([]*domain.CreditMemo, error)
}

// --- Client Payment ---

type ClientPaymentRepository interface {
	Create(ctx context.Context, p *domain.ClientPayment) error
	GetByID(ctx context.Context, id domain.ClientPaymentID) (*domain.ClientPayment, error)
	List(ctx context.Context, clientID *domain.ClientID) ([]*domain.ClientPayment, error)
}

// --- Billing Allocation ---

type BillingAllocationRepository interface {
	Create(ctx context.Context, a *domain.BillingAllocation) error
	ListByPaymentID(ctx context.Context, pID domain.ClientPaymentID) ([]*domain.BillingAllocation, error)
	ListByBillingID(ctx context.Context, bID domain.BillingRecordID) ([]*domain.BillingAllocation, error)
}

// --- Supporting Document ---

type SupportingDocumentRepository interface {
	Create(ctx context.Context, d *domain.LiquidationEvidence) error
	ListByEntity(ctx context.Context, entityType string, entityID string) ([]*domain.LiquidationEvidence, error)
}
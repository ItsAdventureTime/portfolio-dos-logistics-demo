// Package demo wires the in-memory stores and seeds the demo data for local
// demo mode. This avoids needing PostgreSQL to run the demo.
package demo

import (
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/auth"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/service"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/store"
	"log/slog"
)

// BootstrapResult holds the seeded services for demo mode.
type BootstrapResult struct {
	AuthService     *service.AuthService
	QuotationService *service.QuotationService
	WorkflowService *service.WorkflowService
	OTPCfg          auth.OTPConfig
}

// Bootstrap creates in-memory stores, seeds a demo admin user with clients,
// quotations, and workflow records, and returns ready-to-use services.
func Bootstrap(log *slog.Logger) BootstrapResult {
	// --- Auth stores ---
	users := store.NewMemUserStore()
	sessions := store.NewMemSessionStore()
	challenges := store.NewMemEmailChallengeStore()
	audit := store.NewMemAuditStore()

	// Seed the demo admin user.
	hash, _ := auth.HashPassword("Password123!")
	now := time.Now().UTC()
	admin := &domain.User{
		ID:            domain.UserID("demo-admin-001"),
		Username:      "admin",
		Email:         "admin@dosfreightflow.example",
		PasswordHash:  hash,
		DisplayName:   "Demo Administrator",
		EmailVerified: true,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	users.Seed(admin)

	otpCfg := auth.DefaultOTPConfig([]byte("demo-otp-secret-32-bytes-minimum-ok"))
	authSvc := service.NewAuthService(
		users, sessions, challenges, audit,
		otpCfg, 24*time.Hour, 1*time.Hour, true, log,
	)

	// --- Workflow stores ---
	clients := store.NewMemClientStore()
	quotations := store.NewMemQuotationStore()
	fundingSrc := store.NewMemFundingSourceStore()
	budgetReqs := store.NewMemBudgetRequestStore()
	approvals := store.NewMemApprovalDecisionStore()
	disbursements := store.NewMemDisbursementStore()
	payProofs := store.NewMemPaymentProofStore()
	liquidations := store.NewMemLiquidationStore()
	liqEvidence := store.NewMemLiquidationEvidenceStore()
	billing := store.NewMemBillingRecordStore()
	creditMemos := store.NewMemCreditMemoStore()
	payments := store.NewMemClientPaymentStore()
	allocations := store.NewMemBillingAllocationStore()

	// Seed clients
	acme := &domain.Client{
		ID: domain.ClientID("demo-client-acme"), Name: "Acme Freight Logistics Inc.",
		Code: "ACME", ContactEmail: "operations@acme.example",
		ContactPhone: "+1-555-0100", Address: "123 Harbor Drive, Wilmington, DE 19801",
		IsActive: true, CreatedAt: now, UpdatedAt: now, Version: 1,
	}
	globex := &domain.Client{
		ID: domain.ClientID("demo-client-globex"), Name: "Globex Shipping Co.",
		Code: "GLOBX", ContactEmail: "finance@globex.example",
		ContactPhone: "+1-555-0200", Address: "456 Port Authority Blvd, Long Beach, CA 90802",
		IsActive: true, CreatedAt: now, UpdatedAt: now, Version: 1,
	}
	stark := &domain.Client{
		ID: domain.ClientID("demo-client-stark"), Name: "Stark Industries Transport",
		Code: "STARK", ContactEmail: "ap@stark.example",
		ContactPhone: "+1-555-0300", Address: "789 Industrial Parkway, Houston, TX 77001",
		IsActive: true, CreatedAt: now, UpdatedAt: now, Version: 1,
	}
	clients.Seed(acme)
	clients.Seed(globex)
	clients.Seed(stark)

	// Seed funding source
	fs := &domain.FundingSource{
		ID: domain.FundingSourceID("demo-fs-bank-a"), Name: "Bank A Operating Capital",
		Code: "BANKA", IsApproved: true, BalanceCents: 10000000, CurrencyCode: "USD",
		CreatedAt: now, UpdatedAt: now, Version: 1,
	}
	fundingSrc.Seed(fs)

	// Seed quotations
	quotAccepted := &domain.Quotation{
		ID: domain.QuotationID("demo-quot-001"), ClientID: acme.ID,
		QuotationNumber: "Q-ACME-100001", Status: domain.QuotationStatusAccepted,
		CurrencyCode: "USD", Subtotal: 4250000, TaxAmount: 0, Total: 4250000,
		Notes: "Ocean freight: 40ft container, Manila to LA",
		CreatedBy: &admin.ID, CreatedAt: now, UpdatedAt: now, AcceptedAt: &now, Version: 3,
		Lines: []domain.QuotationLine{
			{ID: "demo-ql-001", Description: "Ocean freight 40ft", Quantity: 1, UnitPrice: 3500000, LineTotal: 3500000},
			{ID: "demo-ql-002", Description: "Origin handling", Quantity: 1, UnitPrice: 500000, LineTotal: 500000},
			{ID: "demo-ql-003", Description: "Insurance", Quantity: 1, UnitPrice: 250000, LineTotal: 250000},
		},
	}
	quotDraft := &domain.Quotation{
		ID: domain.QuotationID("demo-quot-002"), ClientID: globex.ID,
		QuotationNumber: "Q-GLOBX-100002", Status: domain.QuotationStatusDraft,
		CurrencyCode: "USD", Subtotal: 1850000, TaxAmount: 0, Total: 1850000,
		Notes: "Air freight: JFK to Heathrow",
		CreatedBy: &admin.ID, CreatedAt: now, UpdatedAt: now, Version: 1,
		Lines: []domain.QuotationLine{
			{ID: "demo-ql-004", Description: "Air freight express", Quantity: 1, UnitPrice: 1500000, LineTotal: 1500000},
			{ID: "demo-ql-005", Description: "Customs clearance", Quantity: 1, UnitPrice: 350000, LineTotal: 350000},
		},
	}
	quotations.Seed(quotAccepted)
	quotations.Seed(quotDraft)

	// Seed a budget request (approved, linked to accepted quotation)
	budgetReq := &domain.BudgetRequest{
		ID: domain.BudgetRequestID("demo-br-001"), QuotationID: quotAccepted.ID,
		ClientID: acme.ID, RequestNumber: "BR-ACME-100001",
		Status: domain.BudgetRequestStatusApproved, CurrencyCode: "USD",
		AmountCents: 4250000, Purpose: "Ocean freight funding",
		CreatedBy: &admin.ID, CreatedAt: now, UpdatedAt: now, ApprovedAt: &now, Version: 3,
	}
	budgetReqs.Seed(budgetReq)

	// Seed a disbursement (released)
	disbursement := &domain.Disbursement{
		ID: domain.DisbursementID("demo-disb-001"), BudgetRequestID: budgetReq.ID,
		FundingSourceID: fs.ID, Status: domain.DisbursementStatusReleased,
		AmountCents: 4250000, CurrencyCode: "USD",
		ReferenceNumber: "WIRE-2026-001", Notes: "Initial disbursement",
		CreatedBy: &admin.ID, CreatedAt: now, UpdatedAt: now, ReleasedAt: &now, Version: 2,
	}
	disbursements.Seed(disbursement)

	// Seed a liquidation (open)
	liquidation := &domain.Liquidation{
		ID: domain.LiquidationID("demo-liq-001"), DisbursementID: disbursement.ID,
		Status: domain.LiquidationStatusOpen, ReleasedAmount: 4250000,
		ActualAmount: 0, VarianceAmount: 4250000,
		CreatedBy: &admin.ID, CreatedAt: now, UpdatedAt: now, Version: 1,
	}
	liquidations.Seed(liquidation)

	// Seed a billing record (finalized)
	billingRec := &domain.BillingRecord{
		ID: domain.BillingRecordID("demo-bill-001"), ClientID: acme.ID,
		BudgetRequestID: &budgetReq.ID, BillingNumber: "INV-ACME-100001",
		Status: domain.BillingStatusFinalized, CurrencyCode: "USD",
		Subtotal: 4250000, TaxAmount: 0, Total: 4250000,
		Notes: "Freight billing for Q-ACME-100001",
		CreatedBy: &admin.ID, CreatedAt: now, UpdatedAt: now,
		FinalizedAt: &now, Version: 4,
		Lines: []domain.BillingLine{
			{ID: "demo-bl-001", Description: "Ocean freight 40ft", Quantity: 1, UnitPrice: 3500000, LineTotal: 3500000},
			{ID: "demo-bl-002", Description: "Origin handling", Quantity: 1, UnitPrice: 500000, LineTotal: 500000},
			{ID: "demo-bl-003", Description: "Insurance", Quantity: 1, UnitPrice: 250000, LineTotal: 250000},
		},
	}
	billing.Seed(billingRec)

	// Seed a client payment
	payment := &domain.ClientPayment{
		ID: domain.ClientPaymentID("demo-pay-001"), ClientID: acme.ID,
		PaymentNumber: "PMT-ACME-100001", AmountCents: 2000000, CurrencyCode: "USD",
		PaymentMethod: "wire", ReferenceNumber: "WIRE-2026-002",
		ReceivedAt: now, CreatedBy: &admin.ID, CreatedAt: now, Version: 1,
	}
	payments.Create(nil, payment)

	// Build services
	quotSvc := service.NewQuotationService(clients, quotations, audit, log)
	wfSvc := service.NewWorkflowService(
		quotations, clients, fundingSrc, budgetReqs, approvals,
		disbursements, payProofs, liquidations, liqEvidence,
		billing, creditMemos, payments, allocations, audit, log,
	)

	return BootstrapResult{
		AuthService:     authSvc,
		QuotationService: quotSvc,
		WorkflowService: wfSvc,
		OTPCfg:          otpCfg,
	}
}
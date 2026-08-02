// Package demo wires the in-memory stores and seeds the demo data for local
// demo mode. This avoids needing PostgreSQL to run the demo.
package demo

import (
	"context"
	"log/slog"
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/auth"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/service"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/store"
)

// BootstrapResult holds the seeded services for demo mode.
type BootstrapResult struct {
	AuthService      *service.AuthService
	QuotationService *service.QuotationService
	WorkflowService  *service.WorkflowService
	OTPCfg           auth.OTPConfig
	ResetFn          func()
}

// Bootstrap creates in-memory stores, seeds realistic demo data, and returns
// ready-to-use services. The ResetFn re-seeds the data from scratch.
func Bootstrap(log *slog.Logger) BootstrapResult {
	stores := newDemoStores()
	services := buildServices(stores, log)
	seedData(stores, log)

	return BootstrapResult{
		AuthService:      services.authService,
		QuotationService: services.quotationService,
		WorkflowService:  services.workflowService,
		OTPCfg:           services.otpCfg,
		ResetFn: func() {
			log.Info("demo reset triggered, re-seeding data")
			seedData(stores, log)
		},
	}
}

type demoStores struct {
	users        *store.MemUserStore
	sessions     *store.MemSessionStore
	challenges   *store.MemEmailChallengeStore
	audit        *store.MemAuditStore
	clients      *store.MemClientStore
	quotations   *store.MemQuotationStore
	fundingSrc   *store.MemFundingSourceStore
	budgetReqs   *store.MemBudgetRequestStore
	approvals    *store.MemApprovalDecisionStore
	disbursements *store.MemDisbursementStore
	payProofs    *store.MemPaymentProofStore
	liquidations *store.MemLiquidationStore
	liqEvidence  *store.MemLiquidationEvidenceStore
	billing      *store.MemBillingRecordStore
	creditMemos  *store.MemCreditMemoStore
	payments     *store.MemClientPaymentStore
	allocations  *store.MemBillingAllocationStore
}

type demoServices struct {
	authService      *service.AuthService
	quotationService *service.QuotationService
	workflowService  *service.WorkflowService
	otpCfg           auth.OTPConfig
}

func newDemoStores() *demoStores {
	return &demoStores{
		users:         store.NewMemUserStore(),
		sessions:      store.NewMemSessionStore(),
		challenges:    store.NewMemEmailChallengeStore(),
		audit:         store.NewMemAuditStore(),
		clients:       store.NewMemClientStore(),
		quotations:    store.NewMemQuotationStore(),
		fundingSrc:    store.NewMemFundingSourceStore(),
		budgetReqs:    store.NewMemBudgetRequestStore(),
		approvals:     store.NewMemApprovalDecisionStore(),
		disbursements: store.NewMemDisbursementStore(),
		payProofs:     store.NewMemPaymentProofStore(),
		liquidations:  store.NewMemLiquidationStore(),
		liqEvidence:   store.NewMemLiquidationEvidenceStore(),
		billing:       store.NewMemBillingRecordStore(),
		creditMemos:   store.NewMemCreditMemoStore(),
		payments:      store.NewMemClientPaymentStore(),
		allocations:   store.NewMemBillingAllocationStore(),
	}
}

func buildServices(s *demoStores, log *slog.Logger) *demoServices {
	otpCfg := auth.DefaultOTPConfig([]byte("demo-otp-secret-32-bytes-minimum-ok"))

	authSvc := service.NewAuthService(
		s.users, s.sessions, s.challenges, s.audit,
		otpCfg, 24*time.Hour, 1*time.Hour, true, log,
	)

	quotSvc := service.NewQuotationService(s.clients, s.quotations, s.audit, log)

	wfSvc := service.NewWorkflowService(
		s.quotations, s.clients, s.fundingSrc, s.budgetReqs, s.approvals,
		s.disbursements, s.payProofs, s.liquidations, s.liqEvidence,
		s.billing, s.creditMemos, s.payments, s.allocations, s.audit, log,
	)

	return &demoServices{
		authService:      authSvc,
		quotationService: quotSvc,
		workflowService:  wfSvc,
		otpCfg:           otpCfg,
	}
}

// seedData populates all in-memory stores with realistic demo records.
// All emails use .example TLD. All amounts are in cents.
func seedData(s *demoStores, log *slog.Logger) {
	now := time.Now().UTC()
	adminID := domain.UserID("demo-admin-001")

	// --- User ---
	hash, _ := auth.HashPassword("Password123!")
	admin := &domain.User{
		ID: adminID, Username: "admin", Email: "admin@dosfreightflow.example",
		PasswordHash: hash, DisplayName: "Demo Administrator",
		EmailVerified: true, IsActive: true, CreatedAt: now, UpdatedAt: now,
	}
	s.users.Seed(admin)

	// --- Clients (5 realistic logistics companies) ---
	clients := []*domain.Client{
		{ID: "demo-client-acme", Name: "Acme Freight Logistics Inc.", Code: "ACME",
			ContactEmail: "operations@acme.example", ContactPhone: "+1-555-0100",
			Address: "123 Harbor Drive, Wilmington, DE 19801", IsActive: true, CreatedAt: now, UpdatedAt: now, Version: 1},
		{ID: "demo-client-globex", Name: "Globex Shipping Co.", Code: "GLOBX",
			ContactEmail: "finance@globex.example", ContactPhone: "+1-555-0200",
			Address: "456 Port Authority Blvd, Long Beach, CA 90802", IsActive: true, CreatedAt: now, UpdatedAt: now, Version: 1},
		{ID: "demo-client-stark", Name: "Stark Industries Transport", Code: "STARK",
			ContactEmail: "ap@stark.example", ContactPhone: "+1-555-0300",
			Address: "789 Industrial Parkway, Houston, TX 77001", IsActive: true, CreatedAt: now, UpdatedAt: now, Version: 1},
		{ID: "demo-client-initech", Name: "Initech Global Freight", Code: "INTEC",
			ContactEmail: "logistics@initech.example", ContactPhone: "+1-555-0400",
			Address: "321 Commerce Way, Atlanta, GA 30303", IsActive: true, CreatedAt: now, UpdatedAt: now, Version: 1},
		{ID: "demo-client-umbrella", Name: "Umbrella Cargo Services", Code: "UMBRA",
			ContactEmail: "billing@umbrella.example", ContactPhone: "+1-555-0500",
			Address: "654 Trade Center Dr, Chicago, IL 60601", IsActive: true, CreatedAt: now, UpdatedAt: now, Version: 1},
	}
	for _, c := range clients {
		s.clients.Seed(c)
	}

	// --- Funding Sources (3) ---
	fundingSources := []*domain.FundingSource{
		{ID: "demo-fs-bank-a", Name: "Bank A Operating Capital", Code: "BANKA",
			IsApproved: true, BalanceCents: 10000000, CurrencyCode: "USD", CreatedAt: now, UpdatedAt: now, Version: 1},
		{ID: "demo-fs-bank-b", Name: "Metro Credit Line", Code: "METRO",
			IsApproved: true, BalanceCents: 5000000, CurrencyCode: "USD", CreatedAt: now, UpdatedAt: now, Version: 1},
		{ID: "demo-fs-pending", Name: "Pacific Reserve Fund", Code: "PACIF",
			IsApproved: false, BalanceCents: 2000000, CurrencyCode: "USD", CreatedAt: now, UpdatedAt: now, Version: 1},
	}
	for _, fs := range fundingSources {
		s.fundingSrc.Seed(fs)
	}

	// --- Quotations (8 at various stages) ---
	quotations := []*domain.Quotation{
		// Accepted - Acme
		{ID: "demo-quot-001", ClientID: "demo-client-acme", QuotationNumber: "Q-ACME-100001",
			Status: domain.QuotationStatusAccepted, CurrencyCode: "USD",
			Subtotal: 4250000, TaxAmount: 0, Total: 4250000,
			Notes: "Ocean freight: 40ft container, Manila to LA",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, AcceptedAt: &now, Version: 3,
			Lines: []domain.QuotationLine{
				{ID: "demo-ql-001", Description: "Ocean freight 40ft", Quantity: 1, UnitPrice: 3500000, LineTotal: 3500000},
				{ID: "demo-ql-002", Description: "Origin handling", Quantity: 1, UnitPrice: 500000, LineTotal: 500000},
				{ID: "demo-ql-003", Description: "Insurance", Quantity: 1, UnitPrice: 250000, LineTotal: 250000},
			}},
		// Draft - Globex
		{ID: "demo-quot-002", ClientID: "demo-client-globex", QuotationNumber: "Q-GLOBX-100002",
			Status: domain.QuotationStatusDraft, CurrencyCode: "USD",
			Subtotal: 1850000, TaxAmount: 0, Total: 1850000,
			Notes: "Air freight: JFK to Heathrow",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, Version: 1,
			Lines: []domain.QuotationLine{
				{ID: "demo-ql-004", Description: "Air freight express", Quantity: 1, UnitPrice: 1500000, LineTotal: 1500000},
				{ID: "demo-ql-005", Description: "Customs clearance", Quantity: 1, UnitPrice: 350000, LineTotal: 350000},
			}},
		// In review - Stark
		{ID: "demo-quot-003", ClientID: "demo-client-stark", QuotationNumber: "Q-STARK-100003",
			Status: domain.QuotationStatusReview, CurrencyCode: "USD",
			Subtotal: 3200000, TaxAmount: 0, Total: 3200000,
			Notes: "Rail freight: Houston to Chicago, bulk cargo",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, Version: 2,
			Lines: []domain.QuotationLine{
				{ID: "demo-ql-006", Description: "Rail freight bulk", Quantity: 3, UnitPrice: 900000, LineTotal: 2700000},
				{ID: "demo-ql-007", Description: "Loading and unloading", Quantity: 1, UnitPrice: 500000, LineTotal: 500000},
			}},
		// Approved - Initech
		{ID: "demo-quot-004", ClientID: "demo-client-initech", QuotationNumber: "Q-INTEC-100004",
			Status: domain.QuotationStatusApproved, CurrencyCode: "USD",
			Subtotal: 2750000, TaxAmount: 0, Total: 2750000,
			Notes: "Intermodal: Atlanta to Miami, refrigerated",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, Version: 3,
			Lines: []domain.QuotationLine{
				{ID: "demo-ql-008", Description: "Intermodal transport", Quantity: 2, UnitPrice: 1100000, LineTotal: 2200000},
				{ID: "demo-ql-009", Description: "Reefer fuel surcharge", Quantity: 1, UnitPrice: 550000, LineTotal: 550000},
			}},
		// Accepted - Umbrella
		{ID: "demo-quot-005", ClientID: "demo-client-umbrella", QuotationNumber: "Q-UMBRA-100005",
			Status: domain.QuotationStatusAccepted, CurrencyCode: "USD",
			Subtotal: 5100000, TaxAmount: 0, Total: 5100000,
			Notes: "Multimodal: Chicago to Toronto, oversize cargo",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, AcceptedAt: &now, Version: 3,
			Lines: []domain.QuotationLine{
				{ID: "demo-ql-010", Description: "Truck freight oversize", Quantity: 1, UnitPrice: 3800000, LineTotal: 3800000},
				{ID: "demo-ql-011", Description: "Permit and escort", Quantity: 1, UnitPrice: 800000, LineTotal: 800000},
				{ID: "demo-ql-012", Description: "Border crossing fees", Quantity: 1, UnitPrice: 500000, LineTotal: 500000},
			}},
		// Draft - Acme (second quote)
		{ID: "demo-quot-006", ClientID: "demo-client-acme", QuotationNumber: "Q-ACME-100006",
			Status: domain.QuotationStatusDraft, CurrencyCode: "USD",
			Subtotal: 980000, TaxAmount: 0, Total: 980000,
			Notes: "Drayage: LA port to inland warehouse",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, Version: 1,
			Lines: []domain.QuotationLine{
				{ID: "demo-ql-013", Description: "Drayage service", Quantity: 2, UnitPrice: 400000, LineTotal: 800000},
				{ID: "demo-ql-014", Description: "Fuel surcharge", Quantity: 1, UnitPrice: 180000, LineTotal: 180000},
			}},
		// Revised - Globex
		{ID: "demo-quot-007", ClientID: "demo-client-globex", QuotationNumber: "Q-GLOBX-100007",
			Status: domain.QuotationStatusRevised, CurrencyCode: "USD",
			Subtotal: 2150000, TaxAmount: 0, Total: 2150000,
			Notes: "Sea freight: Shanghai to Hamburg, revised rates",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, Version: 4,
			Lines: []domain.QuotationLine{
				{ID: "demo-ql-015", Description: "Ocean freight 40ft", Quantity: 1, UnitPrice: 1700000, LineTotal: 1700000},
				{ID: "demo-ql-016", Description: "Terminal handling", Quantity: 1, UnitPrice: 450000, LineTotal: 450000},
			}},
		// Void - Stark
		{ID: "demo-quot-008", ClientID: "demo-client-stark", QuotationNumber: "Q-STARK-100008",
			Status: domain.QuotationStatusVoid, CurrencyCode: "USD",
			Subtotal: 1500000, TaxAmount: 0, Total: 1500000,
			Notes: "Voided: client cancelled shipment",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, VoidedAt: &now, Version: 2,
			Lines: []domain.QuotationLine{
				{ID: "demo-ql-017", Description: "Air freight", Quantity: 1, UnitPrice: 1200000, LineTotal: 1200000},
				{ID: "demo-ql-018", Description: "Handling", Quantity: 1, UnitPrice: 300000, LineTotal: 300000},
			}},
	}
	for _, q := range quotations {
		s.quotations.Seed(q)
	}

	// --- Budget Requests (4 at various stages) ---
	budgetReqs := []*domain.BudgetRequest{
		{ID: "demo-br-001", QuotationID: "demo-quot-001", ClientID: "demo-client-acme",
			RequestNumber: "BR-ACME-100001", Status: domain.BudgetRequestStatusApproved,
			CurrencyCode: "USD", AmountCents: 4250000, Purpose: "Ocean freight funding for Manila to LA",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, ApprovedAt: &now, Version: 3},
		{ID: "demo-br-002", QuotationID: "demo-quot-005", ClientID: "demo-client-umbrella",
			RequestNumber: "BR-UMBRA-100002", Status: domain.BudgetRequestStatusApproved,
			CurrencyCode: "USD", AmountCents: 5100000, Purpose: "Oversize cargo transport to Toronto",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, ApprovedAt: &now, Version: 3},
		{ID: "demo-br-003", QuotationID: "demo-quot-004", ClientID: "demo-client-initech",
			RequestNumber: "BR-INTEC-100003", Status: domain.BudgetRequestStatusControlsReview,
			CurrencyCode: "USD", AmountCents: 2750000, Purpose: "Refrigerated intermodal Atlanta to Miami",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, Version: 2},
		{ID: "demo-br-004", QuotationID: "demo-quot-003", ClientID: "demo-client-stark",
			RequestNumber: "BR-STARK-100004", Status: domain.BudgetRequestStatusDraft,
			CurrencyCode: "USD", AmountCents: 3200000, Purpose: "Rail freight bulk cargo Houston to Chicago",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, Version: 1},
	}
	for _, br := range budgetReqs {
		s.budgetReqs.Seed(br)
	}

	// --- Disbursements (3) ---
	disbursements := []*domain.Disbursement{
		{ID: "demo-disb-001", BudgetRequestID: "demo-br-001", FundingSourceID: "demo-fs-bank-a",
			Status: domain.DisbursementStatusReleased, AmountCents: 4250000, CurrencyCode: "USD",
			ReferenceNumber: "WIRE-2026-001", Notes: "Initial disbursement for Manila to LA",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, ReleasedAt: &now, Version: 2},
		{ID: "demo-disb-002", BudgetRequestID: "demo-br-002", FundingSourceID: "demo-fs-bank-b",
			Status: domain.DisbursementStatusReleased, AmountCents: 5100000, CurrencyCode: "USD",
			ReferenceNumber: "WIRE-2026-002", Notes: "Oversize cargo disbursement",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, ReleasedAt: &now, Version: 2},
		{ID: "demo-disb-003", BudgetRequestID: "demo-br-001", FundingSourceID: "demo-fs-bank-a",
			Status: domain.DisbursementStatusPending, AmountCents: 500000, CurrencyCode: "USD",
			ReferenceNumber: "", Notes: "Additional insurance top-up",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, Version: 1},
	}
	for _, d := range disbursements {
		s.disbursements.Seed(d)
	}

	// --- Liquidations (2) ---
	liquidations := []*domain.Liquidation{
		{ID: "demo-liq-001", DisbursementID: "demo-disb-001",
			Status: domain.LiquidationStatusReconciled,
			ReleasedAmount: 4250000, ActualAmount: 3980000, VarianceAmount: 270000,
			Notes: "Under budget by $2,700 due to negotiated rates",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, Version: 2},
		{ID: "demo-liq-002", DisbursementID: "demo-disb-002",
			Status: domain.LiquidationStatusOpen,
			ReleasedAmount: 5100000, ActualAmount: 0, VarianceAmount: 5100000,
			Notes: "Awaiting actual spending report",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, Version: 1},
	}
	for _, l := range liquidations {
		s.liquidations.Seed(l)
	}

	// --- Billing Records (4) ---
	billingRecords := []*domain.BillingRecord{
		{ID: "demo-bill-001", ClientID: "demo-client-acme", BudgetRequestID: ptrBR("demo-br-001"),
			BillingNumber: "INV-ACME-100001", Status: domain.BillingStatusFinalized,
			CurrencyCode: "USD", Subtotal: 4250000, TaxAmount: 0, Total: 4250000,
			Notes: "Freight billing for Q-ACME-100001",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, FinalizedAt: &now, Version: 4,
			Lines: []domain.BillingLine{
				{ID: "demo-bl-001", Description: "Ocean freight 40ft", Quantity: 1, UnitPrice: 3500000, LineTotal: 3500000},
				{ID: "demo-bl-002", Description: "Origin handling", Quantity: 1, UnitPrice: 500000, LineTotal: 500000},
				{ID: "demo-bl-003", Description: "Insurance", Quantity: 1, UnitPrice: 250000, LineTotal: 250000},
			}},
		{ID: "demo-bill-002", ClientID: "demo-client-umbrella", BudgetRequestID: ptrBR("demo-br-002"),
			BillingNumber: "INV-UMBRA-100002", Status: domain.BillingStatusFinalized,
			CurrencyCode: "USD", Subtotal: 5100000, TaxAmount: 0, Total: 5100000,
			Notes: "Oversize cargo billing for Q-UMBRA-100005",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, FinalizedAt: &now, Version: 4,
			Lines: []domain.BillingLine{
				{ID: "demo-bl-004", Description: "Truck freight oversize", Quantity: 1, UnitPrice: 3800000, LineTotal: 3800000},
				{ID: "demo-bl-005", Description: "Permit and escort", Quantity: 1, UnitPrice: 800000, LineTotal: 800000},
				{ID: "demo-bl-006", Description: "Border crossing fees", Quantity: 1, UnitPrice: 500000, LineTotal: 500000},
			}},
		{ID: "demo-bill-003", ClientID: "demo-client-globex",
			BillingNumber: "INV-GLOBX-100003", Status: domain.BillingStatusApproved,
			CurrencyCode: "USD", Subtotal: 1850000, TaxAmount: 0, Total: 1850000,
			Notes: "Air freight billing for Q-GLOBX-100002",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, ApprovedAt: &now, Version: 3,
			Lines: []domain.BillingLine{
				{ID: "demo-bl-007", Description: "Air freight express", Quantity: 1, UnitPrice: 1500000, LineTotal: 1500000},
				{ID: "demo-bl-008", Description: "Customs clearance", Quantity: 1, UnitPrice: 350000, LineTotal: 350000},
			}},
		{ID: "demo-bill-004", ClientID: "demo-client-stark",
			BillingNumber: "INV-STARK-100004", Status: domain.BillingStatusReview,
			CurrencyCode: "USD", Subtotal: 3200000, TaxAmount: 0, Total: 3200000,
			Notes: "Rail freight billing for Q-STARK-100003",
			CreatedBy: &adminID, CreatedAt: now, UpdatedAt: now, Version: 2,
			Lines: []domain.BillingLine{
				{ID: "demo-bl-009", Description: "Rail freight bulk", Quantity: 3, UnitPrice: 900000, LineTotal: 2700000},
				{ID: "demo-bl-010", Description: "Loading and unloading", Quantity: 1, UnitPrice: 500000, LineTotal: 500000},
			}},
	}
	for _, b := range billingRecords {
		s.billing.Seed(b)
	}

	// --- Client Payments (3) ---
	payments := []*domain.ClientPayment{
		{ID: "demo-pay-001", ClientID: "demo-client-acme", PaymentNumber: "PMT-ACME-100001",
			AmountCents: 2000000, CurrencyCode: "USD", PaymentMethod: "wire", ReferenceNumber: "WIRE-2026-101",
			ReceivedAt: now, CreatedBy: &adminID, CreatedAt: now, Version: 1},
		{ID: "demo-pay-002", ClientID: "demo-client-acme", PaymentNumber: "PMT-ACME-100002",
			AmountCents: 2250000, CurrencyCode: "USD", PaymentMethod: "ach", ReferenceNumber: "ACH-2026-051",
			ReceivedAt: now, CreatedBy: &adminID, CreatedAt: now, Version: 1},
		{ID: "demo-pay-003", ClientID: "demo-client-umbrella", PaymentNumber: "PMT-UMBRA-100003",
			AmountCents: 3100000, CurrencyCode: "USD", PaymentMethod: "wire", ReferenceNumber: "WIRE-2026-102",
			ReceivedAt: now, CreatedBy: &adminID, CreatedAt: now, Version: 1},
	}
	for _, p := range payments {
		s.payments.Create(context.Background(), p)
	}

	log.Info("demo data seeded",
		"clients", 5, "funding_sources", 3, "quotations", 8,
		"budget_requests", 4, "disbursements", 3, "liquidations", 2,
		"billing_records", 4, "payments", 3,
	)
}

func ptrBR(id string) *domain.BudgetRequestID {
	bid := domain.BudgetRequestID(id)
	return &bid
}
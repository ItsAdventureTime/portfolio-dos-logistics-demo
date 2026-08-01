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
	ErrBudgetRequestNotFound   = errors.New("budget request not found")
	ErrDisbursementNotFound     = errors.New("disbursement not found")
	ErrLiquidationNotFound      = errors.New("liquidation not found")
	ErrBillingNotFound          = errors.New("billing record not found")
	ErrPaymentNotFound          = errors.New("client payment not found")
	ErrFundingSourceNotApproved = errors.New("funding source is not approved")
	ErrBillingImmutable         = errors.New("billing record is finalized and cannot be edited in place")
	ErrAllocationClientMismatch = errors.New("billing record does not belong to the payment's client")
)

// WorkflowService handles the funding→approval→disbursement→liquidation→
// billing→collection workflow slices.
type WorkflowService struct {
	quotations   repository.QuotationRepository
	clients      repository.ClientRepository
	fundingSrc   repository.FundingSourceRepository
	budgetReqs   repository.BudgetRequestRepository
	approvals    repository.ApprovalDecisionRepository
	disbursements repository.DisbursementRepository
	payProofs    repository.PaymentProofRepository
	liquidations repository.LiquidationRepository
	liqEvidence  repository.LiquidationEvidenceRepository
	billing      repository.BillingRecordRepository
	creditMemos  repository.CreditMemoRepository
	payments     repository.ClientPaymentRepository
	allocations  repository.BillingAllocationRepository
	audit        repository.AuditEventRepository
	log          *slog.Logger
}

func NewWorkflowService(
	quotations repository.QuotationRepository,
	clients repository.ClientRepository,
	fundingSrc repository.FundingSourceRepository,
	budgetReqs repository.BudgetRequestRepository,
	approvals repository.ApprovalDecisionRepository,
	disbursements repository.DisbursementRepository,
	payProofs repository.PaymentProofRepository,
	liquidations repository.LiquidationRepository,
	liqEvidence repository.LiquidationEvidenceRepository,
	billing repository.BillingRecordRepository,
	creditMemos repository.CreditMemoRepository,
	payments repository.ClientPaymentRepository,
	allocations repository.BillingAllocationRepository,
	audit repository.AuditEventRepository,
	log *slog.Logger,
) *WorkflowService {
	return &WorkflowService{
		quotations: quotations, clients: clients, fundingSrc: fundingSrc,
		budgetReqs: budgetReqs, approvals: approvals, disbursements: disbursements,
		payProofs: payProofs, liquidations: liquidations, liqEvidence: liqEvidence,
		billing: billing, creditMemos: creditMemos, payments: payments,
		allocations: allocations, audit: audit, log: log,
	}
}

// --- Slice 2: Budget Request ---

// CreateBudgetRequest creates a funding request linked to an accepted quotation.
// Per docs/spec/02: "A budget request must link to an accepted quotation and a client."
func (s *WorkflowService) CreateBudgetRequest(ctx context.Context, actor domain.UserID, quotationID domain.QuotationID, clientID domain.ClientID, currencyCode, purpose string, amountCents int64) (*domain.BudgetRequest, error) {
	corrID := observability.CorrelationFrom(ctx)

	// Verify quotation exists and is accepted.
	q, err := s.quotations.GetByID(ctx, quotationID)
	if err != nil || q == nil {
		return nil, ErrQuotationNotFound
	}
	if q.Status != domain.QuotationStatusAccepted {
		return nil, ErrQuotationNotAccepted
	}

	// Verify client exists.
	client, err := s.clients.GetByID(ctx, clientID)
	if err != nil || client == nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	// Verify the quotation belongs to this client.
	if q.ClientID != clientID {
		return nil, fmt.Errorf("quotation does not belong to this client")
	}

	br := &domain.BudgetRequest{
		ID:            domain.BudgetRequestID(uuid.NewString()),
		QuotationID:   quotationID,
		ClientID:      clientID,
		RequestNumber: fmt.Sprintf("BR-%s-%d", client.Code, time.Now().Unix()),
		Status:        domain.BudgetRequestStatusDraft,
		CurrencyCode:  currencyCode,
		AmountCents:   amountCents,
		Purpose:       purpose,
		CreatedBy:     &actor,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		Version:       1,
	}

	if err := s.budgetReqs.Create(ctx, br); err != nil {
		return nil, fmt.Errorf("create budget request: %w", err)
	}

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID, ActorUserID: &actor, ActorRole: "administrator",
		Action: domain.AuditActionBudgetRequestCreated, EntityType: "budget_request",
		EntityID: string(br.ID), Details: map[string]any{"request_number": br.RequestNumber},
	})

	return br, nil
}

// TransitionBudgetRequest transitions a budget request with optimistic concurrency.
func (s *WorkflowService) TransitionBudgetRequest(ctx context.Context, id domain.BudgetRequestID, target domain.BudgetRequestStatus, expectedVersion int, actor domain.UserID, reason string) (*domain.BudgetRequest, error) {
	corrID := observability.CorrelationFrom(ctx)

	br, err := s.budgetReqs.GetByID(ctx, id)
	if err != nil || br == nil {
		return nil, ErrBudgetRequestNotFound
	}

	if !domain.CanTransitionBudgetRequest(br.Status, target) {
		return nil, fmt.Errorf("%w: %s → %s", ErrInvalidTransition, br.Status, target)
	}

	updated, err := s.budgetReqs.UpdateStatus(ctx, id, target, expectedVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVersionConflict, err)
	}

	now := time.Now().UTC()
	switch target {
	case domain.BudgetRequestStatusApproved:
		updated.ApprovedAt = &now
	case domain.BudgetRequestStatusRejected:
		updated.RejectedAt = &now
	case domain.BudgetRequestStatusReturned:
		updated.ReturnedAt = &now
	}

	// Record approval decision (controls review).
	decisionType := domain.ApprovalDecisionApproved
	switch target {
	case domain.BudgetRequestStatusRejected:
		decisionType = domain.ApprovalDecisionRejected
	case domain.BudgetRequestStatusReturned:
		decisionType = domain.ApprovalDecisionReturned
	}
	_ = s.approvals.Create(ctx, &domain.ApprovalDecision{
		ID: uuid.NewString(), BudgetRequestID: id, Decision: decisionType,
		ActorID: actor, Reason: reason, CreatedAt: now,
	})

	action := domain.AuditActionBudgetRequestCreated
	switch target {
	case domain.BudgetRequestStatusControlsReview:
		action = domain.AuditActionBudgetRequestSubmitted
	case domain.BudgetRequestStatusApproved:
		action = domain.AuditActionBudgetRequestApproved
	case domain.BudgetRequestStatusRejected:
		action = domain.AuditActionBudgetRequestRejected
	case domain.BudgetRequestStatusReturned:
		action = domain.AuditActionBudgetRequestReturned
	}

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID, ActorUserID: &actor, ActorRole: "administrator",
		Action: action, EntityType: "budget_request", EntityID: string(id),
		Details: map[string]any{"from_status": string(br.Status), "to_status": string(target), "reason": reason},
	})

	return updated, nil
}

// --- Slice 4: Disbursement ---

// CreateDisbursement creates a disbursement linked to an approved budget request
// and an approved funding source.
// Per docs/spec/02: "A disbursement must reference an approved request and an approved funding source."
func (s *WorkflowService) CreateDisbursement(ctx context.Context, actor domain.UserID, budgetReqID domain.BudgetRequestID, fundingSourceID domain.FundingSourceID, amountCents int64, currencyCode, refNum, notes string) (*domain.Disbursement, error) {
	corrID := observability.CorrelationFrom(ctx)

	br, err := s.budgetReqs.GetByID(ctx, budgetReqID)
	if err != nil || br == nil {
		return nil, ErrBudgetRequestNotFound
	}
	if br.Status != domain.BudgetRequestStatusApproved {
		return nil, fmt.Errorf("budget request must be approved before disbursement")
	}

	fs, err := s.fundingSrc.GetByID(ctx, fundingSourceID)
	if err != nil || fs == nil {
		return nil, fmt.Errorf("funding source not found: %w", err)
	}
	if !fs.IsApproved {
		return nil, ErrFundingSourceNotApproved
	}

	d := &domain.Disbursement{
		ID:             domain.DisbursementID(uuid.NewString()),
		BudgetRequestID: budgetReqID,
		FundingSourceID: fundingSourceID,
		Status:          domain.DisbursementStatusPending,
		AmountCents:    amountCents,
		CurrencyCode:    currencyCode,
		ReferenceNumber: refNum,
		Notes:           notes,
		CreatedBy:       &actor,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		Version:         1,
	}

	if err := s.disbursements.Create(ctx, d); err != nil {
		return nil, fmt.Errorf("create disbursement: %w", err)
	}

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID, ActorUserID: &actor, ActorRole: "administrator",
		Action: domain.AuditActionDisbursementCreated, EntityType: "disbursement",
		EntityID: string(d.ID), Details: map[string]any{"amount": amountCents},
	})

	return d, nil
}

// TransitionDisbursement transitions a disbursement with optimistic concurrency.
func (s *WorkflowService) TransitionDisbursement(ctx context.Context, id domain.DisbursementID, target domain.DisbursementStatus, expectedVersion int, actor domain.UserID) (*domain.Disbursement, error) {
	corrID := observability.CorrelationFrom(ctx)

	d, err := s.disbursements.GetByID(ctx, id)
	if err != nil || d == nil {
		return nil, ErrDisbursementNotFound
	}

	if !domain.CanTransitionDisbursement(d.Status, target) {
		return nil, fmt.Errorf("%w: %s → %s", ErrInvalidTransition, d.Status, target)
	}

	updated, err := s.disbursements.UpdateStatus(ctx, id, target, expectedVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVersionConflict, err)
	}

	now := time.Now().UTC()
	switch target {
	case domain.DisbursementStatusReleased:
		updated.ReleasedAt = &now
	case domain.DisbursementStatusHeld:
		updated.HeldAt = &now
	case domain.DisbursementStatusReturned:
		updated.ReturnedAt = &now
	}

	action := domain.AuditActionDisbursementCreated
	switch target {
	case domain.DisbursementStatusReleased:
		action = domain.AuditActionDisbursementReleased
	case domain.DisbursementStatusHeld:
		action = domain.AuditActionDisbursementHeld
	case domain.DisbursementStatusReturned:
		action = domain.AuditActionDisbursementReturned
	}

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID, ActorUserID: &actor, ActorRole: "administrator",
		Action: action, EntityType: "disbursement", EntityID: string(id),
		Details: map[string]any{"from_status": string(d.Status), "to_status": string(target)},
	})

	return updated, nil
}

// --- Slice 5: Liquidation ---

// CreateLiquidation creates a liquidation for a disbursement.
// Per docs/spec/02: "Liquidation must reconcile released funds with actual spending
// and retain variance evidence."
func (s *WorkflowService) CreateLiquidation(ctx context.Context, actor domain.UserID, disbursementID domain.DisbursementID, releasedAmount int64) (*domain.Liquidation, error) {
	corrID := observability.CorrelationFrom(ctx)

	d, err := s.disbursements.GetByID(ctx, disbursementID)
	if err != nil || d == nil {
		return nil, ErrDisbursementNotFound
	}
	if d.Status != domain.DisbursementStatusReleased {
		return nil, fmt.Errorf("disbursement must be released before liquidation")
	}

	l := &domain.Liquidation{
		ID:             domain.LiquidationID(uuid.NewString()),
		DisbursementID: disbursementID,
		Status:         domain.LiquidationStatusOpen,
		ReleasedAmount: releasedAmount,
		ActualAmount:   0,
		VarianceAmount: releasedAmount, // initial variance = released - 0
		CreatedBy:      &actor,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Version:        1,
	}

	if err := s.liquidations.Create(ctx, l); err != nil {
		return nil, fmt.Errorf("create liquidation: %w", err)
	}

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID, ActorUserID: &actor, ActorRole: "administrator",
		Action: domain.AuditActionLiquidationCreated, EntityType: "liquidation",
		EntityID: string(l.ID), Details: map[string]any{"released_amount": releasedAmount},
	})

	return l, nil
}

// ReconcileLiquidation records the actual spending and calculates variance.
func (s *WorkflowService) ReconcileLiquidation(ctx context.Context, id domain.LiquidationID, actualAmount int64, expectedVersion int, actor domain.UserID) (*domain.Liquidation, error) {
	corrID := observability.CorrelationFrom(ctx)

	l, err := s.liquidations.GetByID(ctx, id)
	if err != nil || l == nil {
		return nil, ErrLiquidationNotFound
	}

	if !domain.CanTransitionLiquidation(l.Status, domain.LiquidationStatusReconciled) {
		return nil, fmt.Errorf("%w: %s → reconciled", ErrInvalidTransition, l.Status)
	}

	variance := l.ReleasedAmount - actualAmount
	updated, err := s.liquidations.UpdateStatus(ctx, id, domain.LiquidationStatusReconciled, actualAmount, variance, expectedVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVersionConflict, err)
	}

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID, ActorUserID: &actor, ActorRole: "administrator",
		Action: domain.AuditActionLiquidationReconciled, EntityType: "liquidation",
		EntityID: string(id), Details: map[string]any{
			"released": l.ReleasedAmount, "actual": actualAmount, "variance": variance,
		},
	})

	return updated, nil
}

// CloseLiquidation closes a reconciled liquidation.
func (s *WorkflowService) CloseLiquidation(ctx context.Context, id domain.LiquidationID, expectedVersion int, actor domain.UserID) (*domain.Liquidation, error) {
	corrID := observability.CorrelationFrom(ctx)

	l, err := s.liquidations.GetByID(ctx, id)
	if err != nil || l == nil {
		return nil, ErrLiquidationNotFound
	}

	if !domain.CanTransitionLiquidation(l.Status, domain.LiquidationStatusClosed) {
		return nil, fmt.Errorf("%w: %s → closed", ErrInvalidTransition, l.Status)
	}

	updated, err := s.liquidations.UpdateStatus(ctx, id, domain.LiquidationStatusClosed, l.ActualAmount, l.VarianceAmount, expectedVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVersionConflict, err)
	}

	now := time.Now().UTC()
	updated.ClosedAt = &now

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID, ActorUserID: &actor, ActorRole: "administrator",
		Action: domain.AuditActionLiquidationClosed, EntityType: "liquidation",
		EntityID: string(id), Details: map[string]any{},
	})

	return updated, nil
}

// --- Slice 6: Billing ---

// CreateBillingRecord creates a new billing record in draft status.
// Per docs/spec/02: "Billing must follow the required approval path before finalization."
func (s *WorkflowService) CreateBillingRecord(ctx context.Context, actor domain.UserID, clientID domain.ClientID, currencyCode string, lines []BillingLineInput, budgetReqID *domain.BudgetRequestID) (*domain.BillingRecord, error) {
	corrID := observability.CorrelationFrom(ctx)

	if len(lines) == 0 {
		return nil, ErrQuotationNoLines // reuse the "at least one charge" invariant
	}

	client, err := s.clients.GetByID(ctx, clientID)
	if err != nil || client == nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	billingLines := make([]domain.BillingLine, len(lines))
	var subtotal int64
	for i, li := range lines {
		lineTotal := li.Quantity * li.UnitPrice
		billingLines[i] = domain.BillingLine{
			ID: uuid.NewString(), Description: li.Description,
			Quantity: li.Quantity, UnitPrice: li.UnitPrice,
			LineTotal: lineTotal, SortOrder: i, CreatedAt: time.Now().UTC(),
		}
		subtotal += lineTotal
	}

	b := &domain.BillingRecord{
		ID:            domain.BillingRecordID(uuid.NewString()),
		ClientID:      clientID,
		BudgetRequestID: budgetReqID,
		BillingNumber: fmt.Sprintf("INV-%s-%d", client.Code, time.Now().Unix()),
		Status:        domain.BillingStatusDraft,
		CurrencyCode:  currencyCode,
		Subtotal:      subtotal,
		Total:         subtotal,
		CreatedBy:     &actor,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		Version:       1,
		Lines:         billingLines,
	}

	if err := s.billing.Create(ctx, b); err != nil {
		return nil, fmt.Errorf("create billing record: %w", err)
	}

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID, ActorUserID: &actor, ActorRole: "administrator",
		Action: domain.AuditActionBillingCreated, EntityType: "billing_record",
		EntityID: string(b.ID), Details: map[string]any{"billing_number": b.BillingNumber},
	})

	return b, nil
}

// TransitionBilling transitions a billing record with optimistic concurrency.
// Per docs/spec/02: "A finalized financial record must not be edited in place;
// use a controlled replacement or void process."
func (s *WorkflowService) TransitionBilling(ctx context.Context, id domain.BillingRecordID, target domain.BillingStatus, expectedVersion int, actor domain.UserID) (*domain.BillingRecord, error) {
	corrID := observability.CorrelationFrom(ctx)

	b, err := s.billing.GetByID(ctx, id)
	if err != nil || b == nil {
		return nil, ErrBillingNotFound
	}

	if !domain.CanTransitionBilling(b.Status, target) {
		return nil, fmt.Errorf("%w: %s → %s", ErrInvalidTransition, b.Status, target)
	}

	updated, err := s.billing.UpdateStatus(ctx, id, target, expectedVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVersionConflict, err)
	}

	now := time.Now().UTC()
	switch target {
	case domain.BillingStatusApproved:
		updated.ApprovedAt = &now
	case domain.BillingStatusFinalized:
		updated.FinalizedAt = &now
	case domain.BillingStatusVoid:
		updated.VoidedAt = &now
	}

	action := domain.AuditActionBillingCreated
	switch target {
	case domain.BillingStatusReview:
		action = domain.AuditActionBillingSubmitted
	case domain.BillingStatusApproved:
		action = domain.AuditActionBillingApproved
	case domain.BillingStatusFinalized:
		action = domain.AuditActionBillingFinalized
	case domain.BillingStatusVoid:
		action = domain.AuditActionBillingVoided
	case domain.BillingStatusReplaced:
		action = domain.AuditActionBillingReplaced
	}

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID, ActorUserID: &actor, ActorRole: "administrator",
		Action: action, EntityType: "billing_record", EntityID: string(id),
		Details: map[string]any{"from_status": string(b.Status), "to_status": string(target)},
	})

	return updated, nil
}

// ReplaceBilling creates a replacement billing record for a finalized one.
// Per docs/spec/02: "use a controlled replacement or void process."
func (s *WorkflowService) ReplaceBilling(ctx context.Context, actor domain.UserID, originalID domain.BillingRecordID, expectedVersion int, currencyCode string, lines []BillingLineInput) (*domain.BillingRecord, error) {
	original, err := s.billing.GetByID(ctx, originalID)
	if err != nil || original == nil {
		return nil, ErrBillingNotFound
	}
	if original.Status != domain.BillingStatusFinalized {
		return nil, ErrBillingImmutable
	}

	// Create the replacement.
	replacement, err := s.CreateBillingRecord(ctx, actor, original.ClientID, currencyCode, lines, original.BudgetRequestID)
	if err != nil {
		return nil, err
	}

	// Mark the original as replaced.
	if err := s.billing.MarkReplaced(ctx, originalID, replacement.ID); err != nil {
		return nil, fmt.Errorf("mark replaced: %w", err)
	}
	_, err = s.billing.UpdateStatus(ctx, originalID, domain.BillingStatusReplaced, expectedVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVersionConflict, err)
	}

	return replacement, nil
}

// --- Slice 7: Collection ---

// RecordClientPayment records a payment from a client.
func (s *WorkflowService) RecordClientPayment(ctx context.Context, actor domain.UserID, clientID domain.ClientID, amountCents int64, currencyCode, method, refNum string) (*domain.ClientPayment, error) {
	corrID := observability.CorrelationFrom(ctx)

	client, err := s.clients.GetByID(ctx, clientID)
	if err != nil || client == nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	p := &domain.ClientPayment{
		ID:             domain.ClientPaymentID(uuid.NewString()),
		ClientID:       clientID,
		PaymentNumber:  fmt.Sprintf("PMT-%s-%d", client.Code, time.Now().Unix()),
		AmountCents:    amountCents,
		CurrencyCode:   currencyCode,
		PaymentMethod:  method,
		ReferenceNumber: refNum,
		ReceivedAt:     time.Now().UTC(),
		CreatedBy:      &actor,
		CreatedAt:      time.Now().UTC(),
		Version:        1,
	}

	if err := s.payments.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID, ActorUserID: &actor, ActorRole: "administrator",
		Action: domain.AuditActionClientPaymentReceived, EntityType: "client_payment",
		EntityID: string(p.ID), Details: map[string]any{"amount": amountCents},
	})

	return p, nil
}

// AllocatePayment allocates a payment to a billing record.
// Per docs/spec/02: "Client payments may allocate only to billing records for the selected client."
func (s *WorkflowService) AllocatePayment(ctx context.Context, actor domain.UserID, paymentID domain.ClientPaymentID, billingID domain.BillingRecordID, amountCents int64) (*domain.BillingAllocation, error) {
	corrID := observability.CorrelationFrom(ctx)

	payment, err := s.payments.GetByID(ctx, paymentID)
	if err != nil || payment == nil {
		return nil, ErrPaymentNotFound
	}

	billing, err := s.billing.GetByID(ctx, billingID)
	if err != nil || billing == nil {
		return nil, ErrBillingNotFound
	}

	// CRITICAL: billing record must belong to the payment's client.
	if billing.ClientID != payment.ClientID {
		return nil, ErrAllocationClientMismatch
	}

	// Billing must be finalized to receive allocations.
	if billing.Status != domain.BillingStatusFinalized {
		return nil, fmt.Errorf("billing record must be finalized to receive allocations")
	}

	a := &domain.BillingAllocation{
		ID:             uuid.NewString(),
		ClientPaymentID: paymentID,
		BillingRecordID: billingID,
		AmountCents:    amountCents,
		AllocatedAt:    time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
	}

	if err := s.allocations.Create(ctx, a); err != nil {
		return nil, fmt.Errorf("create allocation: %w", err)
	}

	_ = s.audit.Create(ctx, &domain.AuditEvent{
		CorrelationID: corrID, ActorUserID: &actor, ActorRole: "administrator",
		Action: domain.AuditActionBillingAllocationCreated, EntityType: "billing_allocation",
		EntityID: a.ID, Details: map[string]any{
			"payment_id": string(paymentID), "billing_id": string(billingID), "amount": amountCents,
		},
	})

	return a, nil
}

type BillingLineInput struct {
	Description string
	Quantity    int64
	UnitPrice   int64
}
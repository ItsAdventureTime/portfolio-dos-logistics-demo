package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
)

// --- Funding Source ---

type MemFundingSourceStore struct {
	mu      sync.RWMutex
	items   map[domain.FundingSourceID]*domain.FundingSource
}

func NewMemFundingSourceStore() *MemFundingSourceStore {
	return &MemFundingSourceStore{items: make(map[domain.FundingSourceID]*domain.FundingSource)}
}

func (s *MemFundingSourceStore) Seed(fs *domain.FundingSource) {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[fs.ID] = fs
}

func (s *MemFundingSourceStore) Create(ctx context.Context, fs *domain.FundingSource) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[fs.ID] = fs
	return nil
}

func (s *MemFundingSourceStore) GetByID(ctx context.Context, id domain.FundingSourceID) (*domain.FundingSource, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	fs, ok := s.items[id]
	if !ok { return nil, fmt.Errorf("funding source not found") }
	return fs, nil
}

func (s *MemFundingSourceStore) List(ctx context.Context) ([]*domain.FundingSource, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	out := make([]*domain.FundingSource, 0, len(s.items))
	for _, fs := range s.items { out = append(out, fs) }
	return out, nil
}

// --- Budget Request ---

type MemBudgetRequestStore struct {
	mu    sync.RWMutex
	items map[domain.BudgetRequestID]*domain.BudgetRequest
}

func NewMemBudgetRequestStore() *MemBudgetRequestStore {
	return &MemBudgetRequestStore{items: make(map[domain.BudgetRequestID]*domain.BudgetRequest)}
}

func (s *MemBudgetRequestStore) Seed(br *domain.BudgetRequest) {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[br.ID] = br
}

func (s *MemBudgetRequestStore) Create(ctx context.Context, br *domain.BudgetRequest) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[br.ID] = br
	return nil
}

func (s *MemBudgetRequestStore) GetByID(ctx context.Context, id domain.BudgetRequestID) (*domain.BudgetRequest, error) {
	return s.GetByIDWithLines(ctx, id)
}

func (s *MemBudgetRequestStore) GetByIDWithLines(ctx context.Context, id domain.BudgetRequestID) (*domain.BudgetRequest, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	br, ok := s.items[id]
	if !ok { return nil, fmt.Errorf("budget request not found") }
	return br, nil
}

func (s *MemBudgetRequestStore) List(ctx context.Context, clientID *domain.ClientID) ([]*domain.BudgetRequest, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	out := make([]*domain.BudgetRequest, 0)
	for _, br := range s.items {
		if clientID == nil || br.ClientID == *clientID {
			out = append(out, br)
		}
	}
	return out, nil
}

func (s *MemBudgetRequestStore) UpdateStatus(ctx context.Context, id domain.BudgetRequestID, status domain.BudgetRequestStatus, expectedVersion int) (*domain.BudgetRequest, error) {
	s.mu.Lock(); defer s.mu.Unlock()
	br, ok := s.items[id]
	if !ok { return nil, fmt.Errorf("budget request not found") }
	if br.Version != expectedVersion { return nil, ErrVersionConflict }
	br.Status = status
	br.Version++
	return br, nil
}

// --- Approval Decision ---

type MemApprovalDecisionStore struct {
	mu    sync.RWMutex
	items map[domain.BudgetRequestID][]*domain.ApprovalDecision
}

func NewMemApprovalDecisionStore() *MemApprovalDecisionStore {
	return &MemApprovalDecisionStore{items: make(map[domain.BudgetRequestID][]*domain.ApprovalDecision)}
}

func (s *MemApprovalDecisionStore) Create(ctx context.Context, d *domain.ApprovalDecision) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[d.BudgetRequestID] = append(s.items[d.BudgetRequestID], d)
	return nil
}

func (s *MemApprovalDecisionStore) ListByRequestID(ctx context.Context, reqID domain.BudgetRequestID) ([]*domain.ApprovalDecision, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.items[reqID], nil
}

// --- Disbursement ---

type MemDisbursementStore struct {
	mu    sync.RWMutex
	items map[domain.DisbursementID]*domain.Disbursement
}

func NewMemDisbursementStore() *MemDisbursementStore {
	return &MemDisbursementStore{items: make(map[domain.DisbursementID]*domain.Disbursement)}
}

func (s *MemDisbursementStore) Seed(d *domain.Disbursement) {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[d.ID] = d
}

func (s *MemDisbursementStore) Create(ctx context.Context, d *domain.Disbursement) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[d.ID] = d
	return nil
}

func (s *MemDisbursementStore) GetByID(ctx context.Context, id domain.DisbursementID) (*domain.Disbursement, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	d, ok := s.items[id]
	if !ok { return nil, fmt.Errorf("disbursement not found") }
	return d, nil
}

func (s *MemDisbursementStore) List(ctx context.Context, budgetRequestID *domain.BudgetRequestID) ([]*domain.Disbursement, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	out := make([]*domain.Disbursement, 0)
	for _, d := range s.items {
		if budgetRequestID == nil || d.BudgetRequestID == *budgetRequestID {
			out = append(out, d)
		}
	}
	return out, nil
}

func (s *MemDisbursementStore) UpdateStatus(ctx context.Context, id domain.DisbursementID, status domain.DisbursementStatus, expectedVersion int) (*domain.Disbursement, error) {
	s.mu.Lock(); defer s.mu.Unlock()
	d, ok := s.items[id]
	if !ok { return nil, fmt.Errorf("disbursement not found") }
	if d.Version != expectedVersion { return nil, ErrVersionConflict }
	d.Status = status
	d.Version++
	return d, nil
}

// --- Payment Proof ---

type MemPaymentProofStore struct {
	mu    sync.RWMutex
	items map[domain.DisbursementID][]*domain.PaymentProof
}

func NewMemPaymentProofStore() *MemPaymentProofStore {
	return &MemPaymentProofStore{items: make(map[domain.DisbursementID][]*domain.PaymentProof)}
}

func (s *MemPaymentProofStore) Create(ctx context.Context, p *domain.PaymentProof) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[p.DisbursementID] = append(s.items[p.DisbursementID], p)
	return nil
}

func (s *MemPaymentProofStore) ListByDisbursementID(ctx context.Context, dID domain.DisbursementID) ([]*domain.PaymentProof, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.items[dID], nil
}

// --- Liquidation ---

type MemLiquidationStore struct {
	mu        sync.RWMutex
	items      map[domain.LiquidationID]*domain.Liquidation
	byDisb    map[domain.DisbursementID]*domain.Liquidation
}

func NewMemLiquidationStore() *MemLiquidationStore {
	return &MemLiquidationStore{
		items:   make(map[domain.LiquidationID]*domain.Liquidation),
		byDisb:  make(map[domain.DisbursementID]*domain.Liquidation),
	}
}

func (s *MemLiquidationStore) Seed(l *domain.Liquidation) {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[l.ID] = l
	s.byDisb[l.DisbursementID] = l
}

func (s *MemLiquidationStore) Create(ctx context.Context, l *domain.Liquidation) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[l.ID] = l
	s.byDisb[l.DisbursementID] = l
	return nil
}

func (s *MemLiquidationStore) GetByID(ctx context.Context, id domain.LiquidationID) (*domain.Liquidation, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	l, ok := s.items[id]
	if !ok { return nil, fmt.Errorf("liquidation not found") }
	return l, nil
}

func (s *MemLiquidationStore) GetByDisbursementID(ctx context.Context, dID domain.DisbursementID) (*domain.Liquidation, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	l, ok := s.byDisb[dID]
	if !ok { return nil, fmt.Errorf("liquidation not found") }
	return l, nil
}

func (s *MemLiquidationStore) UpdateStatus(ctx context.Context, id domain.LiquidationID, status domain.LiquidationStatus, actualAmount, variance int64, expectedVersion int) (*domain.Liquidation, error) {
	s.mu.Lock(); defer s.mu.Unlock()
	l, ok := s.items[id]
	if !ok { return nil, fmt.Errorf("liquidation not found") }
	if l.Version != expectedVersion { return nil, ErrVersionConflict }
	l.Status = status
	l.ActualAmount = actualAmount
	l.VarianceAmount = variance
	l.Version++
	return l, nil
}

// --- Liquidation Evidence ---

type MemLiquidationEvidenceStore struct {
	mu    sync.RWMutex
	items map[domain.LiquidationID][]*domain.LiquidationEvidence
}

func NewMemLiquidationEvidenceStore() *MemLiquidationEvidenceStore {
	return &MemLiquidationEvidenceStore{items: make(map[domain.LiquidationID][]*domain.LiquidationEvidence)}
}

func (s *MemLiquidationEvidenceStore) Create(ctx context.Context, e *domain.LiquidationEvidence) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[e.LiquidationID] = append(s.items[e.LiquidationID], e)
	return nil
}

func (s *MemLiquidationEvidenceStore) ListByLiquidationID(ctx context.Context, lID domain.LiquidationID) ([]*domain.LiquidationEvidence, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.items[lID], nil
}

// --- Billing Record ---

type MemBillingRecordStore struct {
	mu    sync.RWMutex
	items map[domain.BillingRecordID]*domain.BillingRecord
}

func NewMemBillingRecordStore() *MemBillingRecordStore {
	return &MemBillingRecordStore{items: make(map[domain.BillingRecordID]*domain.BillingRecord)}
}

func (s *MemBillingRecordStore) Seed(b *domain.BillingRecord) {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[b.ID] = b
}

func (s *MemBillingRecordStore) Create(ctx context.Context, b *domain.BillingRecord) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[b.ID] = b
	return nil
}

func (s *MemBillingRecordStore) GetByID(ctx context.Context, id domain.BillingRecordID) (*domain.BillingRecord, error) {
	return s.GetByIDWithLines(ctx, id)
}

func (s *MemBillingRecordStore) GetByIDWithLines(ctx context.Context, id domain.BillingRecordID) (*domain.BillingRecord, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	b, ok := s.items[id]
	if !ok { return nil, fmt.Errorf("billing record not found") }
	return b, nil
}

func (s *MemBillingRecordStore) List(ctx context.Context, clientID *domain.ClientID) ([]*domain.BillingRecord, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	out := make([]*domain.BillingRecord, 0)
	for _, b := range s.items {
		if clientID == nil || b.ClientID == *clientID {
			out = append(out, b)
		}
	}
	return out, nil
}

func (s *MemBillingRecordStore) UpdateStatus(ctx context.Context, id domain.BillingRecordID, status domain.BillingStatus, expectedVersion int) (*domain.BillingRecord, error) {
	s.mu.Lock(); defer s.mu.Unlock()
	b, ok := s.items[id]
	if !ok { return nil, fmt.Errorf("billing record not found") }
	if b.Version != expectedVersion { return nil, ErrVersionConflict }
	b.Status = status
	b.Version++
	return b, nil
}

func (s *MemBillingRecordStore) UpdateTotals(ctx context.Context, id domain.BillingRecordID, subtotal, tax, total int64) error {
	s.mu.Lock(); defer s.mu.Unlock()
	b, ok := s.items[id]
	if !ok { return fmt.Errorf("billing record not found") }
	b.Subtotal = subtotal
	b.TaxAmount = tax
	b.Total = total
	return nil
}

func (s *MemBillingRecordStore) MarkReplaced(ctx context.Context, id domain.BillingRecordID, replacedBy domain.BillingRecordID) error {
	s.mu.Lock(); defer s.mu.Unlock()
	b, ok := s.items[id]
	if !ok { return fmt.Errorf("billing record not found") }
	b.ReplacedByID = &replacedBy
	return nil
}

// --- Credit Memo ---

type MemCreditMemoStore struct {
	mu    sync.RWMutex
	items map[domain.ClientID][]*domain.CreditMemo
}

func NewMemCreditMemoStore() *MemCreditMemoStore {
	return &MemCreditMemoStore{items: make(map[domain.ClientID][]*domain.CreditMemo)}
}

func (s *MemCreditMemoStore) Create(ctx context.Context, c *domain.CreditMemo) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[c.ClientID] = append(s.items[c.ClientID], c)
	return nil
}

func (s *MemCreditMemoStore) ListByClientID(ctx context.Context, cID domain.ClientID) ([]*domain.CreditMemo, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.items[cID], nil
}

// --- Client Payment ---

type MemClientPaymentStore struct {
	mu    sync.RWMutex
	items map[domain.ClientPaymentID]*domain.ClientPayment
}

func NewMemClientPaymentStore() *MemClientPaymentStore {
	return &MemClientPaymentStore{items: make(map[domain.ClientPaymentID]*domain.ClientPayment)}
}

func (s *MemClientPaymentStore) Create(ctx context.Context, p *domain.ClientPayment) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items[p.ID] = p
	return nil
}

func (s *MemClientPaymentStore) GetByID(ctx context.Context, id domain.ClientPaymentID) (*domain.ClientPayment, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	p, ok := s.items[id]
	if !ok { return nil, fmt.Errorf("payment not found") }
	return p, nil
}

func (s *MemClientPaymentStore) List(ctx context.Context, clientID *domain.ClientID) ([]*domain.ClientPayment, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	out := make([]*domain.ClientPayment, 0)
	for _, p := range s.items {
		if clientID == nil || p.ClientID == *clientID {
			out = append(out, p)
		}
	}
	return out, nil
}

// --- Billing Allocation ---

type MemBillingAllocationStore struct {
	mu      sync.RWMutex
	byPay   map[domain.ClientPaymentID][]*domain.BillingAllocation
	byBill  map[domain.BillingRecordID][]*domain.BillingAllocation
}

func NewMemBillingAllocationStore() *MemBillingAllocationStore {
	return &MemBillingAllocationStore{
		byPay:  make(map[domain.ClientPaymentID][]*domain.BillingAllocation),
		byBill: make(map[domain.BillingRecordID][]*domain.BillingAllocation),
	}
}

func (s *MemBillingAllocationStore) Create(ctx context.Context, a *domain.BillingAllocation) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.byPay[a.ClientPaymentID] = append(s.byPay[a.ClientPaymentID], a)
	s.byBill[a.BillingRecordID] = append(s.byBill[a.BillingRecordID], a)
	return nil
}

func (s *MemBillingAllocationStore) ListByPaymentID(ctx context.Context, pID domain.ClientPaymentID) ([]*domain.BillingAllocation, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.byPay[pID], nil
}

func (s *MemBillingAllocationStore) ListByBillingID(ctx context.Context, bID domain.BillingRecordID) ([]*domain.BillingAllocation, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.byBill[bID], nil
}
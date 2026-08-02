package domain

import "time"

// --- Funding Source ---

type FundingSourceID string

type FundingSource struct {
	ID           FundingSourceID `json:"id"`
	Name         string          `json:"name"`
	Code         string          `json:"code"`
	IsApproved   bool            `json:"is_approved"`
	BalanceCents int64           `json:"balance_cents"`
	CurrencyCode string          `json:"currency_code"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Version      int             `json:"version"`
}

// --- Budget Request ---

type BudgetRequestID string

type BudgetRequestStatus string

const (
	BudgetRequestStatusDraft          BudgetRequestStatus = "draft"
	BudgetRequestStatusControlsReview  BudgetRequestStatus = "controls_review"
	BudgetRequestStatusApproved       BudgetRequestStatus = "approved"
	BudgetRequestStatusRejected       BudgetRequestStatus = "rejected"
	BudgetRequestStatusReturned       BudgetRequestStatus = "returned"
)

var ValidBudgetRequestTransitions = map[BudgetRequestStatus][]BudgetRequestStatus{
	BudgetRequestStatusDraft:           {BudgetRequestStatusControlsReview},
	BudgetRequestStatusControlsReview:  {BudgetRequestStatusApproved, BudgetRequestStatusRejected, BudgetRequestStatusReturned},
	BudgetRequestStatusApproved:        {}, // terminal
	BudgetRequestStatusRejected:        {}, // terminal
	BudgetRequestStatusReturned:        {BudgetRequestStatusDraft, BudgetRequestStatusControlsReview},
}

func CanTransitionBudgetRequest(from, to BudgetRequestStatus) bool {
	allowed, ok := ValidBudgetRequestTransitions[from]
	if !ok { return false }
	for _, t := range allowed { if t == to { return true } }
	return false
}

type BudgetRequest struct {
	ID            BudgetRequestID      `json:"id"`
	QuotationID   QuotationID          `json:"quotation_id"`
	ClientID      ClientID             `json:"client_id"`
	RequestNumber string               `json:"request_number"`
	Status        BudgetRequestStatus  `json:"status"`
	CurrencyCode  string               `json:"currency_code"`
	AmountCents   int64                `json:"amount_cents"`
	Purpose       string               `json:"purpose"`
	Notes         string               `json:"notes"`
	CreatedBy     *UserID              `json:"created_by,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	ApprovedAt    *time.Time           `json:"approved_at,omitempty"`
	RejectedAt    *time.Time           `json:"rejected_at,omitempty"`
	ReturnedAt    *time.Time           `json:"returned_at,omitempty"`
	Version       int                  `json:"version"`
	Lines         []BudgetLine         `json:"lines"`
}

type BudgetLine struct {
	ID              string             `json:"id"`
	BudgetRequestID BudgetRequestID    `json:"budget_request_id"`
	Description     string             `json:"description"`
	AmountCents     int64              `json:"amount_cents"`
	SortOrder       int                `json:"sort_order"`
	CreatedAt       time.Time          `json:"created_at"`
}

// --- Approval Decision ---

type ApprovalDecisionType string

const (
	ApprovalDecisionApproved ApprovalDecisionType = "approved"
	ApprovalDecisionRejected ApprovalDecisionType = "rejected"
	ApprovalDecisionReturned ApprovalDecisionType = "returned"
)

type ApprovalDecision struct {
	ID              string               `json:"id"`
	BudgetRequestID BudgetRequestID      `json:"budget_request_id"`
	Decision        ApprovalDecisionType `json:"decision"`
	ActorID         UserID               `json:"actor_id"`
	Reason          string               `json:"reason"`
	CreatedAt       time.Time            `json:"created_at"`
}

// --- Disbursement ---

type DisbursementID string

type DisbursementStatus string

const (
	DisbursementStatusPending  DisbursementStatus = "pending"
	DisbursementStatusReleased DisbursementStatus = "released"
	DisbursementStatusHeld     DisbursementStatus = "held"
	DisbursementStatusReturned DisbursementStatus = "returned"
)

var ValidDisbursementTransitions = map[DisbursementStatus][]DisbursementStatus{
	DisbursementStatusPending:  {DisbursementStatusReleased, DisbursementStatusHeld},
	DisbursementStatusReleased: {DisbursementStatusHeld, DisbursementStatusReturned},
	DisbursementStatusHeld:     {DisbursementStatusReleased, DisbursementStatusReturned},
	DisbursementStatusReturned: {}, // terminal
}

func CanTransitionDisbursement(from, to DisbursementStatus) bool {
	allowed, ok := ValidDisbursementTransitions[from]
	if !ok { return false }
	for _, t := range allowed { if t == to { return true } }
	return false
}

type Disbursement struct {
	ID              DisbursementID    `json:"id"`
	BudgetRequestID BudgetRequestID   `json:"budget_request_id"`
	FundingSourceID FundingSourceID   `json:"funding_source_id"`
	Status           DisbursementStatus `json:"status"`
	AmountCents      int64             `json:"amount_cents"`
	CurrencyCode     string            `json:"currency_code"`
	ReferenceNumber  string            `json:"reference_number"`
	Notes            string            `json:"notes"`
	CreatedBy        *UserID           `json:"created_by,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	ReleasedAt       *time.Time        `json:"released_at,omitempty"`
	HeldAt           *time.Time        `json:"held_at,omitempty"`
	ReturnedAt       *time.Time        `json:"returned_at,omitempty"`
	Version          int               `json:"version"`
	PaymentProofs    []PaymentProof    `json:"payment_proofs"`
}

type PaymentProof struct {
	ID             string         `json:"id"`
	DisbursementID DisbursementID `json:"disbursement_id"`
	DocumentName   string         `json:"document_name"`
	StorageKey     string         `json:"storage_key"`
	UploadedBy     *UserID        `json:"uploaded_by,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
}

// --- Liquidation ---

type LiquidationID string

type LiquidationStatus string

const (
	LiquidationStatusOpen       LiquidationStatus = "open"
	LiquidationStatusReconciled LiquidationStatus = "reconciled"
	LiquidationStatusClosed     LiquidationStatus = "closed"
)

var ValidLiquidationTransitions = map[LiquidationStatus][]LiquidationStatus{
	LiquidationStatusOpen:       {LiquidationStatusReconciled, LiquidationStatusClosed},
	LiquidationStatusReconciled:  {LiquidationStatusClosed, LiquidationStatusOpen},
	LiquidationStatusClosed:      {}, // terminal
}

func CanTransitionLiquidation(from, to LiquidationStatus) bool {
	allowed, ok := ValidLiquidationTransitions[from]
	if !ok { return false }
	for _, t := range allowed { if t == to { return true } }
	return false
}

type Liquidation struct {
	ID             LiquidationID        `json:"id"`
	DisbursementID DisbursementID       `json:"disbursement_id"`
	Status         LiquidationStatus    `json:"status"`
	ReleasedAmount int64                `json:"released_amount"`
	ActualAmount   int64                `json:"actual_amount"`
	VarianceAmount int64                `json:"variance_amount"`
	Notes          string               `json:"notes"`
	CreatedBy      *UserID              `json:"created_by,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	ClosedAt       *time.Time           `json:"closed_at,omitempty"`
	Version        int                  `json:"version"`
	Evidence       []LiquidationEvidence `json:"evidence"`
}

type LiquidationEvidence struct {
	ID            string        `json:"id"`
	LiquidationID LiquidationID `json:"liquidation_id"`
	DocumentName  string        `json:"document_name"`
	StorageKey    string        `json:"storage_key"`
	UploadedBy    *UserID       `json:"uploaded_by,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
}

// --- Billing ---

type BillingRecordID string

type BillingStatus string

const (
	BillingStatusDraft     BillingStatus = "draft"
	BillingStatusReview    BillingStatus = "review"
	BillingStatusApproved  BillingStatus = "approved"
	BillingStatusFinalized BillingStatus = "finalized"
	BillingStatusVoid      BillingStatus = "void"
	BillingStatusReplaced  BillingStatus = "replaced"
)

var ValidBillingTransitions = map[BillingStatus][]BillingStatus{
	BillingStatusDraft:     {BillingStatusReview, BillingStatusVoid},
	BillingStatusReview:     {BillingStatusApproved, BillingStatusDraft, BillingStatusVoid},
	BillingStatusApproved:   {BillingStatusFinalized, BillingStatusVoid},
	BillingStatusFinalized:  {BillingStatusVoid, BillingStatusReplaced},
	BillingStatusVoid:       {}, // terminal
	BillingStatusReplaced:   {}, // terminal
}

func CanTransitionBilling(from, to BillingStatus) bool {
	allowed, ok := ValidBillingTransitions[from]
	if !ok { return false }
	for _, t := range allowed { if t == to { return true } }
	return false
}

func IsBillingImmutable(s BillingStatus) bool {
	return s == BillingStatusFinalized || s == BillingStatusVoid || s == BillingStatusReplaced
}

type BillingRecord struct {
	ID             BillingRecordID    `json:"id"`
	ClientID       ClientID           `json:"client_id"`
	BudgetRequestID *BudgetRequestID  `json:"budget_request_id,omitempty"`
	BillingNumber  string             `json:"billing_number"`
	Status         BillingStatus      `json:"status"`
	CurrencyCode   string             `json:"currency_code"`
	Subtotal       int64              `json:"subtotal"`
	TaxAmount      int64              `json:"tax_amount"`
	Total          int64              `json:"total"`
	Notes          string             `json:"notes"`
	CreatedBy      *UserID            `json:"created_by,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	ApprovedAt     *time.Time         `json:"approved_at,omitempty"`
	FinalizedAt    *time.Time         `json:"finalized_at,omitempty"`
	VoidedAt       *time.Time         `json:"voided_at,omitempty"`
	ReplacedByID   *BillingRecordID   `json:"replaced_by_id,omitempty"`
	Version        int                `json:"version"`
	Lines          []BillingLine      `json:"lines"`
}

type BillingLine struct {
	ID              string          `json:"id"`
	BillingRecordID BillingRecordID `json:"billing_record_id"`
	Description     string          `json:"description"`
	Quantity        int64           `json:"quantity"`
	UnitPrice       int64           `json:"unit_price"`
	LineTotal       int64           `json:"line_total"`
	SortOrder       int             `json:"sort_order"`
	CreatedAt       time.Time      `json:"created_at"`
}

// --- Credit Memo ---

type CreditMemoID string

type CreditMemo struct {
	ID             CreditMemoID      `json:"id"`
	ClientID       ClientID           `json:"client_id"`
	BillingRecordID *BillingRecordID `json:"billing_record_id,omitempty"`
	MemoNumber     string            `json:"memo_number"`
	AmountCents    int64              `json:"amount_cents"`
	CurrencyCode   string            `json:"currency_code"`
	Reason         string            `json:"reason"`
	CreatedBy      *UserID            `json:"created_by,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	Version        int                `json:"version"`
}

// --- Client Payment & Allocation ---

type ClientPaymentID string

type ClientPayment struct {
	ID             ClientPaymentID `json:"id"`
	ClientID       ClientID        `json:"client_id"`
	PaymentNumber  string          `json:"payment_number"`
	AmountCents    int64           `json:"amount_cents"`
	CurrencyCode   string          `json:"currency_code"`
	PaymentMethod  string           `json:"payment_method"`
	ReferenceNumber string         `json:"reference_number"`
	ReceivedAt     time.Time       `json:"received_at"`
	CreatedBy      *UserID          `json:"created_by,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	Version        int             `json:"version"`
}

type BillingAllocation struct {
	ID              string           `json:"id"`
	ClientPaymentID  ClientPaymentID  `json:"client_payment_id"`
	BillingRecordID  BillingRecordID `json:"billing_record_id"`
	AmountCents     int64            `json:"amount_cents"`
	AllocatedAt     time.Time        `json:"allocated_at"`
	CreatedAt       time.Time        `json:"created_at"`
}

// --- Audit action constants for workflow ---

const (
	AuditActionBudgetRequestCreated        = "budget_request_created"
	AuditActionBudgetRequestSubmitted       = "budget_request_submitted_for_controls_review"
	AuditActionBudgetRequestApproved        = "budget_request_approved"
	AuditActionBudgetRequestRejected        = "budget_request_rejected"
	AuditActionBudgetRequestReturned        = "budget_request_returned"
	AuditActionDisbursementCreated          = "disbursement_created"
	AuditActionDisbursementReleased         = "disbursement_released"
	AuditActionDisbursementHeld             = "disbursement_held"
	AuditActionDisbursementReturned         = "disbursement_returned"
	AuditActionLiquidationCreated           = "liquidation_created"
	AuditActionLiquidationReconciled        = "liquidation_reconciled"
	AuditActionLiquidationClosed            = "liquidation_closed"
	AuditActionBillingCreated               = "billing_record_created"
	AuditActionBillingSubmitted             = "billing_submitted_for_review"
	AuditActionBillingApproved              = "billing_approved"
	AuditActionBillingFinalized             = "billing_finalized"
	AuditActionBillingVoided                = "billing_voided"
	AuditActionBillingReplaced              = "billing_replaced"
	AuditActionClientPaymentReceived        = "client_payment_received"
	AuditActionBillingAllocationCreated     = "billing_allocation_created"
)
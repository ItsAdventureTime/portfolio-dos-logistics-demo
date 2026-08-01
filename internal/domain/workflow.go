package domain

import "time"

// --- Funding Source ---

type FundingSourceID string

type FundingSource struct {
	ID           FundingSourceID
	Name         string
	Code         string
	IsApproved   bool
	BalanceCents int64
	CurrencyCode string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Version      int
}

// --- Budget Request (Shipment Funding Request) ---

type BudgetRequestID string

type BudgetRequestStatus string

const (
	BudgetRequestStatusDraft         BudgetRequestStatus = "draft"
	BudgetRequestStatusControlsReview BudgetRequestStatus = "controls_review"
	BudgetRequestStatusApproved      BudgetRequestStatus = "approved"
	BudgetRequestStatusRejected      BudgetRequestStatus = "rejected"
	BudgetRequestStatusReturned      BudgetRequestStatus = "returned"
)

var ValidBudgetRequestTransitions = map[BudgetRequestStatus][]BudgetRequestStatus{
	BudgetRequestStatusDraft:          {BudgetRequestStatusControlsReview},
	BudgetRequestStatusControlsReview:  {BudgetRequestStatusApproved, BudgetRequestStatusRejected, BudgetRequestStatusReturned},
	BudgetRequestStatusApproved:        {}, // terminal
	BudgetRequestStatusRejected:        {}, // terminal
	BudgetRequestStatusReturned:        {BudgetRequestStatusDraft, BudgetRequestStatusControlsReview},
}

func CanTransitionBudgetRequest(from, to BudgetRequestStatus) bool {
	allowed, ok := ValidBudgetRequestTransitions[from]
	if !ok {
		return false
	}
	for _, t := range allowed {
		if t == to {
			return true
		}
	}
	return false
}

type BudgetRequest struct {
	ID            BudgetRequestID
	QuotationID   QuotationID
	ClientID      ClientID
	RequestNumber string
	Status        BudgetRequestStatus
	CurrencyCode  string
	AmountCents   int64
	Purpose       string
	Notes         string
	CreatedBy     *UserID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ApprovedAt    *time.Time
	RejectedAt    *time.Time
	ReturnedAt    *time.Time
	Version       int
	Lines         []BudgetLine
}

type BudgetLine struct {
	ID              string
	BudgetRequestID BudgetRequestID
	Description     string
	AmountCents     int64
	SortOrder       int
	CreatedAt       time.Time
}

// --- Approval Decision (Controls Review) ---

type ApprovalDecisionType string

const (
	ApprovalDecisionApproved ApprovalDecisionType = "approved"
	ApprovalDecisionRejected ApprovalDecisionType = "rejected"
	ApprovalDecisionReturned ApprovalDecisionType = "returned"
)

type ApprovalDecision struct {
	ID              string
	BudgetRequestID BudgetRequestID
	Decision        ApprovalDecisionType
	ActorID         UserID
	Reason          string
	CreatedAt       time.Time
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
	if !ok {
		return false
	}
	for _, t := range allowed {
		if t == to {
			return true
		}
	}
	return false
}

type Disbursement struct {
	ID              DisbursementID
	BudgetRequestID  BudgetRequestID
	FundingSourceID  FundingSourceID
	Status           DisbursementStatus
	AmountCents      int64
	CurrencyCode     string
	ReferenceNumber  string
	Notes            string
	CreatedBy        *UserID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ReleasedAt       *time.Time
	HeldAt           *time.Time
	ReturnedAt       *time.Time
	Version          int
	PaymentProofs    []PaymentProof
}

type PaymentProof struct {
	ID             string
	DisbursementID DisbursementID
	DocumentName   string
	StorageKey     string
	UploadedBy     *UserID
	CreatedAt      time.Time
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
	if !ok {
		return false
	}
	for _, t := range allowed {
		if t == to {
			return true
		}
	}
	return false
}

type Liquidation struct {
	ID             LiquidationID
	DisbursementID DisbursementID
	Status         LiquidationStatus
	ReleasedAmount int64
	ActualAmount   int64
	VarianceAmount int64
	Notes          string
	CreatedBy      *UserID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ClosedAt       *time.Time
	Version        int
	Evidence       []LiquidationEvidence
}

type LiquidationEvidence struct {
	ID            string
	LiquidationID LiquidationID
	DocumentName  string
	StorageKey    string
	UploadedBy    *UserID
	CreatedAt     time.Time
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
	if !ok {
		return false
	}
	for _, t := range allowed {
		if t == to {
			return true
		}
	}
	return false
}

func IsBillingImmutable(s BillingStatus) bool {
	return s == BillingStatusFinalized || s == BillingStatusVoid || s == BillingStatusReplaced
}

type BillingRecord struct {
	ID             BillingRecordID
	ClientID       ClientID
	BudgetRequestID *BudgetRequestID
	BillingNumber  string
	Status         BillingStatus
	CurrencyCode   string
	Subtotal       int64
	TaxAmount      int64
	Total          int64
	Notes          string
	CreatedBy      *UserID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ApprovedAt     *time.Time
	FinalizedAt    *time.Time
	VoidedAt       *time.Time
	ReplacedByID   *BillingRecordID
	Version        int
	Lines          []BillingLine
}

type BillingLine struct {
	ID              string
	BillingRecordID BillingRecordID
	Description     string
	Quantity        int64
	UnitPrice       int64
	LineTotal       int64
	SortOrder       int
	CreatedAt       time.Time
}

// --- Credit Memo ---

type CreditMemoID string

type CreditMemo struct {
	ID             CreditMemoID
	ClientID       ClientID
	BillingRecordID *BillingRecordID
	MemoNumber     string
	AmountCents    int64
	CurrencyCode   string
	Reason         string
	CreatedBy      *UserID
	CreatedAt      time.Time
	Version        int
}

// --- Client Payment & Allocation ---

type ClientPaymentID string

type ClientPayment struct {
	ID             ClientPaymentID
	ClientID       ClientID
	PaymentNumber  string
	AmountCents    int64
	CurrencyCode   string
	PaymentMethod  string
	ReferenceNumber string
	ReceivedAt     time.Time
	CreatedBy      *UserID
	CreatedAt      time.Time
	Version        int
}

type BillingAllocation struct {
	ID              string
	ClientPaymentID  ClientPaymentID
	BillingRecordID  BillingRecordID
	AmountCents     int64
	AllocatedAt     time.Time
	CreatedAt       time.Time
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
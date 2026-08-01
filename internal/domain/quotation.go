package domain

import (
	"time"
)

// ClientID uniquely identifies a client.
type ClientID string

// QuotationID uniquely identifies a quotation.
type QuotationID string

// Client is the company whose freight is managed.
type Client struct {
	ID           ClientID
	Name         string
	Code         string
	ContactEmail string
	ContactPhone string
	Address      string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Version      int
}

// QuotationStatus represents the state of a quotation in its lifecycle.
type QuotationStatus string

const (
	QuotationStatusDraft    QuotationStatus = "draft"
	QuotationStatusReview   QuotationStatus = "review"
	QuotationStatusApproved QuotationStatus = "approved"
	QuotationStatusAccepted QuotationStatus = "accepted"
	QuotationStatusRevised  QuotationStatus = "revised"
	QuotationStatusVoid     QuotationStatus = "void"
)

// Quotation is a freight quotation with state machine lifecycle.
type Quotation struct {
	ID             QuotationID
	ClientID       ClientID
	QuotationNumber string
	Status         QuotationStatus
	CurrencyCode   string
	Subtotal       int64 // stored as cents to avoid float issues
	TaxAmount      int64
	Total          int64
	Notes          string
	CreatedBy      *UserID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	AcceptedAt     *time.Time
	VoidedAt       *time.Time
	Version        int
	Lines          []QuotationLine
}

// QuotationLine is an individual charge on a quotation.
type QuotationLine struct {
	ID          string
	QuotationID QuotationID
	Description  string
	Quantity     int64 // stored as fixed-point (e.g., 100 = 1.00)
	UnitPrice    int64 // stored as cents
	LineTotal    int64 // stored as cents
	SortOrder    int
	CreatedAt    time.Time
}

// QuotationTransition represents a valid state transition.
type QuotationTransition struct {
	From QuotationStatus
	To   QuotationStatus
}

// ValidQuotationTransitions defines the legal state transitions per
// docs/spec/02-domain-and-workflows.md.
var ValidQuotationTransitions = map[QuotationStatus][]QuotationStatus{
	QuotationStatusDraft:    {QuotationStatusReview, QuotationStatusVoid},
	QuotationStatusReview:    {QuotationStatusApproved, QuotationStatusDraft, QuotationStatusVoid},
	QuotationStatusApproved:  {QuotationStatusAccepted, QuotationStatusRevised, QuotationStatusVoid},
	QuotationStatusAccepted:  {}, // terminal
	QuotationStatusRevised:   {QuotationStatusReview, QuotationStatusVoid},
	QuotationStatusVoid:      {}, // terminal
}

// CanTransitionQuotation reports whether a transition is valid.
func CanTransitionQuotation(from, to QuotationStatus) bool {
	allowed, ok := ValidQuotationTransitions[from]
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

// IsTerminalQuotationStatus reports whether the status is terminal.
func IsTerminalQuotationStatus(s QuotationStatus) bool {
	return s == QuotationStatusAccepted || s == QuotationStatusVoid
}

// Audit action constants for quotations.
const (
	AuditActionQuotationCreated   = "quotation_created"
	AuditActionQuotationSubmitted = "quotation_submitted_for_review"
	AuditActionQuotationApproved  = "quotation_approved"
	AuditActionQuotationAccepted  = "quotation_accepted"
	AuditActionQuotationRevised   = "quotation_revised"
	AuditActionQuotationVoided    = "quotation_voided"
	AuditActionQuotationRejected  = "quotation_rejected"
	AuditActionQuotationReturned  = "quotation_returned"
)

// Audit action constants for clients.
const (
	AuditActionClientCreated = "client_created"
	AuditActionClientUpdated = "client_updated"
)
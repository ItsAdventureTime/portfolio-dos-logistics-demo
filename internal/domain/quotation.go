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
	ID           ClientID    `json:"id"`
	Name         string      `json:"name"`
	Code         string      `json:"code"`
	ContactEmail string      `json:"contact_email"`
	ContactPhone string      `json:"contact_phone"`
	Address      string      `json:"address"`
	IsActive     bool        `json:"is_active"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Version      int         `json:"version"`
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
	ID              QuotationID      `json:"id"`
	ClientID        ClientID         `json:"client_id"`
	QuotationNumber string           `json:"quotation_number"`
	Status          QuotationStatus  `json:"status"`
	CurrencyCode    string           `json:"currency_code"`
	Subtotal        int64            `json:"subtotal"`
	TaxAmount       int64            `json:"tax_amount"`
	Total           int64            `json:"total"`
	Notes           string           `json:"notes"`
	CreatedBy       *UserID          `json:"created_by,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	AcceptedAt      *time.Time       `json:"accepted_at,omitempty"`
	VoidedAt        *time.Time       `json:"voided_at,omitempty"`
	Version         int              `json:"version"`
	Lines           []QuotationLine  `json:"lines"`
}

// QuotationLine is an individual charge on a quotation.
type QuotationLine struct {
	ID          string    `json:"id"`
	QuotationID QuotationID `json:"quotation_id"`
	Description string    `json:"description"`
	Quantity    int64     `json:"quantity"`
	UnitPrice   int64     `json:"unit_price"`
	LineTotal   int64     `json:"line_total"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

// QuotationTransition represents a valid state transition.
type QuotationTransition struct {
	From QuotationStatus `json:"from"`
	To   QuotationStatus `json:"to"`
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
// Package workflowhttp wires the workflow service routes onto the Chi
// router. Covers quotations, budget requests, disbursements, liquidations,
// billing, collections, and clients.
package workflowhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/httpserver/middleware"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/service"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	quotations *service.QuotationService
	workflow   *service.WorkflowService
}

// Mount registers all workflow routes on the router. All routes require auth.
func Mount(r chi.Router, qSvc *service.QuotationService, wSvc *service.WorkflowService, validateAuth func(ctx context.Context, tokenHash string) (*domain.Session, *domain.User, error)) {
	h := &handler{quotations: qSvc, workflow: wSvc}

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(validateAuth))
		r.Use(middleware.CSRF)

		// Clients
		r.Get("/api/clients", h.ListClients)
		r.Post("/api/clients", h.CreateClient)

		// Quotations
		r.Get("/api/quotations", h.ListQuotations)
		r.Post("/api/quotations", h.CreateQuotation)
		r.Get("/api/quotations/{id}", h.GetQuotation)
		r.Post("/api/quotations/{id}/transition", h.TransitionQuotation)

		// Budget requests
		r.Get("/api/budget-requests", h.ListBudgetRequests)
		r.Post("/api/budget-requests", h.CreateBudgetRequest)
		r.Get("/api/budget-requests/{id}", h.GetBudgetRequest)
		r.Post("/api/budget-requests/{id}/transition", h.TransitionBudgetRequest)

		// Disbursements
		r.Get("/api/disbursements", h.ListDisbursements)
		r.Post("/api/disbursements", h.CreateDisbursement)
		r.Post("/api/disbursements/{id}/transition", h.TransitionDisbursement)

		// Liquidations
		r.Get("/api/liquidations", h.ListLiquidations)
		r.Post("/api/liquidations", h.CreateLiquidation)
		r.Post("/api/liquidations/{id}/reconcile", h.ReconcileLiquidation)
		r.Post("/api/liquidations/{id}/close", h.CloseLiquidation)

		// Funding sources
		r.Get("/api/funding-sources", h.ListFundingSources)

		// Billing
		r.Get("/api/billing", h.ListBilling)
		r.Post("/api/billing", h.CreateBilling)
		r.Post("/api/billing/{id}/transition", h.TransitionBilling)

		// Collections
		r.Get("/api/payments", h.ListPayments)
		r.Post("/api/payments", h.RecordPayment)
		r.Post("/api/payments/{id}/allocate", h.AllocatePayment)

		// Dashboard
		r.Get("/api/dashboard", h.GetDashboard)
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{"code": code, "message": msg},
	})
}

func actorFromCtx(r *http.Request) domain.UserID {
	ac, ok := middleware.AuthFromContext(r.Context())
	if !ok || ac.User == nil {
		return ""
	}
	return ac.User.ID
}

// --- Clients ---

func (h *handler) ListClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.quotations.ListClients(r.Context())
	if err != nil {
		writeErr(w, 500, "internal_error", "An error occurred.")
		return
	}
	if clients == nil {
		clients = []*domain.Client{}
	}
	writeJSON(w, 200, map[string]any{"clients": clients})
}

type createClientReq struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
	Address     string `json:"address"`
}

func (h *handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var req createClientReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	if req.Name == "" || req.Code == "" || req.ContactEmail == "" {
		writeErr(w, 400, "invalid_request", "Name, code, and contact email are required.")
		return
	}
	c, err := h.quotations.CreateClient(r.Context(), req.Name, req.Code, req.ContactEmail, req.ContactPhone, req.Address)
	if err != nil {
		writeErr(w, 400, "invalid_request", err.Error())
		return
	}
	writeJSON(w, 201, c)
}

// --- Quotations ---

func (h *handler) ListQuotations(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	var cid *domain.ClientID
	if clientID != "" {
		c := domain.ClientID(clientID)
		cid = &c
	}
	quotations, err := h.quotations.ListQuotations(r.Context(), cid)
	if err != nil {
		writeErr(w, 500, "internal_error", "An error occurred.")
		return
	}
	if quotations == nil {
		quotations = []*domain.Quotation{}
	}
	writeJSON(w, 200, map[string]any{"quotations": quotations})
}

type createQuotationReq struct {
	ClientID     string              `json:"client_id"`
	CurrencyCode string              `json:"currency_code"`
	Notes        string              `json:"notes"`
	Lines        []service.QuotationLineInput `json:"lines"`
}

func (h *handler) CreateQuotation(w http.ResponseWriter, r *http.Request) {
	var req createQuotationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	q, err := h.quotations.CreateQuotation(r.Context(), actorFromCtx(r), service.CreateQuotationInput{
		ClientID:     domain.ClientID(req.ClientID),
		CurrencyCode: req.CurrencyCode,
		Notes:        req.Notes,
		Lines:        req.Lines,
	})
	if err != nil {
		writeErr(w, 400, "invalid_request", err.Error())
		return
	}
	writeJSON(w, 201, q)
}

func (h *handler) GetQuotation(w http.ResponseWriter, r *http.Request) {
	id := domain.QuotationID(chi.URLParam(r, "id"))
	q, err := h.quotations.GetQuotation(r.Context(), id)
	if err != nil {
		writeErr(w, 404, "not_found", "Quotation not found.")
		return
	}
	writeJSON(w, 200, q)
}

type transitionReq struct {
	TargetStatus string `json:"target_status"`
	Version      int    `json:"version"`
}

func (h *handler) TransitionQuotation(w http.ResponseWriter, r *http.Request) {
	id := domain.QuotationID(chi.URLParam(r, "id"))
	var req transitionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	q, err := h.quotations.TransitionQuotation(r.Context(), id, domain.QuotationStatus(req.TargetStatus), req.Version, actorFromCtx(r))
	if err != nil {
		writeErr(w, 409, "transition_failed", err.Error())
		return
	}
	writeJSON(w, 200, q)
}

// --- Budget Requests ---

func (h *handler) ListBudgetRequests(w http.ResponseWriter, r *http.Request) {
	brs, err := h.workflow.ListBudgetRequests(r.Context())
	if err != nil {
		writeErr(w, 500, "internal_error", "An error occurred.")
		return
	}
	if brs == nil { brs = []*domain.BudgetRequest{} }
	writeJSON(w, 200, map[string]any{"budget_requests": brs})
}

type createBudgetReq struct {
	QuotationID  string `json:"quotation_id"`
	ClientID     string `json:"client_id"`
	CurrencyCode string `json:"currency_code"`
	Purpose      string `json:"purpose"`
	AmountCents  int64  `json:"amount_cents"`
}

func (h *handler) CreateBudgetRequest(w http.ResponseWriter, r *http.Request) {
	var req createBudgetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	br, err := h.workflow.CreateBudgetRequest(r.Context(), actorFromCtx(r),
		domain.QuotationID(req.QuotationID), domain.ClientID(req.ClientID),
		req.CurrencyCode, req.Purpose, req.AmountCents)
	if err != nil {
		writeErr(w, 400, "invalid_request", err.Error())
		return
	}
	writeJSON(w, 201, br)
}

func (h *handler) GetBudgetRequest(w http.ResponseWriter, r *http.Request) {
	writeErr(w, 404, "not_found", "Budget request not found.")
}

func (h *handler) TransitionBudgetRequest(w http.ResponseWriter, r *http.Request) {
	id := domain.BudgetRequestID(chi.URLParam(r, "id"))
	var req transitionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	br, err := h.workflow.TransitionBudgetRequest(r.Context(), id, domain.BudgetRequestStatus(req.TargetStatus), req.Version, actorFromCtx(r), "")
	if err != nil {
		writeErr(w, 409, "transition_failed", err.Error())
		return
	}
	writeJSON(w, 200, br)
}

// --- Disbursements ---

func (h *handler) ListDisbursements(w http.ResponseWriter, r *http.Request) {
	ds, err := h.workflow.ListDisbursements(r.Context())
	if err != nil {
		writeErr(w, 500, "internal_error", "An error occurred.")
		return
	}
	if ds == nil { ds = []*domain.Disbursement{} }
	writeJSON(w, 200, map[string]any{"disbursements": ds})
}

type createDisbursementReq struct {
	BudgetRequestID  string `json:"budget_request_id"`
	FundingSourceID  string `json:"funding_source_id"`
	AmountCents     int64  `json:"amount_cents"`
	CurrencyCode    string `json:"currency_code"`
	ReferenceNumber string `json:"reference_number"`
	Notes           string `json:"notes"`
}

func (h *handler) CreateDisbursement(w http.ResponseWriter, r *http.Request) {
	var req createDisbursementReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	d, err := h.workflow.CreateDisbursement(r.Context(), actorFromCtx(r),
		domain.BudgetRequestID(req.BudgetRequestID), domain.FundingSourceID(req.FundingSourceID),
		req.AmountCents, req.CurrencyCode, req.ReferenceNumber, req.Notes)
	if err != nil {
		writeErr(w, 400, "invalid_request", err.Error())
		return
	}
	writeJSON(w, 201, d)
}

func (h *handler) TransitionDisbursement(w http.ResponseWriter, r *http.Request) {
	id := domain.DisbursementID(chi.URLParam(r, "id"))
	var req transitionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	d, err := h.workflow.TransitionDisbursement(r.Context(), id, domain.DisbursementStatus(req.TargetStatus), req.Version, actorFromCtx(r))
	if err != nil {
		writeErr(w, 409, "transition_failed", err.Error())
		return
	}
	writeJSON(w, 200, d)
}

// --- Liquidations ---

func (h *handler) ListLiquidations(w http.ResponseWriter, r *http.Request) {
	ls, err := h.workflow.ListLiquidations(r.Context())
	if err != nil {
		writeErr(w, 500, "internal_error", "An error occurred.")
		return
	}
	if ls == nil { ls = []*domain.Liquidation{} }
	writeJSON(w, 200, map[string]any{"liquidations": ls})
}

// --- Funding Sources ---

func (h *handler) ListFundingSources(w http.ResponseWriter, r *http.Request) {
	fss, err := h.workflow.ListFundingSources(r.Context())
	if err != nil {
		writeErr(w, 500, "internal_error", "An error occurred.")
		return
	}
	if fss == nil { fss = []*domain.FundingSource{} }
	writeJSON(w, 200, map[string]any{"funding_sources": fss})
}

type createLiquidationReq struct {
	DisbursementID  string `json:"disbursement_id"`
	ReleasedAmount  int64  `json:"released_amount"`
}

func (h *handler) CreateLiquidation(w http.ResponseWriter, r *http.Request) {
	var req createLiquidationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	l, err := h.workflow.CreateLiquidation(r.Context(), actorFromCtx(r),
		domain.DisbursementID(req.DisbursementID), req.ReleasedAmount)
	if err != nil {
		writeErr(w, 400, "invalid_request", err.Error())
		return
	}
	writeJSON(w, 201, l)
}

type reconcileReq struct {
	ActualAmount  int64 `json:"actual_amount"`
	Version       int   `json:"version"`
}

func (h *handler) ReconcileLiquidation(w http.ResponseWriter, r *http.Request) {
	id := domain.LiquidationID(chi.URLParam(r, "id"))
	var req reconcileReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	l, err := h.workflow.ReconcileLiquidation(r.Context(), id, req.ActualAmount, req.Version, actorFromCtx(r))
	if err != nil {
		writeErr(w, 409, "transition_failed", err.Error())
		return
	}
	writeJSON(w, 200, l)
}

func (h *handler) CloseLiquidation(w http.ResponseWriter, r *http.Request) {
	id := domain.LiquidationID(chi.URLParam(r, "id"))
	var req struct {
		Version int `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	l, err := h.workflow.CloseLiquidation(r.Context(), id, req.Version, actorFromCtx(r))
	if err != nil {
		writeErr(w, 409, "transition_failed", err.Error())
		return
	}
	writeJSON(w, 200, l)
}

// --- Billing ---

func (h *handler) ListBilling(w http.ResponseWriter, r *http.Request) {
	bs, err := h.workflow.ListBillingRecords(r.Context())
	if err != nil {
		writeErr(w, 500, "internal_error", "An error occurred.")
		return
	}
	if bs == nil { bs = []*domain.BillingRecord{} }
	writeJSON(w, 200, map[string]any{"billing_records": bs})
}

type createBillingReq struct {
	ClientID      string                        `json:"client_id"`
	CurrencyCode  string                        `json:"currency_code"`
	Lines         []service.BillingLineInput     `json:"lines"`
}

func (h *handler) CreateBilling(w http.ResponseWriter, r *http.Request) {
	var req createBillingReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	b, err := h.workflow.CreateBillingRecord(r.Context(), actorFromCtx(r),
		domain.ClientID(req.ClientID), req.CurrencyCode, req.Lines, nil)
	if err != nil {
		writeErr(w, 400, "invalid_request", err.Error())
		return
	}
	writeJSON(w, 201, b)
}

func (h *handler) TransitionBilling(w http.ResponseWriter, r *http.Request) {
	id := domain.BillingRecordID(chi.URLParam(r, "id"))
	var req transitionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	b, err := h.workflow.TransitionBilling(r.Context(), id, domain.BillingStatus(req.TargetStatus), req.Version, actorFromCtx(r))
	if err != nil {
		writeErr(w, 409, "transition_failed", err.Error())
		return
	}
	writeJSON(w, 200, b)
}

// --- Collections ---

func (h *handler) ListPayments(w http.ResponseWriter, r *http.Request) {
	ps, err := h.workflow.ListPayments(r.Context())
	if err != nil {
		writeErr(w, 500, "internal_error", "An error occurred.")
		return
	}
	if ps == nil { ps = []*domain.ClientPayment{} }
	writeJSON(w, 200, map[string]any{"payments": ps})
}

type recordPaymentReq struct {
	ClientID        string `json:"client_id"`
	AmountCents    int64  `json:"amount_cents"`
	CurrencyCode   string `json:"currency_code"`
	PaymentMethod  string `json:"payment_method"`
	ReferenceNumber string `json:"reference_number"`
}

func (h *handler) RecordPayment(w http.ResponseWriter, r *http.Request) {
	var req recordPaymentReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	p, err := h.workflow.RecordClientPayment(r.Context(), actorFromCtx(r),
		domain.ClientID(req.ClientID), req.AmountCents, req.CurrencyCode,
		req.PaymentMethod, req.ReferenceNumber)
	if err != nil {
		writeErr(w, 400, "invalid_request", err.Error())
		return
	}
	writeJSON(w, 201, p)
}

type allocateReq struct {
	BillingRecordID string `json:"billing_record_id"`
	AmountCents    int64  `json:"amount_cents"`
}

func (h *handler) AllocatePayment(w http.ResponseWriter, r *http.Request) {
	paymentID := domain.ClientPaymentID(chi.URLParam(r, "id"))
	var req allocateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid_request", "The request body could not be read.")
		return
	}
	a, err := h.workflow.AllocatePayment(r.Context(), actorFromCtx(r),
		paymentID, domain.BillingRecordID(req.BillingRecordID), req.AmountCents)
	if err != nil {
		writeErr(w, 400, "invalid_request", err.Error())
		return
	}
	writeJSON(w, 201, a)
}

// --- Dashboard ---

func (h *handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	// Aggregate from seeded demo data
	clients, _ := h.quotations.ListClients(r.Context())
	quotations, _ := h.quotations.ListQuotations(r.Context(), nil)

	activeQuotations := 0
	pendingApprovals := 0
	totalDisbursed := int64(0)
	outstandingReceivables := int64(0)

	for _, q := range quotations {
		if q.Status != domain.QuotationStatusVoid {
			activeQuotations++
		}
	}

	writeJSON(w, 200, map[string]any{
		"kpis": []map[string]any{
			{"label": "Active quotations", "value": strconv.Itoa(activeQuotations), "trend": "+3 this week", "trend_direction": "up"},
			{"label": "Pending approvals", "value": strconv.Itoa(pendingApprovals), "trend": "Awaiting review", "trend_direction": "neutral"},
			{"label": "Outstanding receivables", "value": "$" + formatCents(outstandingReceivables), "trend": "No overdue", "trend_direction": "neutral"},
			{"label": "Disbursed this month", "value": "$" + formatCents(totalDisbursed), "trend": "No activity yet", "trend_direction": "neutral"},
		},
		"activity": []map[string]any{},
		"client_count": len(clients),
		"quotation_count": len(quotations),
	})
}

func formatCents(cents int64) string {
	return strconv.FormatInt(cents/100, 10) + ".00"
}
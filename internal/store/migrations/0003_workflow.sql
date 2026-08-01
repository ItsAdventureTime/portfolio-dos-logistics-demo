-- +goose Up
-- +goose StatementBegin

-- Funding sources: approved sources of funds for disbursement.
CREATE TABLE funding_sources (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    code            TEXT NOT NULL UNIQUE,
    is_approved     BOOLEAN NOT NULL DEFAULT FALSE,
    balance_cents   NUMERIC(15,2) NOT NULL DEFAULT 0,
    currency_code   TEXT NOT NULL REFERENCES currencies(code),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    version         INTEGER NOT NULL DEFAULT 1
);

-- Budget requests (shipment funding requests).
CREATE TABLE budget_requests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quotation_id    UUID NOT NULL REFERENCES quotations(id) ON DELETE RESTRICT,
    client_id       UUID NOT NULL REFERENCES clients(id) ON DELETE RESTRICT,
    request_number  TEXT NOT NULL UNIQUE,
    status          TEXT NOT NULL DEFAULT 'draft',
    currency_code   TEXT NOT NULL REFERENCES currencies(code),
    amount_cents    NUMERIC(15,2) NOT NULL,
    purpose         TEXT NOT NULL,
    notes           TEXT,
    created_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    approved_at     TIMESTAMPTZ,
    rejected_at     TIMESTAMPTZ,
    returned_at     TIMESTAMPTZ,
    version         INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_budget_requests_quotation_id ON budget_requests(quotation_id);
CREATE INDEX idx_budget_requests_client_id ON budget_requests(client_id);
CREATE INDEX idx_budget_requests_status ON budget_requests(status);

-- Budget lines: detail items in a funding request.
CREATE TABLE budget_lines (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_request_id UUID NOT NULL REFERENCES budget_requests(id) ON DELETE CASCADE,
    description     TEXT NOT NULL,
    amount_cents    NUMERIC(15,2) NOT NULL,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_budget_lines_request_id ON budget_lines(budget_request_id);

-- Approval decisions: controls review records.
CREATE TABLE approval_decisions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_request_id UUID NOT NULL REFERENCES budget_requests(id) ON DELETE CASCADE,
    decision        TEXT NOT NULL, -- approved, rejected, returned
    actor_id        UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    reason          TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_approval_decisions_request_id ON approval_decisions(budget_request_id);

-- Disbursements: released funds.
CREATE TABLE disbursements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_request_id UUID NOT NULL REFERENCES budget_requests(id) ON DELETE RESTRICT,
    funding_source_id UUID NOT NULL REFERENCES funding_sources(id) ON DELETE RESTRICT,
    status          TEXT NOT NULL DEFAULT 'pending',
    amount_cents    NUMERIC(15,2) NOT NULL,
    currency_code   TEXT NOT NULL REFERENCES currencies(code),
    reference_number TEXT,
    notes           TEXT,
    created_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    released_at     TIMESTAMPTZ,
    held_at         TIMESTAMPTZ,
    returned_at     TIMESTAMPTZ,
    version         INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_disbursements_request_id ON disbursements(budget_request_id);
CREATE INDEX idx_disbursements_status ON disbursements(status);

-- Payment proofs: evidence of disbursement.
CREATE TABLE payment_proofs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disbursement_id UUID NOT NULL REFERENCES disbursements(id) ON DELETE CASCADE,
    document_name   TEXT NOT NULL,
    storage_key     TEXT NOT NULL,
    uploaded_by     UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_payment_proofs_disbursement_id ON payment_proofs(disbursement_id);

-- Liquidations: reconciliation of released vs actual spending.
CREATE TABLE liquidations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disbursement_id UUID NOT NULL REFERENCES disbursements(id) ON DELETE RESTRICT,
    status          TEXT NOT NULL DEFAULT 'open',
    released_amount NUMERIC(15,2) NOT NULL,
    actual_amount   NUMERIC(15,2) NOT NULL DEFAULT 0,
    variance_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    notes           TEXT,
    created_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    closed_at       TIMESTAMPTZ,
    version         INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_liquidations_disbursement_id ON liquidations(disbursement_id);

-- Liquidation evidence: supporting documents for liquidation.
CREATE TABLE liquidation_evidence (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    liquidation_id  UUID NOT NULL REFERENCES liquidations(id) ON DELETE CASCADE,
    document_name   TEXT NOT NULL,
    storage_key     TEXT NOT NULL,
    uploaded_by     UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_liquidation_evidence_liquidation_id ON liquidation_evidence(liquidation_id);

-- Billing records: invoices for the client.
CREATE TABLE billing_records (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id       UUID NOT NULL REFERENCES clients(id) ON DELETE RESTRICT,
    budget_request_id UUID REFERENCES budget_requests(id) ON DELETE SET NULL,
    billing_number  TEXT NOT NULL UNIQUE,
    status          TEXT NOT NULL DEFAULT 'draft',
    currency_code   TEXT NOT NULL REFERENCES currencies(code),
    subtotal        NUMERIC(15,2) NOT NULL DEFAULT 0,
    tax_amount      NUMERIC(15,2) NOT NULL DEFAULT 0,
    total           NUMERIC(15,2) NOT NULL DEFAULT 0,
    notes           TEXT,
    created_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    approved_at     TIMESTAMPTZ,
    finalized_at    TIMESTAMPTZ,
    voided_at       TIMESTAMPTZ,
    replaced_by_id  UUID REFERENCES billing_records(id) ON DELETE SET NULL,
    version         INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_billing_records_client_id ON billing_records(client_id);
CREATE INDEX idx_billing_records_status ON billing_records(status);

-- Billing lines: individual charges on a billing record.
CREATE TABLE billing_lines (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    billing_record_id UUID NOT NULL REFERENCES billing_records(id) ON DELETE CASCADE,
    description     TEXT NOT NULL,
    quantity        NUMERIC(10,2) NOT NULL DEFAULT 1,
    unit_price      NUMERIC(15,2) NOT NULL DEFAULT 0,
    line_total      NUMERIC(15,2) NOT NULL DEFAULT 0,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_billing_lines_record_id ON billing_lines(billing_record_id);

-- Credit memos: adjustments to billing.
CREATE TABLE credit_memos (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id       UUID NOT NULL REFERENCES clients(id) ON DELETE RESTRICT,
    billing_record_id UUID REFERENCES billing_records(id) ON DELETE SET NULL,
    memo_number     TEXT NOT NULL UNIQUE,
    amount_cents    NUMERIC(15,2) NOT NULL,
    currency_code   TEXT NOT NULL REFERENCES currencies(code),
    reason          TEXT,
    created_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    version         INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_credit_memos_client_id ON credit_memos(client_id);

-- Client payments: payments received from clients.
CREATE TABLE client_payments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id       UUID NOT NULL REFERENCES clients(id) ON DELETE RESTRICT,
    payment_number  TEXT NOT NULL UNIQUE,
    amount_cents    NUMERIC(15,2) NOT NULL,
    currency_code   TEXT NOT NULL REFERENCES currencies(code),
    payment_method  TEXT,
    reference_number TEXT,
    received_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    version         INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_client_payments_client_id ON client_payments(client_id);

-- Billing allocations: how payments are allocated to billing records.
CREATE TABLE billing_allocations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_payment_id UUID NOT NULL REFERENCES client_payments(id) ON DELETE CASCADE,
    billing_record_id UUID NOT NULL REFERENCES billing_records(id) ON DELETE RESTRICT,
    amount_cents    NUMERIC(15,2) NOT NULL,
    allocated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_billing_allocations_payment_id ON billing_allocations(client_payment_id);
CREATE INDEX idx_billing_allocations_billing_id ON billing_allocations(billing_record_id);

-- Supporting documents: file uploads scoped to entities.
CREATE TABLE supporting_documents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type     TEXT NOT NULL,
    entity_id       UUID NOT NULL,
    document_name   TEXT NOT NULL,
    storage_key     TEXT NOT NULL,
    content_type    TEXT NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    uploaded_by     UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_supporting_documents_entity ON supporting_documents(entity_type, entity_id);

-- Tax profiles: tax configuration for billing.
CREATE TABLE tax_profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    tax_rate        NUMERIC(5,4) NOT NULL DEFAULT 0,
    description     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tax_profiles;
DROP TABLE IF EXISTS supporting_documents;
DROP TABLE IF EXISTS billing_allocations;
DROP TABLE IF EXISTS client_payments;
DROP TABLE IF EXISTS credit_memos;
DROP TABLE IF EXISTS billing_lines;
DROP TABLE IF EXISTS billing_records;
DROP TABLE IF EXISTS liquidation_evidence;
DROP TABLE IF EXISTS liquidations;
DROP TABLE IF EXISTS payment_proofs;
DROP TABLE IF EXISTS disbursements;
DROP TABLE IF EXISTS approval_decisions;
DROP TABLE IF EXISTS budget_lines;
DROP TABLE IF EXISTS budget_requests;
DROP TABLE IF EXISTS funding_sources;
-- +goose StatementEnd
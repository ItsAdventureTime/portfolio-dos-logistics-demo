import type { LoginResponse, User, Session } from '../types'

const API_BASE = ''

function getCSRFToken(): string | null {
  const match = document.cookie.match(/dos_csrf=([^;]+)/)
  return match ? match[1] : null
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...((options.headers as Record<string, string>) || {}),
  }

  if (options.method && options.method !== 'GET') {
    const csrf = getCSRFToken()
    if (csrf) {
      headers['X-CSRF-Token'] = csrf
    }
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers, credentials: 'same-origin' })

  if (res.status === 401) {
    throw new Error('Authentication required. Please sign in.')
  }
  if (res.status === 403) {
    throw new Error('You do not have permission to perform this action.')
  }
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    const message = body?.error?.message || 'Something went wrong. Please try again.'
    throw new Error(message)
  }

  return res.json() as Promise<T>
}

export const api = {
  // Auth
  login: (identifier: string, password: string) =>
    request<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ identifier, password }),
    }),

  verifyEmail: (identifier: string, code: string) =>
    request<LoginResponse>('/auth/verify-email', {
      method: 'POST',
      body: JSON.stringify({ identifier, code }),
    }),

  logout: () => request<{ status: string }>('/auth/logout', { method: 'POST' }),

  me: () => request<User>('/auth/me'),

  session: () => request<Session>('/auth/session'),

  setRolePreview: (role_preview: string) =>
    request<{ status: string; role_preview: string }>('/auth/role-preview', {
      method: 'POST',
      body: JSON.stringify({ role_preview }),
    }),

  // Dashboard
  getDashboard: () => request<DashboardData>('/api/dashboard'),

  // Clients
  listClients: () => request<{ clients: Client[] }>('/api/clients'),
  createClient: (data: { name: string; code: string; contact_email: string; contact_phone?: string; address?: string }) =>
    request<Client>('/api/clients', { method: 'POST', body: JSON.stringify(data) }),

  // Quotations
  listQuotations: (client_id?: string) =>
    request<{ quotations: Quotation[] }>(`/api/quotations${client_id ? `?client_id=${client_id}` : ''}`),
  createQuotation: (data: { client_id: string; currency_code: string; notes?: string; lines: { description: string; quantity: number; unit_price: number }[] }) =>
    request<Quotation>('/api/quotations', { method: 'POST', body: JSON.stringify(data) }),
  getQuotation: (id: string) => request<Quotation>(`/api/quotations/${id}`),
  transitionQuotation: (id: string, target_status: string, version: number) =>
    request<Quotation>(`/api/quotations/${id}/transition`, { method: 'POST', body: JSON.stringify({ target_status, version }) }),

  // Budget requests
  listBudgetRequests: () => request<{ budget_requests: BudgetRequest[] }>('/api/budget-requests'),
  createBudgetRequest: (data: { quotation_id: string; client_id: string; currency_code: string; purpose: string; amount_cents: number }) =>
    request<BudgetRequest>('/api/budget-requests', { method: 'POST', body: JSON.stringify(data) }),
  transitionBudgetRequest: (id: string, target_status: string, version: number) =>
    request<BudgetRequest>(`/api/budget-requests/${id}/transition`, { method: 'POST', body: JSON.stringify({ target_status, version }) }),

  // Disbursements
  listDisbursements: () => request<{ disbursements: Disbursement[] }>('/api/disbursements'),
  createDisbursement: (data: { budget_request_id: string; funding_source_id: string; amount_cents: number; currency_code: string; reference_number?: string; notes?: string }) =>
    request<Disbursement>('/api/disbursements', { method: 'POST', body: JSON.stringify(data) }),
  transitionDisbursement: (id: string, target_status: string, version: number) =>
    request<Disbursement>(`/api/disbursements/${id}/transition`, { method: 'POST', body: JSON.stringify({ target_status, version }) }),

  // Liquidations
  listLiquidations: () => request<{ liquidations: Liquidation[] }>('/api/liquidations'),
  createLiquidation: (disbursement_id: string, released_amount: number) =>
    request<Liquidation>('/api/liquidations', { method: 'POST', body: JSON.stringify({ disbursement_id, released_amount }) }),
  reconcileLiquidation: (id: string, actual_amount: number, version: number) =>
    request<Liquidation>(`/api/liquidations/${id}/reconcile`, { method: 'POST', body: JSON.stringify({ actual_amount, version }) }),
  closeLiquidation: (id: string, version: number) =>
    request<Liquidation>(`/api/liquidations/${id}/close`, { method: 'POST', body: JSON.stringify({ version }) }),

  // Funding sources
  listFundingSources: () => request<{ funding_sources: FundingSource[] }>('/api/funding-sources'),

  // Billing
  listBilling: () => request<{ billing_records: BillingRecord[] }>('/api/billing'),
  createBilling: (data: { client_id: string; currency_code: string; lines: { description: string; quantity: number; unit_price: number }[] }) =>
    request<BillingRecord>('/api/billing', { method: 'POST', body: JSON.stringify(data) }),
  transitionBilling: (id: string, target_status: string, version: number) =>
    request<BillingRecord>(`/api/billing/${id}/transition`, { method: 'POST', body: JSON.stringify({ target_status, version }) }),

  // Collections
  listPayments: () => request<{ payments: ClientPayment[] }>('/api/payments'),
  recordPayment: (data: { client_id: string; amount_cents: number; currency_code: string; payment_method?: string; reference_number?: string }) =>
    request<ClientPayment>('/api/payments', { method: 'POST', body: JSON.stringify(data) }),
  allocatePayment: (payment_id: string, billing_record_id: string, amount_cents: number) =>
    request<BillingAllocation>(`/api/payments/${payment_id}/allocate`, { method: 'POST', body: JSON.stringify({ billing_record_id, amount_cents }) }),
}

export interface DashboardData {
  kpis: { label: string; value: string; trend: string; trend_direction: string }[]
  activity: { action: string; entity: string; actor: string; timestamp: string }[]
  client_count: number
  quotation_count: number
}

export interface Client {
  id: string
  name: string
  code: string
  contact_email: string
  contact_phone?: string
  address?: string
  is_active: boolean
}

export interface Quotation {
  id: string
  client_id: string
  quotation_number: string
  status: string
  currency_code: string
  subtotal: number
  tax_amount: number
  total: number
  notes: string
  version: number
  lines: QuotationLine[]
}

export interface QuotationLine {
  id: string
  description: string
  quantity: number
  unit_price: number
  line_total: number
}

export interface FundingSource {
  id: string
  name: string
  code: string
  is_approved: boolean
  balance_cents: number
  currency_code: string
}

export interface BudgetRequest {
  id: string
  quotation_id: string
  client_id: string
  request_number: string
  status: string
  currency_code: string
  amount_cents: number
  purpose: string
  version: number
}

export interface Disbursement {
  id: string
  budget_request_id: string
  funding_source_id: string
  status: string
  amount_cents: number
  currency_code: string
  reference_number?: string
  notes?: string
  version: number
}

export interface Liquidation {
  id: string
  disbursement_id: string
  status: string
  released_amount: number
  actual_amount: number
  variance_amount: number
  notes?: string
  version: number
}

export interface BillingRecord {
  id: string
  client_id: string
  billing_number: string
  status: string
  currency_code: string
  total: number
  notes?: string
  version: number
}

export interface ClientPayment {
  id: string
  client_id: string
  payment_number: string
  amount_cents: number
  currency_code: string
  payment_method?: string
  reference_number?: string
  version: number
}

export interface BillingAllocation {
  id: string
  client_payment_id: string
  billing_record_id: string
  amount_cents: number
}

export { type LoginResponse, type User, type Session }
export interface User {
  user_id: string
  username: string
  display_name: string
  email: string
  email_verified: boolean
  role_preview: string
}

export interface Session {
  user_id: string
  username: string
  display_name: string
  email: string
  role_preview: string
}

export interface LoginResponse {
  status: 'authenticated' | 'otp_required'
  display_name?: string
  csrf_token?: string
  message?: string
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

export interface Client {
  id: string
  name: string
  code: string
  contact_email: string
  contact_phone?: string
  address?: string
  is_active: boolean
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
  version: number
}

export interface Liquidation {
  id: string
  disbursement_id: string
  status: string
  released_amount: number
  actual_amount: number
  variance_amount: number
  version: number
}

export interface BillingRecord {
  id: string
  client_id: string
  billing_number: string
  status: string
  currency_code: string
  total: number
  version: number
}

export interface ClientPayment {
  id: string
  client_id: string
  payment_number: string
  amount_cents: number
  currency_code: string
  version: number
}

export interface BillingAllocation {
  id: string
  client_payment_id: string
  billing_record_id: string
  amount_cents: number
}

export const rolePreviewLabels: Record<string, string> = {
  '': 'Full administrator',
  logistics_coordinator: 'Logistics coordinator',
  freight_ops_approver: 'Freight operations approver',
  disbursement_controller: 'Disbursement controller',
  finance_ops_lead: 'Finance operations lead',
}
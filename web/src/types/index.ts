import { type z } from 'zod'

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
  is_active: boolean
}

export const rolePreviewLabels: Record<string, string> = {
  '': 'Full administrator',
  logistics_coordinator: 'Logistics coordinator',
  freight_ops_approver: 'Freight operations approver',
  disbursement_controller: 'Disbursement controller',
  finance_ops_lead: 'Finance operations lead',
}

export type { z }
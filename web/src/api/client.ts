import type { LoginResponse, User, Session } from '../types'

const API_BASE = '/api'

function getCSRFToken(): string | null {
  const match = document.cookie.match(/__Host-dos_csrf=([^;]+)/)
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

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers, credentials: 'include' })

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
}

export { type LoginResponse, type User, type Session }
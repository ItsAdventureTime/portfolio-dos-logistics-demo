import { useState, useEffect } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { api } from './api/client'
import { LoginPage } from './auth/LoginPage'
import { AppShell } from './app/AppShell'
import { DashboardPage } from './dashboard/DashboardPage'
import { QuotationsPage } from './quotations/QuotationsPage'
import { FundingRequestsPage } from './funding/FundingRequestsPage'
import { DisbursementsPage } from './disbursements/DisbursementsPage'
import { LiquidationsPage } from './liquidations/LiquidationsPage'
import { BillingPage } from './billing/BillingPage'
import { CollectionsPage } from './collections/CollectionsPage'
import { DocumentsPage } from './documents/DocumentsPage'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { staleTime: 30000, retry: 1 },
  },
})

type Route = 'dashboard' | 'quotations' | 'funding' | 'disbursements' | 'liquidations' | 'billing' | 'collections' | 'documents'

function getRouteFromHash(): Route {
  const hash = window.location.hash.slice(1)
  const path = hash.split('/')[1] || ''
  if (path === 'quotations') return 'quotations'
  if (path === 'funding') return 'funding'
  if (path === 'disbursements') return 'disbursements'
  if (path === 'liquidations') return 'liquidations'
  if (path === 'billing') return 'billing'
  if (path === 'collections') return 'collections'
  if (path === 'documents') return 'documents'
  return 'dashboard'
}

export default function App() {
  const [authenticated, setAuthenticated] = useState<boolean | null>(null)
  const [route, setRoute] = useState<Route>(getRouteFromHash())

  useEffect(() => {
    const checkAuth = async () => {
      try {
        await api.me()
        setAuthenticated(true)
      } catch {
        setAuthenticated(false)
      }
    }
    checkAuth()
  }, [])

  useEffect(() => {
    const onHashChange = () => setRoute(getRouteFromHash())
    window.addEventListener('hashchange', onHashChange)
    return () => window.removeEventListener('hashchange', onHashChange)
  }, [])

  const handleLogout = async () => {
    try {
      await api.logout()
    } catch {
      // Session may be expired
    }
    setAuthenticated(false)
  }

  if (authenticated === null) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[hsl(var(--ui-background))]">
        <p className="text-[hsl(var(--text-03))]">Loading...</p>
      </div>
    )
  }

  if (!authenticated) {
    return <LoginPage onSuccess={() => setAuthenticated(true)} />
  }

  return (
    <QueryClientProvider client={queryClient}>
      <AppShell onLogout={handleLogout}>
        {route === 'dashboard' && <DashboardPage />}
        {route === 'quotations' && <QuotationsPage />}
        {route === 'funding' && <FundingRequestsPage />}
        {route === 'disbursements' && <DisbursementsPage />}
        {route === 'liquidations' && <LiquidationsPage />}
        {route === 'billing' && <BillingPage />}
        {route === 'collections' && <CollectionsPage />}
        {route === 'documents' && <DocumentsPage />}
      </AppShell>
    </QueryClientProvider>
  )
}
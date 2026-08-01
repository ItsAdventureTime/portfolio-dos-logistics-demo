import { type ReactNode, useState, useEffect, useCallback } from 'react'
import { api } from '../api/client'
import { type User, rolePreviewLabels } from '../types'
import { StatusPill } from '../components/feedback/Badge'
import { Button } from '../components/ui/Button'

interface AppShellProps {
  children: ReactNode
  onLogout: () => void
}

export function AppShell({ children, onLogout }: AppShellProps) {
  const [user, setUser] = useState<User | null>(null)
  const [navOpen, setNavOpen] = useState(false)

  const loadUser = useCallback(async () => {
    try {
      const u = await api.me()
      setUser(u)
    } catch {
      onLogout()
    }
  }, [onLogout])

  useEffect(() => {
    loadUser()
  }, [loadUser])

  const handleRoleChange = async (value: string) => {
    try {
      await api.setRolePreview(value)
      await loadUser()
    } catch (e) {
      console.error(e)
    }
  }

  const navItems = [
    { label: 'Dashboard', path: '#/', icon: '▤' },
    { label: 'Quotations', path: '#/quotations', icon: '⊟' },
    { label: 'Funding requests', path: '#/funding', icon: '⊙' },
    { label: 'Disbursements', path: '#/disbursements', icon: '↗' },
    { label: 'Liquidations', path: '#/liquidations', icon: '⊙' },
    { label: 'Billing', path: '#/billing', icon: '⊞' },
    { label: 'Collections', path: '#/collections', icon: '◆' },
    { label: 'Documents', path: '#/documents', icon: '☰' },
  ]

  const envMarker = 'DEMONSTRATION ENVIRONMENT'

  return (
    <div className="flex min-h-screen flex-col bg-[hsl(var(--surface))]">
      <header className="flex items-center justify-between border-b border-[hsl(var(--border))] bg-[hsl(var(--surface-elevated))] px-4 py-3 lg:px-6">
        <div className="flex items-center gap-3">
          <button
            onClick={() => setNavOpen(!navOpen)}
            className="lg:hidden rounded-lg p-2 text-[hsl(var(--content))] hover:bg-[hsl(var(--surface-panel))] focus:outline-none focus-visible:ring-2 focus-visible:ring-[hsl(var(--focus))]"
            aria-label="Toggle navigation"
            aria-expanded={navOpen}
          >
            <span aria-hidden="true">{navOpen ? '✕' : '☰'}</span>
          </button>
          <span className="text-lg font-semibold text-[hsl(var(--content))]">DOS FreightFlow Control</span>
          <StatusPill status="warning" label={envMarker} />
        </div>
        <div className="flex items-center gap-4">
          {user && (
            <>
              <select
                value={user.role_preview}
                onChange={(e) => handleRoleChange(e.target.value)}
                aria-label="Role preview"
                className="rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--surface))] px-3 py-1.5 text-sm text-[hsl(var(--content))] focus:outline-none focus-visible:ring-2 focus-visible:ring-[hsl(var(--focus))]"
              >
                <option value="">Full administrator</option>
                {Object.entries(rolePreviewLabels)
                  .filter(([k]) => k !== '')
                  .map(([k, label]) => (
                    <option key={k} value={k}>{label}</option>
                  ))}
              </select>
              <span className="hidden text-sm text-[hsl(var(--content-muted))] sm:block">
                {user.display_name}
              </span>
              <Button onPress={onLogout} size="sm">
                Sign out
              </Button>
            </>
          )}
        </div>
      </header>

      <div className="flex flex-1">
        <nav
          className={`${navOpen ? 'block' : 'hidden'} fixed inset-y-0 left-0 z-30 w-64 border-r border-[hsl(var(--border))] bg-[hsl(var(--surface-elevated))] pt-16 lg:static lg:block lg:pt-0`}
          aria-label="Primary navigation"
        >
          <ul className="flex flex-col gap-1 p-4">
            {navItems.map((item) => (
              <li key={item.path}>
                <a
                  href={item.path}
                  className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm text-[hsl(var(--content-muted))] hover:bg-[hsl(var(--surface-panel))] hover:text-[hsl(var(--content))] focus:outline-none focus-visible:ring-2 focus-visible:ring-[hsl(var(--focus))] focus-visible:ring-offset-2 focus-visible:ring-offset-[hsl(var(--surface-elevated))]"
                >
                  <span aria-hidden="true" className="text-base">{item.icon}</span>
                  {item.label}
                </a>
              </li>
            ))}
          </ul>
        </nav>

        <main className="flex-1 overflow-x-hidden p-4 lg:p-6" id="main-content" tabIndex={-1}>
          {children}
        </main>
      </div>
    </div>
  )
}
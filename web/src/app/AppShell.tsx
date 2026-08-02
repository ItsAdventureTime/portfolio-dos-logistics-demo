import { type ReactNode, useState, useEffect, useCallback } from 'react'
import { api } from '../api/client'
import { type User, rolePreviewLabels } from '../types'
import { StatusPill } from '../components/feedback/Badge'
import { Button } from '../components/ui/Button'

interface AppShellProps {
  children: ReactNode
  onLogout: () => void
}

const navSections = [
  {
    label: 'Operations',
    items: [
      { label: 'Dashboard', path: '#/', icon: '\u25A4' },
      { label: 'Quotations', path: '#/quotations', icon: '\u229E' },
      { label: 'Funding requests', path: '#/funding', icon: '\u2299' },
    ],
  },
  {
    label: 'Finance',
    items: [
      { label: 'Disbursements', path: '#/disbursements', icon: '\u2197' },
      { label: 'Liquidations', path: '#/liquidations', icon: '\u2299' },
      { label: 'Billing', path: '#/billing', icon: '\u229E' },
      { label: 'Collections', path: '#/collections', icon: '\u25C6' },
    ],
  },
  {
    label: 'Support',
    items: [
      { label: 'Documents', path: '#/documents', icon: '\u2630' },
    ],
  },
]

export function AppShell({ children, onLogout }: AppShellProps) {
  const [user, setUser] = useState<User | null>(null)
  const [navOpen, setNavOpen] = useState(false)
  const [currentHash, setCurrentHash] = useState(window.location.hash)

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

  useEffect(() => {
    const onHashChange = () => setCurrentHash(window.location.hash)
    window.addEventListener('hashchange', onHashChange)
    return () => window.removeEventListener('hashchange', onHashChange)
  }, [])

  const handleRoleChange = async (value: string) => {
    try {
      await api.setRolePreview(value)
      await loadUser()
    } catch (e) {
      console.error(e)
    }
  }

  const envMarker = 'DEMONSTRATION'

  const isActive = (path: string) => {
    const route = path.replace('#/', '')
    if (route === '') return currentHash === '' || currentHash === '#/' || currentHash === '#'
    return currentHash.startsWith(`#/${route}`)
  }

  return (
    <div className="flex min-h-screen bg-[hsl(var(--ui-background))]">
      {/* Sidebar */}
      <aside
        className={`${navOpen ? 'fixed inset-y-0 left-0 z-40 w-64' : 'hidden'} lg:relative lg:flex lg:w-64 lg:flex-shrink-0`}
        aria-label="Primary navigation"
      >
        <div className="flex w-64 flex-col border-r border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))]">
          {/* DOS brand header */}
          <div className="flex h-14 items-center gap-2.5 border-b border-[hsl(var(--border-subtle))] px-4">
            <img
              src="/dos-mark.png"
              alt="DOS"
              className="h-7 w-7 rounded-[var(--radius-sm)]"
            />
            <div className="flex flex-col leading-tight">
              <span className="text-sm font-semibold text-[hsl(var(--text-01))]">DOS FreightFlow</span>
              <span className="text-[10px] text-[hsl(var(--text-03))]">DelegateOps Business Support</span>
            </div>
          </div>

          {/* Navigation */}
          <nav className="flex-1 overflow-y-auto py-4">
            {navSections.map((section) => (
              <div key={section.label} className="mb-5">
                <h2 className="px-4 pb-1.5 text-[11px] font-semibold uppercase tracking-wider text-[hsl(var(--text-03))]">
                  {section.label}
                </h2>
                <ul className="flex flex-col gap-0.5 px-2">
                  {section.items.map((item) => (
                    <li key={item.path}>
                      <a
                        href={item.path}
                        className={`flex h-9 items-center gap-3 rounded-[var(--radius-sm)] px-3 text-sm transition-all var(--duration-fast-01) var(--ease-productive) focus:outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[hsl(var(--focus))] ${
                          isActive(item.path)
                            ? 'bg-[hsl(var(--interactive-01))]/12 text-[hsl(var(--interactive-04))] font-medium'
                            : 'text-[hsl(var(--text-02))] hover:bg-[hsl(var(--hover-ui))] hover:text-[hsl(var(--text-01))]'
                        }`}
                      >
                        <span aria-hidden="true" className={`text-base ${isActive(item.path) ? 'text-[hsl(var(--interactive-04))]' : 'text-[hsl(var(--text-03))]'}`}>{item.icon}</span>
                        {item.label}
                      </a>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </nav>

          {/* User info at bottom */}
          {user && (
            <div className="border-t border-[hsl(var(--border-subtle))] p-3">
              <div className="flex items-center gap-2.5">
                <div className="flex h-8 w-8 items-center justify-center rounded-full bg-[hsl(var(--interactive-01))] text-xs font-semibold text-[hsl(var(--text-on-color))]">
                  {user.display_name.charAt(0)}
                </div>
                <div className="flex flex-col leading-tight">
                  <span className="text-sm font-medium text-[hsl(var(--text-01))]">{user.display_name}</span>
                  <span className="text-[11px] text-[hsl(var(--text-03))]">{user.email}</span>
                </div>
              </div>
            </div>
          )}
        </div>
      </aside>

      {/* Mobile overlay */}
      {navOpen && (
        <div
          className="fixed inset-0 z-30 bg-black/50 backdrop-blur-sm transition-opacity var(--duration-moderate-02) var(--ease-exit) lg:hidden"
          onClick={() => setNavOpen(false)}
          aria-hidden="true"
        />
      )}

      {/* Main content area */}
      <div className="flex flex-1 flex-col overflow-hidden">
        {/* Top bar */}
        <header className="flex h-14 items-center justify-between border-b border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] px-4 shadow-[var(--shadow-01)] lg:px-6">
          <div className="flex items-center gap-3">
            <button
              onClick={() => setNavOpen(!navOpen)}
              className="lg:hidden flex h-8 w-8 items-center justify-center rounded-[var(--radius-sm)] text-[hsl(var(--text-01))] transition-colors var(--duration-fast-01) var(--ease-productive) hover:bg-[hsl(var(--hover-ui))] focus:outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-[hsl(var(--focus))]"
              aria-label="Toggle navigation"
              aria-expanded={navOpen}
            >
              <span aria-hidden="true">{navOpen ? '\u2715' : '\u2630'}</span>
            </button>
            <StatusPill status="warning" label={envMarker} />
          </div>
          <div className="flex items-center gap-3">
            {user && (
              <>
                <label className="hidden items-center gap-2 sm:flex">
                  <span className="text-xs text-[hsl(var(--text-03))]">Viewing as</span>
                  <select
                    value={user.role_preview}
                    onChange={(e) => handleRoleChange(e.target.value)}
                    aria-label="Role preview"
                    className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-strong))] bg-[hsl(var(--ui-02))] px-2.5 py-1.5 text-xs text-[hsl(var(--text-01))] transition-colors var(--duration-fast-01) var(--ease-productive) focus:outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-[hsl(var(--focus))] hover:border-[hsl(var(--border-interactive))]"
                  >
                    <option value="">Full administrator</option>
                    {Object.entries(rolePreviewLabels)
                      .filter(([k]) => k !== '')
                      .map(([k, label]) => (
                        <option key={k} value={k}>{label}</option>
                      ))}
                  </select>
                </label>
                <Button onPress={onLogout} variant="ghost" size="sm">
                  Sign out
                </Button>
              </>
            )}
          </div>
        </header>

        {/* Page content */}
        <main className="flex-1 overflow-y-auto p-4 lg:p-6" id="main-content" tabIndex={-1}>
          <div
            key={currentHash}
            style={{ animation: 'var(--duration-moderate-01) var(--ease-entrance) fadeIn' }}
          >
            {children}
          </div>
        </main>
      </div>

      <style>{`
        @keyframes fadeIn {
          from { opacity: 0; transform: translateY(4px); }
          to { opacity: 1; transform: translateY(0); }
        }
      `}</style>
    </div>
  )
}
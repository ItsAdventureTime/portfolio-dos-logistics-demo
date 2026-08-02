interface KPI {
  label: string
  value: string
  trend?: string
  trendDirection?: 'up' | 'down' | 'neutral'
}

interface ActivityItem {
  action: string
  entity: string
  actor: string
  timestamp: string
}

const mockKPIs: KPI[] = [
  { label: 'Active quotations', value: '12', trend: '+3 this week', trendDirection: 'up' },
  { label: 'Pending approvals', value: '4', trend: '2 overdue', trendDirection: 'down' },
  { label: 'Outstanding receivables', value: '$284,500', trend: '-$12,000', trendDirection: 'down' },
  { label: 'Disbursed this month', value: '$156,200', trend: '+$8,400', trendDirection: 'up' },
]

const mockActivity: ActivityItem[] = [
  { action: 'Quotation accepted', entity: 'Q-ACME-100001', actor: 'Juan Dela Cruz', timestamp: '2 hours ago' },
  { action: 'Disbursement released', entity: 'DISB-2026-0042', actor: 'Admin', timestamp: '5 hours ago' },
  { action: 'Billing finalized', entity: 'INV-GLOBX-100002', actor: 'Admin', timestamp: '1 day ago' },
  { action: 'Client payment received', entity: 'PMT-STARK-100003', actor: 'Admin', timestamp: '2 days ago' },
]

export function DashboardPage() {
  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Dashboard</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Overview of freight operations and financial controls</p>
      </div>

      {/* KPI strip - Carbon pattern: 4 cards in a row */}
      <section aria-labelledby="kpi-heading">
        <h2 id="kpi-heading" className="sr-only">Key performance indicators</h2>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {mockKPIs.map((kpi) => (
            <div
              key={kpi.label}
              className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-4"
            >
              <p className="text-xs font-medium uppercase tracking-wide text-[hsl(var(--text-03))]">{kpi.label}</p>
              <p className="mt-2 text-2xl font-semibold text-[hsl(var(--text-01))]">{kpi.value}</p>
              {kpi.trend && (
                <p
                  className={`mt-1.5 flex items-center gap-1 text-xs ${kpi.trendDirection === 'up' ? 'text-[hsl(var(--support-success))]' : kpi.trendDirection === 'down' ? 'text-[hsl(var(--support-error))]' : 'text-[hsl(var(--text-03))]'}`}
                >
                  <span aria-hidden="true">{kpi.trendDirection === 'up' ? '\u2191' : kpi.trendDirection === 'down' ? '\u2193' : '\u2014'}</span>
                  {kpi.trend}
                </p>
              )}
            </div>
          ))}
        </div>
      </section>

      {/* Activity feed */}
      <section
        aria-labelledby="activity-heading"
        className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))]"
      >
        <div className="border-b border-[hsl(var(--border-subtle))] px-4 py-3">
          <h2 id="activity-heading" className="text-sm font-semibold text-[hsl(var(--text-01))]">Recent activity</h2>
        </div>
        <ul className="flex flex-col">
          {mockActivity.map((item, i) => (
            <li
              key={i}
              className={`flex items-center justify-between px-4 py-3 ${i < mockActivity.length - 1 ? 'border-b border-[hsl(var(--border-subtle))]/60' : ''}`}
            >
              <div className="flex flex-col gap-0.5">
                <p className="text-sm text-[hsl(var(--text-01))]">{item.action}</p>
                <p className="text-xs text-[hsl(var(--text-03))]">{item.entity} \u00B7 {item.actor}</p>
              </div>
              <span className="text-xs whitespace-nowrap text-[hsl(var(--text-03))]">{item.timestamp}</span>
            </li>
          ))}
        </ul>
      </section>
    </div>
  )
}
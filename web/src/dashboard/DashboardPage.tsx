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
        <h1 className="text-xl font-semibold text-[hsl(var(--content))]">Dashboard</h1>
        <p className="mt-1 text-sm text-[hsl(var(--content-muted))]">Overview of freight operations and financial controls</p>
      </div>

      <section aria-labelledby="kpi-heading">
        <h2 id="kpi-heading" className="sr-only">Key performance indicators</h2>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {mockKPIs.map((kpi) => (
            <div
              key={kpi.label}
              className="rounded-xl border border-[hsl(var(--border))] bg-[hsl(var(--surface-elevated))] p-5"
            >
              <p className="text-sm text-[hsl(var(--content-muted))]">{kpi.label}</p>
              <p className="mt-2 text-2xl font-semibold text-[hsl(var(--content))]">{kpi.value}</p>
              {kpi.trend && (
                <p
                  className={`mt-1 text-xs ${kpi.trendDirection === 'up' ? 'text-[hsl(var(--status-success))]' : kpi.trendDirection === 'down' ? 'text-[hsl(var(--status-error))]' : 'text-[hsl(var(--content-muted))]'}`}
                >
                  {kpi.trend}
                </p>
              )}
            </div>
          ))}
        </div>
      </section>

      <section aria-labelledby="activity-heading" className="rounded-xl border border-[hsl(var(--border))] bg-[hsl(var(--surface-elevated))] p-5">
        <h2 id="activity-heading" className="mb-4 text-lg font-medium text-[hsl(var(--content))]">Recent activity</h2>
        <ul className="flex flex-col gap-3">
          {mockActivity.map((item, i) => (
            <li key={i} className="flex items-center justify-between border-b border-[hsl(var(--border))]/50 pb-3 last:border-0 last:pb-0">
              <div>
                <p className="text-sm text-[hsl(var(--content))]">{item.action}</p>
                <p className="text-xs text-[hsl(var(--content-muted))]">{item.entity} · {item.actor}</p>
              </div>
              <span className="text-xs text-[hsl(var(--content-muted))]">{item.timestamp}</span>
            </li>
          ))}
        </ul>
      </section>
    </div>
  )
}
import { useQuery } from '@tanstack/react-query'
import { api, type DashboardData } from '../api/client'

export function DashboardPage() {
  const { data, isLoading, isError } = useQuery<DashboardData>({
    queryKey: ['dashboard'],
    queryFn: api.getDashboard,
  })

  if (isLoading) {
    return (
      <div className="flex flex-col gap-6">
        <div>
          <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Dashboard</h1>
          <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Overview of freight operations and financial controls</p>
        </div>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {[0, 1, 2, 3].map((i) => (
            <div key={i} className="h-28 animate-pulse rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-02))]" />
          ))}
        </div>
      </div>
    )
  }

  if (isError || !data) {
    return (
      <div className="flex flex-col gap-6">
        <div>
          <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Dashboard</h1>
        </div>
        <div className="rounded-[var(--radius-sm)] border border-[hsl(var(--support-error))]/30 bg-[hsl(var(--support-error))]/10 p-4 text-sm text-[hsl(var(--support-error))]">
          Unable to load dashboard data. Please try again.
        </div>
      </div>
    )
  }

  const kpis = data.kpis || []
  const activity = data.activity || []

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Dashboard</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Overview of freight operations and financial controls</p>
      </div>

      {/* KPI strip */}
      <section aria-labelledby="kpi-heading">
        <h2 id="kpi-heading" className="sr-only">Key performance indicators</h2>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {kpis.map((kpi, i) => (
            <div
              key={i}
              className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-4"
            >
              <p className="text-xs font-medium uppercase tracking-wide text-[hsl(var(--text-03))]">{kpi.label}</p>
              <p className="mt-2 text-2xl font-semibold text-[hsl(var(--text-01))]">{kpi.value}</p>
              {kpi.trend && (
                <p
                  className={`mt-1.5 flex items-center gap-1 text-xs ${kpi.trend_direction === 'up' ? 'text-[hsl(var(--support-success))]' : kpi.trend_direction === 'down' ? 'text-[hsl(var(--support-error))]' : 'text-[hsl(var(--text-03))]'}`}
                >
                  <span aria-hidden="true">
                    {kpi.trend_direction === 'up' ? '\u2191' : kpi.trend_direction === 'down' ? '\u2193' : '\u2014'}
                  </span>
                  {kpi.trend}
                </p>
              )}
            </div>
          ))}
        </div>
      </section>

      {/* Activity feed */}
      {activity.length > 0 && (
        <section
          aria-labelledby="activity-heading"
          className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))]"
        >
          <div className="border-b border-[hsl(var(--border-subtle))] px-4 py-3">
            <h2 id="activity-heading" className="text-sm font-semibold text-[hsl(var(--text-01))]">Recent activity</h2>
          </div>
          <ul className="flex flex-col">
            {activity.map((item, i) => (
              <li
                key={i}
                className={`flex items-center justify-between px-4 py-3 ${i < activity.length - 1 ? 'border-b border-[hsl(var(--border-subtle))]/60' : ''}`}
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
      )}

      {/* Summary stats */}
      <div className="flex gap-4 text-sm text-[hsl(var(--text-03))]">
        <span>{data.client_count} clients</span>
        <span>{data.quotation_count} quotations</span>
      </div>
    </div>
  )
}
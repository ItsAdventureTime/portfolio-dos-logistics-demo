import { useQuery } from '@tanstack/react-query'
import { api, type Liquidation } from '../api/client'
import { DataTable } from '../components/table/DataTable'
import { Badge, type StatusType } from '../components/feedback/Badge'
import { DetailPanel, type DocumentItem, type TimelineEvent } from '../components/feedback/DetailPanel'
import { type ColumnDef } from '@tanstack/react-table'

const statusMap: Record<string, { status: StatusType; label: string }> = {
  open: { status: 'neutral', label: 'Open' },
  reconciled: { status: 'info', label: 'Reconciled' },
  closed: { status: 'success', label: 'Closed' },
}

function formatCurrency(cents: number): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(cents / 100)
}

function simulatedDocs(status: string): DocumentItem[] {
  if (status === 'open') return []
  return [
    { name: 'Liquidation_Report.pdf', type: 'PDF', size: '210 KB' },
    { name: 'Actual_Spending_Receipts.zip', type: 'DOC', size: '2.4 MB' },
    { name: 'Variance_Explanation.pdf', type: 'PDF', size: '45 KB' },
  ]
}

function simulatedTimeline(status: string): TimelineEvent[] {
  return [
    { action: 'Liquidation created', actor: 'Demo Administrator', timestamp: '2026-07-31 09:00' },
    ...(status !== 'open' ? [{ action: 'Reconciled with actual spending', actor: 'Demo Administrator', timestamp: '2026-07-31 15:00' }] : []),
    ...(status === 'closed' ? [{ action: 'Liquidation closed', actor: 'Demo Administrator', timestamp: '2026-08-01 10:00' }] : []),
  ]
}

const columns: ColumnDef<Liquidation, unknown>[] = [
  { header: 'Status', accessorKey: 'status', cell: ({ row }) => { const c = statusMap[row.original.status] || statusMap.open; return <Badge status={c.status}>{c.label}</Badge> } },
  { header: 'Released', accessorKey: 'released_amount', cell: ({ row }) => <span className="text-right text-[hsl(var(--text-01))]">{formatCurrency(row.original.released_amount)}</span> },
  { header: 'Actual', accessorKey: 'actual_amount', cell: ({ row }) => <span className="text-right text-[hsl(var(--text-01))]">{formatCurrency(row.original.actual_amount)}</span> },
  { header: 'Variance', accessorKey: 'variance_amount', cell: ({ row }) => <span className={`text-right ${row.original.variance_amount > 0 ? 'text-[hsl(var(--support-warning))]' : 'text-[hsl(var(--text-01))]'}`}>{formatCurrency(row.original.variance_amount)}</span> },
  { header: 'Version', accessorKey: 'version', cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">v{row.original.version}</span> },
]

export function LiquidationsPage() {
  const { data, isLoading, isError } = useQuery<{ liquidations: Liquidation[] }>({
    queryKey: ['liquidations'],
    queryFn: () => api.listLiquidations(),
  })

  const items = data?.liquidations || []

  if (isLoading) return <div className="h-64 animate-pulse rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-02))]" />
  if (isError) return <div className="rounded-[var(--radius-sm)] border border-[hsl(var(--support-error))]/30 bg-[hsl(var(--support-error))]/10 p-4 text-sm text-[hsl(var(--support-error))]">Unable to load liquidations.</div>

  const renderDetail = (l: Liquidation) => (
    <DetailPanel
      title="Liquidation detail"
      subtitle={l.notes || 'No notes'}
      documents={simulatedDocs(l.status)}
      timeline={simulatedTimeline(l.status)}
      metadata={[
        { label: 'Status', value: statusMap[l.status]?.label || l.status },
        { label: 'Released', value: formatCurrency(l.released_amount) },
        { label: 'Actual', value: formatCurrency(l.actual_amount) },
        { label: 'Variance', value: formatCurrency(l.variance_amount) },
      ]}
    />
  )

  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Liquidations</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Reconcile released funds with actual spending. Click a row to see receipts and variance reports.</p>
      </div>
      <div className="hidden sm:block">
        <DataTable columns={columns} data={items} caption="Liquidations" emptyMessage="No liquidations found." zebra density="normal" renderDetail={renderDetail} />
      </div>
      <div className="flex flex-col gap-3 sm:hidden">
        {items.length === 0 ? <p className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-6 text-center text-sm text-[hsl(var(--text-03))]">No liquidations found.</p> : items.map((l) => { const c = statusMap[l.status] || statusMap.open; return (
          <details key={l.id} className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))]">
            <summary className="cursor-pointer p-4">
              <div className="flex items-center justify-between"><Badge status={c.status}>{c.label}</Badge><span className="text-[hsl(var(--text-03))]">v{l.version}</span></div>
              <div className="mt-3 flex flex-col gap-1 text-sm">
                <div className="flex justify-between"><span className="text-[hsl(var(--text-03))]">Released</span><span className="text-[hsl(var(--text-01))]">{formatCurrency(l.released_amount)}</span></div>
                <div className="flex justify-between"><span className="text-[hsl(var(--text-03))]">Actual</span><span className="text-[hsl(var(--text-01))]">{formatCurrency(l.actual_amount)}</span></div>
                <div className="flex justify-between"><span className="text-[hsl(var(--text-03))]">Variance</span><span className={l.variance_amount > 0 ? 'text-[hsl(var(--support-warning))]' : 'text-[hsl(var(--text-01))]'}>{formatCurrency(l.variance_amount)}</span></div>
              </div>
            </summary>
            <div className="border-t border-[hsl(var(--border-subtle))] p-4">{renderDetail(l)}</div>
          </details>
        ) })}
      </div>
    </div>
  )
}
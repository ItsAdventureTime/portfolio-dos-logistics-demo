import { DataTable } from '../components/table/DataTable'
import { Badge, type StatusType } from '../components/feedback/Badge'
import { type ColumnDef } from '@tanstack/react-table'

interface Liquidation {
  id: string
  status: string
  released_amount: number
  actual_amount: number
  variance_amount: number
  version: number
}

const statusMap: Record<string, { status: StatusType; label: string }> = {
  open: { status: 'neutral', label: 'Open' },
  reconciled: { status: 'info', label: 'Reconciled' },
  closed: { status: 'success', label: 'Closed' },
}

function formatCurrency(cents: number): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(cents / 100)
}

const columns: ColumnDef<Liquidation, unknown>[] = [
  { header: 'Status', accessorKey: 'status', cell: ({ row }) => { const c = statusMap[row.original.status] || statusMap.open; return <Badge status={c.status}>{c.label}</Badge> } },
  { header: 'Released', accessorKey: 'released_amount', cell: ({ row }) => <span className="text-right text-[hsl(var(--text-01))]">{formatCurrency(row.original.released_amount)}</span> },
  { header: 'Actual', accessorKey: 'actual_amount', cell: ({ row }) => <span className="text-right text-[hsl(var(--text-01))]">{formatCurrency(row.original.actual_amount)}</span> },
  { header: 'Variance', accessorKey: 'variance_amount', cell: ({ row }) => <span className={`text-right ${row.original.variance_amount > 0 ? 'text-[hsl(var(--support-warning))]' : 'text-[hsl(var(--text-01))]'}`}>{formatCurrency(row.original.variance_amount)}</span> },
  { header: 'Version', accessorKey: 'version', cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">v{row.original.version}</span> },
]

const mockData: Liquidation[] = [
  { id: 'demo-liq-001', status: 'open', released_amount: 4250000, actual_amount: 0, variance_amount: 4250000, version: 1 },
]

export function LiquidationsPage() {
  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Liquidations</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Reconcile released funds with actual spending</p>
      </div>
      <div className="hidden sm:block">
        <DataTable columns={columns} data={mockData} caption="Liquidations" emptyMessage="No liquidations found." zebra density="normal" />
      </div>
      <div className="flex flex-col gap-3 sm:hidden">
        {mockData.map((l) => { const c = statusMap[l.status] || statusMap.open; return (
          <div key={l.id} className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-4">
            <div className="flex items-center justify-between"><Badge status={c.status}>{c.label}</Badge><span className="text-[hsl(var(--text-03))]">v{l.version}</span></div>
            <div className="mt-3 flex flex-col gap-1 text-sm">
              <div className="flex justify-between"><span className="text-[hsl(var(--text-03))]">Released</span><span className="text-[hsl(var(--text-01))]">{formatCurrency(l.released_amount)}</span></div>
              <div className="flex justify-between"><span className="text-[hsl(var(--text-03))]">Actual</span><span className="text-[hsl(var(--text-01))]">{formatCurrency(l.actual_amount)}</span></div>
              <div className="flex justify-between"><span className="text-[hsl(var(--text-03))]">Variance</span><span className={l.variance_amount > 0 ? 'text-[hsl(var(--support-warning))]' : 'text-[hsl(var(--text-01))]'}>{formatCurrency(l.variance_amount)}</span></div>
            </div>
          </div>
        ) })}
      </div>
    </div>
  )
}
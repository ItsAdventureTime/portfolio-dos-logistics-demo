import { DataTable } from '../components/table/DataTable'
import { Badge, type StatusType } from '../components/feedback/Badge'
import { type ColumnDef } from '@tanstack/react-table'

interface Disbursement {
  id: string
  status: string
  amount_cents: number
  currency_code: string
  version: number
}

const statusMap: Record<string, { status: StatusType; label: string }> = {
  pending: { status: 'neutral', label: 'Pending' },
  released: { status: 'success', label: 'Released' },
  held: { status: 'warning', label: 'Held' },
  returned: { status: 'error', label: 'Returned' },
}

function formatCurrency(cents: number, code: string): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: code }).format(cents / 100)
}

const columns: ColumnDef<Disbursement, unknown>[] = [
  { header: 'Status', accessorKey: 'status', cell: ({ row }) => { const c = statusMap[row.original.status] || statusMap.pending; return <Badge status={c.status}>{c.label}</Badge> } },
  { header: 'Amount', accessorKey: 'amount_cents', cell: ({ row }) => <span className="text-right text-[hsl(var(--text-01))]">{formatCurrency(row.original.amount_cents, row.original.currency_code)}</span> },
  { header: 'Version', accessorKey: 'version', cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">v{row.original.version}</span> },
]

const mockData: Disbursement[] = [
  { id: 'demo-disb-001', status: 'released', amount_cents: 4250000, currency_code: 'USD', version: 2 },
]

export function DisbursementsPage() {
  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Disbursements</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Released funds from approved budget requests</p>
      </div>
      <div className="hidden sm:block">
        <DataTable columns={columns} data={mockData} caption="Disbursements" emptyMessage="No disbursements found." zebra density="normal" />
      </div>
      <div className="flex flex-col gap-3 sm:hidden">
        {mockData.map((d) => { const c = statusMap[d.status] || statusMap.pending; return (
          <div key={d.id} className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-4">
            <div className="flex items-center justify-between"><span className="font-medium text-[hsl(var(--text-01))]">{formatCurrency(d.amount_cents, d.currency_code)}</span><Badge status={c.status}>{c.label}</Badge></div>
            <p className="mt-2 text-xs text-[hsl(var(--text-03))]">v{d.version}</p>
          </div>
        ) })}
      </div>
    </div>
  )
}
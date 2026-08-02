import { useQuery } from '@tanstack/react-query'
import { api, type Disbursement } from '../api/client'
import { DataTable } from '../components/table/DataTable'
import { Badge, type StatusType } from '../components/feedback/Badge'
import { type ColumnDef } from '@tanstack/react-table'

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
  { header: 'Reference', accessorKey: 'reference_number', cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">{row.original.reference_number || '-'}</span> },
  { header: 'Version', accessorKey: 'version', cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">v{row.original.version}</span> },
]

export function DisbursementsPage() {
  const { data, isLoading, isError } = useQuery<{ disbursements: Disbursement[] }>({
    queryKey: ['disbursements'],
    queryFn: () => api.listDisbursements(),
  })

  const items = data?.disbursements || []

  if (isLoading) return <div className="h-64 animate-pulse rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-02))]" />
  if (isError) return <div className="rounded-[var(--radius-sm)] border border-[hsl(var(--support-error))]/30 bg-[hsl(var(--support-error))]/10 p-4 text-sm text-[hsl(var(--support-error))]">Unable to load disbursements.</div>

  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Disbursements</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Released funds from approved budget requests</p>
      </div>
      <div className="hidden sm:block">
        <DataTable columns={columns} data={items} caption="Disbursements" emptyMessage="No disbursements found." zebra density="normal" />
      </div>
      <div className="flex flex-col gap-3 sm:hidden">
        {items.length === 0 ? <p className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-6 text-center text-sm text-[hsl(var(--text-03))]">No disbursements found.</p> : items.map((d) => { const c = statusMap[d.status] || statusMap.pending; return (
          <div key={d.id} className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-4">
            <div className="flex items-center justify-between"><Badge status={c.status}>{c.label}</Badge><span className="text-[hsl(var(--text-03))]">v{d.version}</span></div>
            <p className="mt-2 text-sm text-[hsl(var(--text-01))]">{formatCurrency(d.amount_cents, d.currency_code)}</p>
            {d.reference_number && <p className="mt-1 text-xs text-[hsl(var(--text-03))]">Ref: {d.reference_number}</p>}
          </div>
        ) })}
      </div>
    </div>
  )
}
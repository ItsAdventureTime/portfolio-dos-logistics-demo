import { useQuery } from '@tanstack/react-query'
import { api, type BudgetRequest } from '../api/client'
import { DataTable } from '../components/table/DataTable'
import { Badge, type StatusType } from '../components/feedback/Badge'
import { type ColumnDef } from '@tanstack/react-table'

const statusMap: Record<string, { status: StatusType; label: string }> = {
  draft: { status: 'neutral', label: 'Draft' },
  controls_review: { status: 'info', label: 'In review' },
  approved: { status: 'success', label: 'Approved' },
  rejected: { status: 'error', label: 'Rejected' },
  returned: { status: 'warning', label: 'Returned' },
}

function formatCurrency(cents: number, code: string): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: code }).format(cents / 100)
}

const columns: ColumnDef<BudgetRequest, unknown>[] = [
  { header: 'Number', accessorKey: 'request_number', cell: ({ row }) => <span className="font-medium text-[hsl(var(--text-01))]">{row.original.request_number}</span> },
  { header: 'Status', accessorKey: 'status', cell: ({ row }) => { const c = statusMap[row.original.status] || statusMap.draft; return <Badge status={c.status}>{c.label}</Badge> } },
  { header: 'Purpose', accessorKey: 'purpose', cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">{row.original.purpose}</span> },
  { header: 'Amount', accessorKey: 'amount_cents', cell: ({ row }) => <span className="text-right text-[hsl(var(--text-01))]">{formatCurrency(row.original.amount_cents, row.original.currency_code)}</span> },
  { header: 'Version', accessorKey: 'version', cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">v{row.original.version}</span> },
]

export function FundingRequestsPage() {
  const { data, isLoading, isError } = useQuery<{ budget_requests: BudgetRequest[] }>({
    queryKey: ['budget-requests'],
    queryFn: () => api.listBudgetRequests(),
  })

  const items = data?.budget_requests || []

  if (isLoading) return <div className="h-64 animate-pulse rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-02))]" />
  if (isError) return <div className="rounded-[var(--radius-sm)] border border-[hsl(var(--support-error))]/30 bg-[hsl(var(--support-error))]/10 p-4 text-sm text-[hsl(var(--support-error))]">Unable to load funding requests.</div>

  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Funding requests</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Shipment funding linked to accepted quotations</p>
      </div>
      <div className="hidden sm:block">
        <DataTable columns={columns} data={items} caption="Budget requests" emptyMessage="No funding requests found." zebra density="normal" />
      </div>
      <div className="flex flex-col gap-3 sm:hidden">
        {items.length === 0 ? <p className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-6 text-center text-sm text-[hsl(var(--text-03))]">No funding requests found.</p> : items.map((br) => { const c = statusMap[br.status] || statusMap.draft; return (
          <div key={br.id} className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-4">
            <div className="flex items-center justify-between"><span className="font-medium text-[hsl(var(--text-01))]">{br.request_number}</span><Badge status={c.status}>{c.label}</Badge></div>
            <p className="mt-2 text-sm text-[hsl(var(--text-03))]">{br.purpose}</p>
            <div className="mt-3 flex items-center justify-between text-sm"><span className="text-[hsl(var(--text-01))]">{formatCurrency(br.amount_cents, br.currency_code)}</span><span className="text-[hsl(var(--text-03))]">v{br.version}</span></div>
          </div>
        ) })}
      </div>
    </div>
  )
}
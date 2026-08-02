import { DataTable } from '../components/table/DataTable'
import { Badge, type StatusType } from '../components/feedback/Badge'
import { type ColumnDef } from '@tanstack/react-table'
import type { BudgetRequest } from '../api/client'

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
  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Funding requests</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Shipment funding linked to accepted quotations</p>
      </div>
      <div className="hidden sm:block">
        <DataTable columns={columns} data={[]} caption="Budget requests" emptyMessage="No funding requests found. Create one from an accepted quotation." zebra density="normal" />
      </div>
      <p className="text-sm text-[hsl(var(--text-03))] sm:hidden">No funding requests found. Create one from an accepted quotation.</p>
    </div>
  )
}
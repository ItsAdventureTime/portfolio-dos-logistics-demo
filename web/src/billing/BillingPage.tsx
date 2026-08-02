import { DataTable } from '../components/table/DataTable'
import { Badge, type StatusType } from '../components/feedback/Badge'
import { type ColumnDef } from '@tanstack/react-table'

interface BillingRecord {
  id: string
  billing_number: string
  status: string
  currency_code: string
  total: number
  version: number
}

const statusMap: Record<string, { status: StatusType; label: string }> = {
  draft: { status: 'neutral', label: 'Draft' },
  review: { status: 'info', label: 'In review' },
  approved: { status: 'info', label: 'Approved' },
  finalized: { status: 'success', label: 'Finalized' },
  void: { status: 'error', label: 'Void' },
  replaced: { status: 'warning', label: 'Replaced' },
}

function formatCurrency(cents: number, code: string): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: code }).format(cents / 100)
}

const columns: ColumnDef<BillingRecord, unknown>[] = [
  { header: 'Number', accessorKey: 'billing_number', cell: ({ row }) => <span className="font-medium text-[hsl(var(--text-01))]">{row.original.billing_number}</span> },
  { header: 'Status', accessorKey: 'status', cell: ({ row }) => { const c = statusMap[row.original.status] || statusMap.draft; return <Badge status={c.status}>{c.label}</Badge> } },
  { header: 'Total', accessorKey: 'total', cell: ({ row }) => <span className="text-right font-medium text-[hsl(var(--text-01))]">{formatCurrency(row.original.total, row.original.currency_code)}</span> },
  { header: 'Version', accessorKey: 'version', cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">v{row.original.version}</span> },
]

const mockData: BillingRecord[] = [
  { id: 'demo-bill-001', billing_number: 'INV-ACME-100001', status: 'finalized', currency_code: 'USD', total: 4250000, version: 4 },
]

export function BillingPage() {
  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Billing</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Manage invoices and billing records</p>
      </div>
      <div className="hidden sm:block">
        <DataTable columns={columns} data={mockData} caption="Billing records" emptyMessage="No billing records found." zebra density="normal" />
      </div>
      <div className="flex flex-col gap-3 sm:hidden">
        {mockData.map((b) => { const c = statusMap[b.status] || statusMap.draft; return (
          <div key={b.id} className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-4">
            <div className="flex items-center justify-between"><span className="font-medium text-[hsl(var(--text-01))]">{b.billing_number}</span><Badge status={c.status}>{c.label}</Badge></div>
            <div className="mt-3 flex items-center justify-between text-sm"><span className="text-[hsl(var(--text-01))]">{formatCurrency(b.total, b.currency_code)}</span><span className="text-[hsl(var(--text-03))]">v{b.version}</span></div>
          </div>
        ) })}
      </div>
    </div>
  )
}
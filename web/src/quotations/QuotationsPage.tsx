import { DataTable } from '../components/table/DataTable'
import { Badge, type StatusType } from '../components/feedback/Badge'
import { type ColumnDef } from '@tanstack/react-table'
import { type Quotation } from '../types'

const quotationStatusMap: Record<string, { status: StatusType; label: string }> = {
  draft: { status: 'neutral', label: 'Draft' },
  review: { status: 'info', label: 'In review' },
  approved: { status: 'info', label: 'Approved' },
  accepted: { status: 'success', label: 'Accepted' },
  revised: { status: 'warning', label: 'Revised' },
  void: { status: 'error', label: 'Void' },
}

function formatCurrency(cents: number, code: string): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: code }).format(cents / 100)
}

const columns: ColumnDef<Quotation, unknown>[] = [
  {
    header: 'Number',
    accessorKey: 'quotation_number',
    cell: ({ row }) => (
      <span className="font-medium text-[hsl(var(--text-01))]">{row.original.quotation_number}</span>
    ),
  },
  {
    header: 'Status',
    accessorKey: 'status',
    cell: ({ row }) => {
      const config = quotationStatusMap[row.original.status] || quotationStatusMap.draft
      return <Badge status={config.status}>{config.label}</Badge>
    },
  },
  {
    header: 'Currency',
    accessorKey: 'currency_code',
    cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">{row.original.currency_code}</span>,
  },
  {
    header: 'Total',
    accessorKey: 'total',
    cell: ({ row }) => (
      <span className="text-right font-medium text-[hsl(var(--text-01))]">
        {formatCurrency(row.original.total, row.original.currency_code)}
      </span>
    ),
  },
  {
    header: 'Version',
    accessorKey: 'version',
    cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">v{row.original.version}</span>,
  },
]

const mockQuotations: Quotation[] = [
  {
    id: 'demo-quot-001',
    client_id: 'demo-client-acme',
    quotation_number: 'Q-ACME-100001',
    status: 'accepted',
    currency_code: 'USD',
    subtotal: 4250000,
    tax_amount: 0,
    total: 4250000,
    notes: 'Ocean freight: 40ft container, Manila to LA',
    version: 3,
    lines: [],
  },
  {
    id: 'demo-quot-002',
    client_id: 'demo-client-globex',
    quotation_number: 'Q-GLOBX-100002',
    status: 'draft',
    currency_code: 'USD',
    subtotal: 1850000,
    tax_amount: 0,
    total: 1850000,
    notes: 'Air freight: JFK to Heathrow',
    version: 1,
    lines: [],
  },
]

export function QuotationsPage() {
  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Quotations</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Create, review, and manage freight quotations</p>
      </div>

      {/* Desktop: TanStack Table */}
      <div className="hidden sm:block">
        <DataTable
          columns={columns}
          data={mockQuotations}
          caption="Freight quotations"
          emptyMessage="No quotations found. Create one to get started."
          zebra={true}
          density="normal"
        />
      </div>

      {/* Mobile: card-based alternative */}
      <div className="flex flex-col gap-3 sm:hidden">
        {mockQuotations.length === 0 ? (
          <p className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-6 text-center text-sm text-[hsl(var(--text-03))]">
            No quotations found. Create one to get started.
          </p>
        ) : (
          mockQuotations.map((q) => {
            const config = quotationStatusMap[q.status] || quotationStatusMap.draft
            return (
              <div
                key={q.id}
                className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-4"
              >
                <div className="flex items-center justify-between">
                  <span className="font-medium text-[hsl(var(--text-01))]">{q.quotation_number}</span>
                  <Badge status={config.status}>{config.label}</Badge>
                </div>
                <p className="mt-2 text-sm text-[hsl(var(--text-03))]">{q.notes}</p>
                <div className="mt-3 flex items-center justify-between text-sm">
                  <span className="text-[hsl(var(--text-01))]">{formatCurrency(q.total, q.currency_code)}</span>
                  <span className="text-[hsl(var(--text-03))]">v{q.version}</span>
                </div>
              </div>
            )
          })
        )}
      </div>
    </div>
  )
}
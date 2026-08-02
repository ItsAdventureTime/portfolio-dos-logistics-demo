import { useQuery } from '@tanstack/react-query'
import { DataTable } from '../components/table/DataTable'
import { Badge, type StatusType } from '../components/feedback/Badge'
import { type ColumnDef } from '@tanstack/react-table'
import { api, type Quotation } from '../api/client'

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

export function QuotationsPage() {
  const { data, isLoading, isError } = useQuery<{ quotations: Quotation[] }>({
    queryKey: ['quotations'],
    queryFn: () => api.listQuotations(),
  })

  const quotations = data?.quotations || []

  if (isLoading) {
    return (
      <div className="flex flex-col gap-5">
        <div>
          <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Quotations</h1>
          <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Create, review, and manage freight quotations</p>
        </div>
        <div className="h-64 animate-pulse rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-02))]" />
      </div>
    )
  }

  if (isError) {
    return (
      <div className="flex flex-col gap-5">
        <div>
          <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Quotations</h1>
        </div>
        <div className="rounded-[var(--radius-sm)] border border-[hsl(var(--support-error))]/30 bg-[hsl(var(--support-error))]/10 p-4 text-sm text-[hsl(var(--support-error))]">
          Unable to load quotations. Please try again.
        </div>
      </div>
    )
  }

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
          data={quotations}
          caption="Freight quotations"
          emptyMessage="No quotations found. Create one to get started."
          zebra={true}
          density="normal"
        />
      </div>

      {/* Mobile: card-based alternative */}
      <div className="flex flex-col gap-3 sm:hidden">
        {quotations.length === 0 ? (
          <p className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-6 text-center text-sm text-[hsl(var(--text-03))]">
            No quotations found. Create one to get started.
          </p>
        ) : (
          quotations.map((q) => {
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
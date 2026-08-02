import { DataTable } from '../components/table/DataTable'
import { type ColumnDef } from '@tanstack/react-table'

interface ClientPayment {
  id: string
  payment_number: string
  amount_cents: number
  currency_code: string
  version: number
}

function formatCurrency(cents: number, code: string): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: code }).format(cents / 100)
}

const columns: ColumnDef<ClientPayment, unknown>[] = [
  { header: 'Number', accessorKey: 'payment_number', cell: ({ row }) => <span className="font-medium text-[hsl(var(--text-01))]">{row.original.payment_number}</span> },
  { header: 'Amount', accessorKey: 'amount_cents', cell: ({ row }) => <span className="text-right font-medium text-[hsl(var(--text-01))]">{formatCurrency(row.original.amount_cents, row.original.currency_code)}</span> },
  { header: 'Version', accessorKey: 'version', cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">v{row.original.version}</span> },
]

const mockData: ClientPayment[] = [
  { id: 'demo-pay-001', payment_number: 'PMT-ACME-100001', amount_cents: 2000000, currency_code: 'USD', version: 1 },
]

export function CollectionsPage() {
  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Collections</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Client payments and billing allocations</p>
      </div>
      <div className="hidden sm:block">
        <DataTable columns={columns} data={mockData} caption="Client payments" emptyMessage="No payments recorded." zebra density="normal" />
      </div>
      <div className="flex flex-col gap-3 sm:hidden">
        {mockData.map((p) => (
          <div key={p.id} className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-4">
            <div className="flex items-center justify-between"><span className="font-medium text-[hsl(var(--text-01))]">{p.payment_number}</span><span className="text-[hsl(var(--text-03))]">v{p.version}</span></div>
            <p className="mt-2 text-sm text-[hsl(var(--text-01))]">{formatCurrency(p.amount_cents, p.currency_code)}</p>
          </div>
        ))}
      </div>
    </div>
  )
}
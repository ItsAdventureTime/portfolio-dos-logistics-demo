import { useQuery } from '@tanstack/react-query'
import { api, type ClientPayment } from '../api/client'
import { DataTable } from '../components/table/DataTable'
import { DetailPanel, type DocumentItem, type TimelineEvent } from '../components/feedback/DetailPanel'
import { type ColumnDef } from '@tanstack/react-table'

function formatCurrency(cents: number, code: string): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: code }).format(cents / 100)
}

function simulatedDocs(payId: string): DocumentItem[] {
  const docs: Record<string, DocumentItem[]> = {
    'demo-pay-001': [
      { name: 'Wire_Confirmation_ACME_001.pdf', type: 'PDF', size: '42 KB' },
      { name: 'Bank_Statement_July.pdf', type: 'PDF', size: '180 KB' },
    ],
    'demo-pay-002': [
      { name: 'ACH_Confirmation_ACME_002.pdf', type: 'PDF', size: '35 KB' },
    ],
    'demo-pay-003': [
      { name: 'Wire_Confirmation_UMBRA_003.pdf', type: 'PDF', size: '48 KB' },
      { name: 'Remittance_Advice_UMBRA.pdf', type: 'PDF', size: '62 KB' },
    ],
  }
  return docs[payId] || []
}

function simulatedTimeline(): TimelineEvent[] {
  return [
    { action: 'Payment received', actor: 'Demo Administrator', timestamp: '2026-08-01 09:00' },
    { action: 'Payment verified', actor: 'System', timestamp: '2026-08-01 09:05' },
    { action: 'Allocated to billing records', actor: 'Demo Administrator', timestamp: '2026-08-01 10:00' },
  ]
}

const columns: ColumnDef<ClientPayment, unknown>[] = [
  { header: 'Number', accessorKey: 'payment_number', cell: ({ row }) => <span className="font-medium text-[hsl(var(--text-01))]">{row.original.payment_number}</span> },
  { header: 'Amount', accessorKey: 'amount_cents', cell: ({ row }) => <span className="text-right font-medium text-[hsl(var(--text-01))]">{formatCurrency(row.original.amount_cents, row.original.currency_code)}</span> },
  { header: 'Version', accessorKey: 'version', cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">v{row.original.version}</span> },
]

export function CollectionsPage() {
  const { data, isLoading, isError } = useQuery<{ payments: ClientPayment[] }>({
    queryKey: ['payments'],
    queryFn: () => api.listPayments(),
  })

  const items = data?.payments || []

  if (isLoading) return <div className="h-64 animate-pulse rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-02))]" />
  if (isError) return <div className="rounded-[var(--radius-sm)] border border-[hsl(var(--support-error))]/30 bg-[hsl(var(--support-error))]/10 p-4 text-sm text-[hsl(var(--support-error))]">Unable to load payments.</div>

  const renderDetail = (p: ClientPayment) => (
    <DetailPanel
      title={p.payment_number}
      subtitle="Client payment received"
      documents={simulatedDocs(p.id)}
      timeline={simulatedTimeline()}
      metadata={[
        { label: 'Amount', value: formatCurrency(p.amount_cents, p.currency_code) },
        { label: 'Method', value: p.payment_method || 'Not specified' },
        { label: 'Reference', value: p.reference_number || 'None' },
        { label: 'Version', value: `v${p.version}` },
      ]}
    />
  )

  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Collections</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Client payments and billing allocations. Click a row to see payment proofs and allocation history.</p>
      </div>
      <div className="hidden sm:block">
        <DataTable columns={columns} data={items} caption="Client payments" emptyMessage="No payments recorded." zebra density="normal" renderDetail={renderDetail} />
      </div>
      <div className="flex flex-col gap-3 sm:hidden">
        {items.length === 0 ? <p className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-6 text-center text-sm text-[hsl(var(--text-03))]">No payments recorded.</p> : items.map((p) => (
          <details key={p.id} className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))]">
            <summary className="cursor-pointer p-4">
              <div className="flex items-center justify-between"><span className="font-medium text-[hsl(var(--text-01))]">{p.payment_number}</span><span className="text-[hsl(var(--text-03))]">v{p.version}</span></div>
              <p className="mt-2 text-sm text-[hsl(var(--text-01))]">{formatCurrency(p.amount_cents, p.currency_code)}</p>
            </summary>
            <div className="border-t border-[hsl(var(--border-subtle))] p-4">{renderDetail(p)}</div>
          </details>
        ))}
      </div>
    </div>
  )
}
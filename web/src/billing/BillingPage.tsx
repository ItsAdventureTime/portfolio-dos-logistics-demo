import { useQuery } from '@tanstack/react-query'
import { api, type BillingRecord } from '../api/client'
import { DataTable } from '../components/table/DataTable'
import { Badge, type StatusType } from '../components/feedback/Badge'
import { DetailPanel, type DocumentItem, type TimelineEvent } from '../components/feedback/DetailPanel'
import { type ColumnDef } from '@tanstack/react-table'

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

function simulatedDocs(billId: string): DocumentItem[] {
  const docs: Record<string, DocumentItem[]> = {
    'demo-bill-001': [
      { name: 'Invoice_ACME_100001.pdf', type: 'PDF', size: '95 KB' },
      { name: 'Delivery_Confirmation.pdf', type: 'PDF', size: '78 KB' },
    ],
    'demo-bill-002': [
      { name: 'Invoice_UMBRA_100002.pdf', type: 'PDF', size: '88 KB' },
      { name: 'Customs_Docs_Canada.pdf', type: 'PDF', size: '320 KB' },
    ],
    'demo-bill-003': [
      { name: 'Invoice_GLOBX_100003.pdf', type: 'PDF', size: '72 KB' },
    ],
    'demo-bill-004': [
      { name: 'Invoice_STARK_100004_Draft.pdf', type: 'PDF', size: '65 KB' },
    ],
  }
  return docs[billId] || []
}

function simulatedTimeline(status: string): TimelineEvent[] {
  return [
    { action: 'Billing record created', actor: 'Demo Administrator', timestamp: '2026-07-30 09:00' },
    ...(status !== 'draft' ? [{ action: 'Submitted for review', actor: 'Demo Administrator', timestamp: '2026-07-30 11:00' }] : []),
    ...(status === 'approved' || status === 'finalized' ? [{ action: 'Approved by finance', actor: 'Demo Administrator', timestamp: '2026-07-31 09:30' }] : []),
    ...(status === 'finalized' ? [{ action: 'Invoice finalized and sent', actor: 'Demo Administrator', timestamp: '2026-07-31 14:00' }] : []),
  ]
}

const columns: ColumnDef<BillingRecord, unknown>[] = [
  { header: 'Number', accessorKey: 'billing_number', cell: ({ row }) => <span className="font-medium text-[hsl(var(--text-01))]">{row.original.billing_number}</span> },
  { header: 'Status', accessorKey: 'status', cell: ({ row }) => { const c = statusMap[row.original.status] || statusMap.draft; return <Badge status={c.status}>{c.label}</Badge> } },
  { header: 'Total', accessorKey: 'total', cell: ({ row }) => <span className="text-right font-medium text-[hsl(var(--text-01))]">{formatCurrency(row.original.total, row.original.currency_code)}</span> },
  { header: 'Version', accessorKey: 'version', cell: ({ row }) => <span className="text-[hsl(var(--text-03))]">v{row.original.version}</span> },
]

export function BillingPage() {
  const { data, isLoading, isError } = useQuery<{ billing_records: BillingRecord[] }>({
    queryKey: ['billing'],
    queryFn: () => api.listBilling(),
  })

  const items = data?.billing_records || []

  if (isLoading) return <div className="h-64 animate-pulse rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-02))]" />
  if (isError) return <div className="rounded-[var(--radius-sm)] border border-[hsl(var(--support-error))]/30 bg-[hsl(var(--support-error))]/10 p-4 text-sm text-[hsl(var(--support-error))]">Unable to load billing records.</div>

  const renderDetail = (b: BillingRecord) => (
    <DetailPanel
      title={b.billing_number}
      subtitle={b.notes || 'No notes'}
      documents={simulatedDocs(b.id)}
      timeline={simulatedTimeline(b.status)}
      metadata={[
        { label: 'Status', value: statusMap[b.status]?.label || b.status },
        { label: 'Total', value: formatCurrency(b.total, b.currency_code) },
        { label: 'Currency', value: b.currency_code },
        { label: 'Version', value: `v${b.version}` },
      ]}
    />
  )

  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Billing</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Manage invoices and billing records. Click a row to see invoice documents and approval history.</p>
      </div>
      <div className="hidden sm:block">
        <DataTable columns={columns} data={items} caption="Billing records" emptyMessage="No billing records found." zebra density="normal" renderDetail={renderDetail} />
      </div>
      <div className="flex flex-col gap-3 sm:hidden">
        {items.length === 0 ? <p className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-6 text-center text-sm text-[hsl(var(--text-03))]">No billing records found.</p> : items.map((b) => { const c = statusMap[b.status] || statusMap.draft; return (
          <details key={b.id} className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))]">
            <summary className="cursor-pointer p-4">
              <div className="flex items-center justify-between"><span className="font-medium text-[hsl(var(--text-01))]">{b.billing_number}</span><Badge status={c.status}>{c.label}</Badge></div>
              <div className="mt-3 flex items-center justify-between text-sm"><span className="text-[hsl(var(--text-01))]">{formatCurrency(b.total, b.currency_code)}</span><span className="text-[hsl(var(--text-03))]">v{b.version}</span></div>
            </summary>
            <div className="border-t border-[hsl(var(--border-subtle))] p-4">{renderDetail(b)}</div>
          </details>
        ) })}
      </div>
    </div>
  )
}
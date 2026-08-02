import { useQuery } from '@tanstack/react-query'
import { api, type BudgetRequest } from '../api/client'
import { DataTable } from '../components/table/DataTable'
import { Badge, type StatusType } from '../components/feedback/Badge'
import { DetailPanel, type DocumentItem, type TimelineEvent } from '../components/feedback/DetailPanel'
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

function simulatedDocs(brId: string): DocumentItem[] {
  const docs: Record<string, DocumentItem[]> = {
    'demo-br-001': [
      { name: 'Funding_Justification_ACME.pdf', type: 'PDF', size: '125 KB' },
      { name: 'Quotation_Reference_Q-ACME-100001.pdf', type: 'PDF', size: '88 KB' },
    ],
    'demo-br-002': [
      { name: 'Oversize_Cargo_Quote.pdf', type: 'PDF', size: '340 KB' },
      { name: 'Insurance_Certificate.pdf', type: 'PDF', size: '95 KB' },
    ],
    'demo-br-003': [
      { name: 'Reefer_Fuel_Costs.pdf', type: 'PDF', size: '72 KB' },
    ],
  }
  return docs[brId] || []
}

function simulatedTimeline(status: string): TimelineEvent[] {
  return [
    { action: 'Budget request created', actor: 'Demo Administrator', timestamp: '2026-07-29 10:00' },
    ...(status !== 'draft' ? [{ action: 'Submitted for controls review', actor: 'Demo Administrator', timestamp: '2026-07-29 14:00' }] : []),
    ...(status === 'approved' ? [{ action: 'Approved by controls review', actor: 'Demo Administrator', timestamp: '2026-07-30 09:00' }] : []),
    ...(status === 'returned' ? [{ action: 'Returned for clarification', actor: 'Demo Administrator', timestamp: '2026-07-30 11:00' }] : []),
  ]
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

  const renderDetail = (br: BudgetRequest) => (
    <DetailPanel
      title={br.request_number}
      subtitle={br.purpose}
      documents={simulatedDocs(br.id)}
      timeline={simulatedTimeline(br.status)}
      metadata={[
        { label: 'Status', value: statusMap[br.status]?.label || br.status },
        { label: 'Amount', value: formatCurrency(br.amount_cents, br.currency_code) },
        { label: 'Purpose', value: br.purpose },
        { label: 'Version', value: `v${br.version}` },
      ]}
    />
  )

  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Funding requests</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Shipment funding linked to accepted quotations. Click a row to see justifications and approval history.</p>
      </div>
      <div className="hidden sm:block">
        <DataTable columns={columns} data={items} caption="Budget requests" emptyMessage="No funding requests found." zebra density="normal" renderDetail={renderDetail} />
      </div>
      <div className="flex flex-col gap-3 sm:hidden">
        {items.length === 0 ? <p className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-6 text-center text-sm text-[hsl(var(--text-03))]">No funding requests found.</p> : items.map((br) => { const c = statusMap[br.status] || statusMap.draft; return (
          <details key={br.id} className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))]">
            <summary className="cursor-pointer p-4">
              <div className="flex items-center justify-between"><span className="font-medium text-[hsl(var(--text-01))]">{br.request_number}</span><Badge status={c.status}>{c.label}</Badge></div>
              <p className="mt-2 text-sm text-[hsl(var(--text-03))]">{br.purpose}</p>
              <div className="mt-3 flex items-center justify-between text-sm"><span className="text-[hsl(var(--text-01))]">{formatCurrency(br.amount_cents, br.currency_code)}</span><span className="text-[hsl(var(--text-03))]">v{br.version}</span></div>
            </summary>
            <div className="border-t border-[hsl(var(--border-subtle))] p-4">{renderDetail(br)}</div>
          </details>
        ) })}
      </div>
    </div>
  )
}
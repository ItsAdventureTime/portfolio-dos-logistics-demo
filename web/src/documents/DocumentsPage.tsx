import { useState } from 'react'
import { DocumentList, type DocItem } from '../components/feedback/DocViewer'
import { Badge } from '../components/feedback/Badge'

// Simulated documents across different workflow entities
const allDocuments: (DocItem & { entity: string; entityType: string; uploadedBy: string; date: string })[] = [
  // Quotation documents
  { name: 'Bill_of_Lading_MNL_LA.pdf', type: 'PDF', size: '248 KB', entity: 'Q-ACME-100001', entityType: 'Quotation', uploadedBy: 'Demo Administrator', date: '2026-07-28' },
  { name: 'Container_Photo_40ft.jpg', type: 'IMG', size: '1.2 MB', entity: 'Q-ACME-100001', entityType: 'Quotation', uploadedBy: 'Juan Dela Cruz', date: '2026-07-28' },
  { name: 'Rate_Confirmation_ACME.pdf', type: 'PDF', size: '96 KB', entity: 'Q-ACME-100001', entityType: 'Quotation', uploadedBy: 'Demo Administrator', date: '2026-07-28' },
  { name: 'Airway_Bill_JFK_LHR.pdf', type: 'PDF', size: '182 KB', entity: 'Q-GLOBX-100002', entityType: 'Quotation', uploadedBy: 'Demo Administrator', date: '2026-07-29' },
  { name: 'Rail_Quote_Bulk_HOU_CHI.pdf', type: 'PDF', size: '145 KB', entity: 'Q-STARK-100003', entityType: 'Quotation', uploadedBy: 'Demo Administrator', date: '2026-07-29' },
  { name: 'Cargo_Manifest_Stark.pdf', type: 'PDF', size: '310 KB', entity: 'Q-STARK-100003', entityType: 'Quotation', uploadedBy: 'Juan Dela Cruz', date: '2026-07-29' },
  { name: 'Oversize_Permit_ON.pdf', type: 'PDF', size: '520 KB', entity: 'Q-UMBRA-100005', entityType: 'Quotation', uploadedBy: 'Demo Administrator', date: '2026-07-30' },

  // Disbursement documents
  { name: 'Wire_Confirmation_001.pdf', type: 'PDF', size: '42 KB', entity: 'DISB-001', entityType: 'Disbursement', uploadedBy: 'Demo Administrator', date: '2026-07-30' },
  { name: 'Bank_Receipt_BankA.pdf', type: 'PDF', size: '58 KB', entity: 'DISB-001', entityType: 'Disbursement', uploadedBy: 'Demo Administrator', date: '2026-07-30' },
  { name: 'Disbursement_Authorization.pdf', type: 'PDF', size: '120 KB', entity: 'DISB-001', entityType: 'Disbursement', uploadedBy: 'Demo Administrator', date: '2026-07-30' },
  { name: 'Wire_Confirmation_002.pdf', type: 'PDF', size: '38 KB', entity: 'DISB-002', entityType: 'Disbursement', uploadedBy: 'Demo Administrator', date: '2026-07-31' },
  { name: 'Payment_Proof_UMBRA.jpg', type: 'IMG', size: '850 KB', entity: 'DISB-002', entityType: 'Disbursement', uploadedBy: 'Juan Dela Cruz', date: '2026-07-31' },

  // Billing documents
  { name: 'Invoice_ACME_100001.pdf', type: 'PDF', size: '95 KB', entity: 'INV-ACME-100001', entityType: 'Billing', uploadedBy: 'Demo Administrator', date: '2026-07-31' },
  { name: 'Delivery_Confirmation.pdf', type: 'PDF', size: '78 KB', entity: 'INV-ACME-100001', entityType: 'Billing', uploadedBy: 'Juan Dela Cruz', date: '2026-07-31' },
  { name: 'Invoice_UMBRA_100002.pdf', type: 'PDF', size: '88 KB', entity: 'INV-UMBRA-100002', entityType: 'Billing', uploadedBy: 'Demo Administrator', date: '2026-07-31' },
  { name: 'Customs_Docs_Canada.pdf', type: 'PDF', size: '320 KB', entity: 'INV-UMBRA-100002', entityType: 'Billing', uploadedBy: 'Demo Administrator', date: '2026-07-31' },

  // Liquidation documents
  { name: 'Liquidation_Report.pdf', type: 'PDF', size: '210 KB', entity: 'LIQ-001', entityType: 'Liquidation', uploadedBy: 'Demo Administrator', date: '2026-07-31' },
  { name: 'Actual_Spending_Receipts.zip', type: 'DOC', size: '2.4 MB', entity: 'LIQ-001', entityType: 'Liquidation', uploadedBy: 'Juan Dela Cruz', date: '2026-07-31' },
  { name: 'Variance_Explanation.pdf', type: 'PDF', size: '45 KB', entity: 'LIQ-001', entityType: 'Liquidation', uploadedBy: 'Demo Administrator', date: '2026-07-31' },

  // Payment documents
  { name: 'Wire_Confirmation_ACME_001.pdf', type: 'PDF', size: '42 KB', entity: 'PMT-ACME-100001', entityType: 'Payment', uploadedBy: 'Demo Administrator', date: '2026-08-01' },
  { name: 'Bank_Statement_July.pdf', type: 'PDF', size: '180 KB', entity: 'PMT-ACME-100001', entityType: 'Payment', uploadedBy: 'Demo Administrator', date: '2026-08-01' },
  { name: 'Remittance_Advice_UMBRA.pdf', type: 'PDF', size: '62 KB', entity: 'PMT-UMBRA-100003', entityType: 'Payment', uploadedBy: 'Demo Administrator', date: '2026-08-01' },
]

export function DocumentsPage() {
  const [filter, setFilter] = useState<string>('all')

  const entityTypes = ['all', ...Array.from(new Set(allDocuments.map(d => d.entityType)))]
  const filtered = filter === 'all' ? allDocuments : allDocuments.filter(d => d.entityType === filter)

  // Group by entity
  const grouped = filtered.reduce((acc, doc) => {
    const key = doc.entity
    if (!acc[key]) acc[key] = []
    acc[key].push(doc)
    return acc
  }, {} as Record<string, typeof allDocuments>)

  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Documents</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Supporting documents across all workflow entities. Click any file to preview.</p>
      </div>

      {/* Filter tabs */}
      <div className="flex flex-wrap gap-2">
        {entityTypes.map((type) => (
          <button
            key={type}
            onClick={() => setFilter(type)}
            className={`rounded-full border px-3 py-1 text-xs font-medium transition-all var(--duration-fast-01) var(--ease-productive) ${
              filter === type
                ? 'border-[hsl(var(--interactive-01))] bg-[hsl(var(--interactive-01))]/12 text-[hsl(var(--interactive-04))]'
                : 'border-[hsl(var(--border-subtle))] text-[hsl(var(--text-03))] hover:border-[hsl(var(--border-interactive))] hover:text-[hsl(var(--text-02))]'
            }`}
          >
            {type === 'all' ? 'All documents' : type}
            <span className="ml-1.5 text-[10px] opacity-70">
              ({type === 'all' ? allDocuments.length : allDocuments.filter(d => d.entityType === type).length})
            </span>
          </button>
        ))}
      </div>

      {/* Summary stats */}
      <div className="flex gap-4 text-sm text-[hsl(var(--text-03))]">
        <span>{filtered.length} documents</span>
        <span>{Object.keys(grouped).length} entities</span>
      </div>

      {/* Document groups */}
      <div className="flex flex-col gap-4">
        {Object.entries(grouped).map(([entity, docs]) => (
          <div
            key={entity}
            className="rounded-[var(--radius-md)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-4 shadow-[var(--shadow-01)]"
          >
            <div className="mb-3 flex items-center justify-between border-b border-[hsl(var(--border-subtle))] pb-2">
              <div className="flex items-center gap-2">
                <span className="font-mono text-sm font-medium text-[hsl(var(--text-01))]">{entity}</span>
                <Badge status="neutral">{docs[0].entityType}</Badge>
              </div>
              <span className="text-xs text-[hsl(var(--text-03))]">{docs.length} file{docs.length > 1 ? 's' : ''}</span>
            </div>
            <DocumentList documents={docs} />
          </div>
        ))}
      </div>

      {filtered.length === 0 && (
        <div className="rounded-[var(--radius-md)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-8 text-center">
          <p className="text-sm text-[hsl(var(--text-03))]">No documents found for this filter.</p>
        </div>
      )}
    </div>
  )
}
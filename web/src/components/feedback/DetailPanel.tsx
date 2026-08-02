interface LineItem {
  id: string
  description: string
  quantity: number
  unit_price: number
  line_total: number
}

interface TimelineEvent {
  action: string
  actor: string
  timestamp: string
}

interface DocumentItem {
  name: string
  type: string
  size: string
}

interface DetailPanelProps {
  title: string
  subtitle?: string
  lines?: LineItem[]
  timeline?: TimelineEvent[]
  documents?: DocumentItem[]
  metadata?: { label: string; value: string }[]
}

function formatCurrency(cents: number, code = 'USD'): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: code }).format(cents / 100)
}

export function DetailPanel({ title, subtitle, lines, timeline, documents, metadata }: DetailPanelProps) {
  return (
    <div className="flex flex-col gap-5">
      {/* Header */}
      <div className="flex items-start justify-between border-b border-[hsl(var(--border-subtle))] pb-3">
        <div>
          <h3 className="text-sm font-semibold text-[hsl(var(--text-01))]">{title}</h3>
          {subtitle && <p className="mt-0.5 text-xs text-[hsl(var(--text-03))]">{subtitle}</p>}
        </div>
      </div>

      {/* Metadata grid */}
      {metadata && metadata.length > 0 && (
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
          {metadata.map((m, i) => (
            <div key={i}>
              <p className="text-[11px] font-medium uppercase tracking-wide text-[hsl(var(--text-03))]">{m.label}</p>
              <p className="mt-0.5 text-sm text-[hsl(var(--text-01))]">{m.value}</p>
            </div>
          ))}
        </div>
      )}

      {/* Line items */}
      {lines && lines.length > 0 && (
        <div>
          <h4 className="mb-2 text-xs font-semibold uppercase tracking-wide text-[hsl(var(--text-03))]">Line items</h4>
          <div className="overflow-x-auto rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))]">
            <table className="w-full text-xs">
              <thead className="bg-[hsl(var(--ui-03))]">
                <tr>
                  <th className="px-3 py-2 text-left font-semibold text-[hsl(var(--text-02))]">Description</th>
                  <th className="px-3 py-2 text-right font-semibold text-[hsl(var(--text-02))]">Qty</th>
                  <th className="px-3 py-2 text-right font-semibold text-[hsl(var(--text-02))]">Unit price</th>
                  <th className="px-3 py-2 text-right font-semibold text-[hsl(var(--text-02))]">Total</th>
                </tr>
              </thead>
              <tbody>
                {lines.map((line, i) => (
                  <tr key={line.id} className={i % 2 === 1 ? 'bg-[hsl(var(--ui-02))]/50' : ''}>
                    <td className="px-3 py-2 text-[hsl(var(--text-01))]">{line.description}</td>
                    <td className="px-3 py-2 text-right text-[hsl(var(--text-02))]">{line.quantity}</td>
                    <td className="px-3 py-2 text-right text-[hsl(var(--text-02))]">{formatCurrency(line.unit_price)}</td>
                    <td className="px-3 py-2 text-right font-medium text-[hsl(var(--text-01))]">{formatCurrency(line.line_total)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Documents */}
      {documents && documents.length > 0 && (
        <div>
          <h4 className="mb-2 text-xs font-semibold uppercase tracking-wide text-[hsl(var(--text-03))]">Documents</h4>
          <ul className="flex flex-col gap-2">
            {documents.map((doc, i) => (
              <li key={i} className="flex items-center justify-between rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-02))]/50 px-3 py-2">
                <div className="flex items-center gap-2.5">
                  <span className="flex h-7 w-7 items-center justify-center rounded-[var(--radius-sm)] bg-[hsl(var(--interactive-01))]/10 text-xs text-[hsl(var(--interactive-04))]" aria-hidden="true">
                    {doc.type === 'PDF' ? 'PDF' : doc.type === 'IMG' ? 'IMG' : 'DOC'}
                  </span>
                  <div>
                    <p className="text-sm text-[hsl(var(--text-01))]">{doc.name}</p>
                    <p className="text-[11px] text-[hsl(var(--text-03))]">{doc.type} - {doc.size}</p>
                  </div>
                </div>
                <button className="text-xs text-[hsl(var(--interactive-04))] hover:underline" aria-label={`Download ${doc.name}`}>
                  View
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Timeline */}
      {timeline && timeline.length > 0 && (
        <div>
          <h4 className="mb-2 text-xs font-semibold uppercase tracking-wide text-[hsl(var(--text-03))]">Activity timeline</h4>
          <ul className="flex flex-col gap-2.5">
            {timeline.map((event, i) => (
              <li key={i} className="flex gap-3">
                <div className="flex flex-col items-center">
                  <span className="flex h-2 w-2 rounded-full bg-[hsl(var(--interactive-01))]" aria-hidden="true" />
                  {i < timeline.length - 1 && <span className="w-px flex-1 bg-[hsl(var(--border-subtle))]" aria-hidden="true" />}
                </div>
                <div className="pb-3">
                  <p className="text-sm text-[hsl(var(--text-01))]">{event.action}</p>
                  <p className="text-[11px] text-[hsl(var(--text-03))]">{event.actor} - {event.timestamp}</p>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  )
}

export { type LineItem, type TimelineEvent, type DocumentItem }
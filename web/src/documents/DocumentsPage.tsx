export function DocumentsPage() {
  return (
    <div className="flex flex-col gap-5">
      <div>
        <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">Documents</h1>
        <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">Supporting documents for quotations, disbursements, and billing</p>
      </div>
      <div className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-8 text-center">
        <p className="text-sm text-[hsl(var(--text-03))]">No documents uploaded yet.</p>
        <p className="mt-1 text-xs text-[hsl(var(--text-03))]">Upload payment proofs, liquidation evidence, and supporting files from the relevant workflow page.</p>
      </div>
    </div>
  )
}
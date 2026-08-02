import { useState } from 'react'
import { Button } from '../ui/Button'

interface DocViewerModalProps {
  isOpen: boolean
  onClose: () => void
  fileName: string
  fileType: string
  fileSize: string
}

// Generates a simulated PDF-like content view in the browser.
// Creates a blob URL with a text representation that opens in a new tab,
// or shows an inline preview for images.
function generateSimulatedDocument(name: string, type: string): string {
  const content = `%PDF-1.4
DOS FreightFlow Control - Simulated Document
==============================================

Document: ${name}
Type: ${type}
Generated: ${new Date().toISOString()}

This is a simulated document for demonstration purposes.
In production, this would be served from private object storage
behind authenticated authorization checks.

----------------------------------------------
DOS FreightFlow Control
DelegateOps Business Support Services
----------------------------------------------`

  const blob = new Blob([content], { type: type === 'IMG' ? 'image/png' : 'application/pdf' })
  return URL.createObjectURL(blob)
}

export function DocViewerModal({ isOpen, onClose, fileName, fileType, fileSize }: DocViewerModalProps) {
  if (!isOpen) return null

  const handleView = () => {
    const url = generateSimulatedDocument(fileName, fileType)
    window.open(url, '_blank')
  }

  const handleDownload = () => {
    const url = generateSimulatedDocument(fileName, fileType)
    const a = document.createElement('a')
    a.href = url
    a.download = fileName
    a.click()
    URL.revokeObjectURL(url)
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
      onClick={onClose}
      style={{ animation: 'var(--duration-moderate-02) var(--ease-entrance) fadeIn' }}
    >
      <div
        className="w-full max-w-md rounded-[var(--radius-md)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-6 shadow-[var(--shadow-04)]"
        onClick={(e) => e.stopPropagation()}
        style={{ animation: 'var(--duration-moderate-02) var(--ease-entrance) scaleIn' }}
      >
        <div className="mb-4 flex items-start justify-between">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-[var(--radius-sm)] bg-[hsl(var(--interactive-01))]/12 text-sm font-semibold text-[hsl(var(--interactive-04))]">
              {fileType === 'PDF' ? 'PDF' : fileType === 'IMG' ? 'IMG' : 'DOC'}
            </div>
            <div>
              <h3 className="text-sm font-semibold text-[hsl(var(--text-01))]">{fileName}</h3>
              <p className="text-xs text-[hsl(var(--text-03))]">{fileType} - {fileSize}</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="flex h-8 w-8 items-center justify-center rounded-[var(--radius-sm)] text-[hsl(var(--text-03))] transition-colors var(--duration-fast-01) var(--ease-productive) hover:bg-[hsl(var(--hover-ui))] hover:text-[hsl(var(--text-01))]"
            aria-label="Close"
          >
            <span aria-hidden="true">{'\u2715'}</span>
          </button>
        </div>

        <div className="mb-5 rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-02))]/50 p-4">
          <p className="text-xs text-[hsl(var(--text-03))]">
            This document is stored in private object storage and requires authentication to access.
            In the demo environment, a simulated preview is available.
          </p>
        </div>

        <div className="flex gap-3">
          <Button onPress={handleView} size="sm" className="flex-1">
            Open preview
          </Button>
          <Button onPress={handleDownload} variant="secondary" size="sm" className="flex-1">
            Download
          </Button>
        </div>
      </div>

      <style>{`
        @keyframes fadeIn {
          from { opacity: 0; }
          to { opacity: 1; }
        }
        @keyframes scaleIn {
          from { opacity: 0; transform: scale(0.95); }
          to { opacity: 1; transform: scale(1); }
        }
      `}</style>
    </div>
  )
}

// Reusable document list item with modal viewer
interface DocItem {
  name: string
  type: string
  size: string
}

export function DocumentList({ documents }: { documents: DocItem[] }) {
  const [selectedDoc, setSelectedDoc] = useState<DocItem | null>(null)

  if (!documents || documents.length === 0) return null

  return (
    <>
      <ul className="flex flex-col gap-2">
        {documents.map((doc, i) => (
          <li key={i}>
            <button
              onClick={() => setSelectedDoc(doc)}
              className="flex w-full items-center justify-between rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-02))]/50 px-3 py-2.5 text-left transition-all var(--duration-fast-01) var(--ease-productive) hover:border-[hsl(var(--border-interactive))] hover:bg-[hsl(var(--hover-ui))] focus:outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[hsl(var(--focus))]"
            >
              <div className="flex items-center gap-2.5">
                <span className="flex h-7 w-7 items-center justify-center rounded-[var(--radius-sm)] bg-[hsl(var(--interactive-01))]/12 text-[10px] font-semibold text-[hsl(var(--interactive-04))]" aria-hidden="true">
                  {doc.type === 'PDF' ? 'PDF' : doc.type === 'IMG' ? 'IMG' : 'DOC'}
                </span>
                <div>
                  <p className="text-sm text-[hsl(var(--text-01))]">{doc.name}</p>
                  <p className="text-[11px] text-[hsl(var(--text-03))]">{doc.type} - {doc.size}</p>
                </div>
              </div>
              <span className="text-xs text-[hsl(var(--interactive-04))]">View</span>
            </button>
          </li>
        ))}
      </ul>

      <DocViewerModal
        isOpen={!!selectedDoc}
        onClose={() => setSelectedDoc(null)}
        fileName={selectedDoc?.name || ''}
        fileType={selectedDoc?.type || ''}
        fileSize={selectedDoc?.size || ''}
      />
    </>
  )
}

export { type DocItem, type DocViewerModalProps }
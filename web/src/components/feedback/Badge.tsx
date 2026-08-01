import { type ReactNode } from 'react'

type StatusType = 'success' | 'warning' | 'error' | 'info' | 'neutral'

interface BadgeProps {
  children: ReactNode
  status?: StatusType
  className?: string
}

const statusConfig: Record<StatusType, { bg: string; text: string; label: string }> = {
  success: { bg: 'bg-[hsl(var(--status-success))]/15', text: 'text-[hsl(var(--status-success))]', label: 'Success' },
  warning: { bg: 'bg-[hsl(var(--status-warning))]/15', text: 'text-[hsl(var(--status-warning))]', label: 'Warning' },
  error: { bg: 'bg-[hsl(var(--status-error))]/15', text: 'text-[hsl(var(--status-error))]', label: 'Error' },
  info: { bg: 'bg-[hsl(var(--status-info))]/15', text: 'text-[hsl(var(--status-info))]', label: 'Info' },
  neutral: { bg: 'bg-[hsl(var(--surface-panel))]', text: 'text-[hsl(var(--content-muted))]', label: 'Neutral' },
}

export function Badge({ children, status = 'neutral', className = '' }: BadgeProps) {
  const config = statusConfig[status]
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-md px-2 py-0.5 text-xs font-medium ${config.bg} ${config.text} ${className}`}
    >
      <span className="inline-block h-1.5 w-1.5 rounded-full bg-current" aria-hidden="true" />
      {children}
      <span className="sr-only">{config.label}</span>
    </span>
  )
}

export function StatusPill({ status, label }: { status: StatusType; label: string }) {
  return (
    <Badge status={status}>{label}</Badge>
  )
}

export type { StatusType }
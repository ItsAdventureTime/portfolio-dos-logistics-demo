import { type ReactNode } from 'react'

type StatusType = 'success' | 'warning' | 'error' | 'info' | 'neutral'

interface BadgeProps {
  children: ReactNode
  status?: StatusType
  className?: string
}

const statusConfig: Record<StatusType, { bg: string; text: string; border: string; label: string }> = {
  success: {
    bg: 'bg-[hsl(var(--support-success))]/12',
    text: 'text-[hsl(var(--support-success))]',
    border: 'border-[hsl(var(--support-success))]/30',
    label: 'Success',
  },
  warning: {
    bg: 'bg-[hsl(var(--support-warning))]/12',
    text: 'text-[hsl(var(--support-warning))]',
    border: 'border-[hsl(var(--support-warning))]/30',
    label: 'Warning',
  },
  error: {
    bg: 'bg-[hsl(var(--support-error))]/12',
    text: 'text-[hsl(var(--support-error))]',
    border: 'border-[hsl(var(--support-error))]/30',
    label: 'Error',
  },
  info: {
    bg: 'bg-[hsl(var(--support-info))]/12',
    text: 'text-[hsl(var(--support-info))]',
    border: 'border-[hsl(var(--support-info))]/30',
    label: 'Info',
  },
  neutral: {
    bg: 'bg-[hsl(var(--ui-03))]',
    text: 'text-[hsl(var(--text-02))]',
    border: 'border-[hsl(var(--border-subtle))]',
    label: 'Neutral',
  },
}

export function Badge({ children, status = 'neutral', className = '' }: BadgeProps) {
  const config = statusConfig[status]
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-[var(--radius-sm)] border px-2 py-0.5 text-xs font-medium ${config.bg} ${config.text} ${config.border} ${className}`}
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
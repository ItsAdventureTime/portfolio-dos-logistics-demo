import { type ReactNode } from 'react'

type StatusType = 'success' | 'warning' | 'error' | 'info' | 'neutral'

interface BadgeProps {
  children: ReactNode
  status?: StatusType
  className?: string
}

const statusConfig: Record<StatusType, { bg: string; text: string; border: string; dot: string; label: string }> = {
  success: {
    bg: 'bg-[hsl(var(--support-success))]/12',
    text: 'text-[hsl(var(--support-success))]',
    border: 'border-[hsl(var(--support-success))]/25',
    dot: 'bg-[hsl(var(--support-success))]',
    label: 'Success',
  },
  warning: {
    bg: 'bg-[hsl(var(--support-warning))]/12',
    text: 'text-[hsl(var(--support-warning))]',
    border: 'border-[hsl(var(--support-warning))]/25',
    dot: 'bg-[hsl(var(--support-warning))]',
    label: 'Warning',
  },
  error: {
    bg: 'bg-[hsl(var(--support-error))]/12',
    text: 'text-[hsl(var(--support-error))]',
    border: 'border-[hsl(var(--support-error))]/25',
    dot: 'bg-[hsl(var(--support-error))]',
    label: 'Error',
  },
  info: {
    bg: 'bg-[hsl(var(--support-info))]/12',
    text: 'text-[hsl(var(--support-info))]',
    border: 'border-[hsl(var(--support-info))]/25',
    dot: 'bg-[hsl(var(--support-info))]',
    label: 'Info',
  },
  neutral: {
    bg: 'bg-[hsl(var(--ui-03))]',
    text: 'text-[hsl(var(--text-02))]',
    border: 'border-[hsl(var(--border-subtle))]',
    dot: 'bg-[hsl(var(--text-03))]',
    label: 'Neutral',
  },
}

export function Badge({ children, status = 'neutral', className = '' }: BadgeProps) {
  const config = statusConfig[status]
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium ${config.bg} ${config.text} ${config.border} ${className}`}
    >
      <span className={`inline-block h-1.5 w-1.5 rounded-full ${config.dot}`} aria-hidden="true" />
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
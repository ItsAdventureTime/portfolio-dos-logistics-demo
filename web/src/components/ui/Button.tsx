import { Button as RACButton, type ButtonProps } from 'react-aria-components'
import { type ReactNode } from 'react'

type Variant = 'primary' | 'secondary' | 'destructive' | 'ghost'
type Size = 'sm' | 'md' | 'lg'

interface DOSButtonProps extends ButtonProps {
  variant?: Variant
  size?: Size
  children: ReactNode
}

const variantClasses: Record<Variant, string> = {
  primary: 'bg-[hsl(var(--action))] text-[hsl(var(--action-foreground))] hover:opacity-90',
  secondary:
    'bg-[hsl(var(--surface-elevated))] text-[hsl(var(--content))] border border-[hsl(var(--border))] hover:border-[hsl(var(--border-interactive))]',
  destructive: 'bg-[hsl(var(--status-error))] text-white hover:opacity-90',
  ghost: 'bg-transparent text-[hsl(var(--content))] hover:bg-[hsl(var(--surface-elevated))]',
}

const sizeClasses: Record<Size, string> = {
  sm: 'px-3 py-1.5 text-sm',
  md: 'px-4 py-2 text-base',
  lg: 'px-6 py-3 text-lg',
}

export function Button({ variant = 'primary', size = 'md', className = '', children, ...props }: DOSButtonProps) {
  return (
    <RACButton
      className={`inline-flex items-center justify-center gap-2 rounded-lg font-medium transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-[hsl(var(--focus))] focus-visible:ring-offset-2 focus-visible:ring-offset-[hsl(var(--surface))] disabled:opacity-50 disabled:cursor-not-allowed ${variantClasses[variant]} ${sizeClasses[size]} ${className}`}
      {...props}
    >
      {children}
    </RACButton>
  )
}
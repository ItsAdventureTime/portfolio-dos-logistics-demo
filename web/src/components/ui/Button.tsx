import { Button as RACButton, type ButtonProps } from 'react-aria-components'
import { type ReactNode } from 'react'

type Variant = 'primary' | 'secondary' | 'tertiary' | 'destructive' | 'ghost'
type Size = 'sm' | 'md' | 'lg'

interface DOSButtonProps extends ButtonProps {
  variant?: Variant
  size?: Size
  children: ReactNode
}

const variantClasses: Record<Variant, string> = {
  primary: 'bg-[hsl(var(--interactive-01))] text-[hsl(var(--text-on-color))] hover:bg-[hsl(var(--hover-primary))] border border-transparent shadow-[var(--shadow-01)]',
  secondary: 'bg-[hsl(var(--ui-02))] text-[hsl(var(--text-01))] border border-[hsl(var(--border-strong))] hover:bg-[hsl(var(--hover-ui))] hover:border-[hsl(var(--border-interactive))]',
  tertiary: 'bg-transparent text-[hsl(var(--interactive-04))] border border-transparent hover:bg-[hsl(var(--hover-ui))]',
  destructive: 'bg-[hsl(var(--support-error))] text-white border border-transparent hover:opacity-90 shadow-[var(--shadow-01)]',
  ghost: 'bg-transparent text-[hsl(var(--text-02))] border border-transparent hover:bg-[hsl(var(--hover-ui))] hover:text-[hsl(var(--text-01))]',
}

const sizeClasses: Record<Size, string> = {
  sm: 'px-3 py-1.5 text-xs font-medium',
  md: 'px-4 py-2.5 text-sm font-medium',
  lg: 'px-5 py-3 text-sm font-medium',
}

export function Button({ variant = 'primary', size = 'md', className = '', children, ...props }: DOSButtonProps) {
  return (
    <RACButton
      className={`inline-flex items-center justify-center gap-2 rounded-[var(--radius-sm)] font-medium transition-all var(--duration-fast-01) var(--ease-productive) active:scale-[0.98] focus:outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[hsl(var(--focus))] disabled:opacity-40 disabled:cursor-not-allowed disabled:active:scale-100 ${variantClasses[variant]} ${sizeClasses[size]} ${className}`}
      {...props}
    >
      {children}
    </RACButton>
  )
}
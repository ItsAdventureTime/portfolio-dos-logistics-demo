import { TextField, Input as RACInput, Label, Text, type InputProps } from 'react-aria-components'
import { type ReactNode } from 'react'

interface DOSFieldProps {
  label: string
  name?: string
  type?: string
  value?: string
  defaultValue?: string
  placeholder?: string
  isRequired?: boolean
  isInvalid?: boolean
  errorMessage?: string
  description?: string
  onChange?: (value: string) => void
  autoComplete?: string
  className?: string
}

export function Field({
  label,
  name,
  type = 'text',
  value,
  defaultValue,
  placeholder,
  isRequired,
  isInvalid,
  errorMessage,
  description,
  onChange,
  autoComplete,
  className = '',
}: DOSFieldProps) {
  return (
    <TextField
      name={name}
      type={type}
      value={value}
      defaultValue={defaultValue}
      isRequired={isRequired}
      isInvalid={isInvalid}
      onChange={onChange}
      className={`flex flex-col gap-1.5 ${className}`}
    >
      <Label className="text-sm font-medium text-[hsl(var(--content))]">
        {label}
        {isRequired && <span className="ml-0.5 text-[hsl(var(--status-error))]" aria-hidden="true">*</span>}
      </Label>
      <RACInput
        placeholder={placeholder}
        autoComplete={autoComplete}
        className="w-full rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--surface-elevated))] px-3 py-2 text-[hsl(var(--content))] placeholder:text-[hsl(var(--content-muted))] focus:outline-none focus:border-[hsl(var(--focus))] focus:ring-1 focus:ring-[hsl(var(--focus))] data-[invalid]:border-[hsl(var(--status-error))]"
      />
      {description && !isInvalid && (
        <Text slot="description" className="text-xs text-[hsl(var(--content-muted))]">
          {description}
        </Text>
      )}
      {errorMessage && (
        <Text slot="errorMessage" className="text-xs text-[hsl(var(--status-error))]">
          {errorMessage}
        </Text>
      )}
    </TextField>
  )
}

interface DOSInputProps extends InputProps {
  className?: string
}

export function Input({ className = '', ...props }: DOSInputProps) {
  return (
    <RACInput
      className={`w-full rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--surface-elevated))] px-3 py-2 text-[hsl(var(--content))] placeholder:text-[hsl(var(--content-muted))] focus:outline-none focus:border-[hsl(var(--focus))] focus:ring-1 focus:ring-[hsl(var(--focus))] data-[invalid]:border-[hsl(var(--status-error))] ${className}`}
      {...props}
    />
  )
}

export { Label, TextField }
export type { DOSFieldProps }
export type { ReactNode }
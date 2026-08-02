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
      className={`flex flex-col gap-2 ${className}`}
    >
      <Label className="text-sm font-medium text-[hsl(var(--text-01))]">
        {label}
        {isRequired && <span className="ml-0.5 text-[hsl(var(--support-error))]" aria-hidden="true">*</span>}
      </Label>
      <RACInput
        placeholder={placeholder}
        autoComplete={autoComplete}
        className="w-full rounded-[var(--radius-sm)] border border-[hsl(var(--border-strong))] bg-[hsl(var(--ui-01))] px-3.5 py-2.5 text-sm text-[hsl(var(--text-01))] placeholder:text-[hsl(var(--text-03))] focus:outline-none focus:border-[hsl(var(--focus))] focus:ring-1 focus:ring-[hsl(var(--focus))] data-[invalid]:border-[hsl(var(--support-error))]"
      />
      {description && !isInvalid && (
        <Text slot="description" className="text-xs text-[hsl(var(--text-03))]">
          {description}
        </Text>
      )}
      {errorMessage && (
        <Text slot="errorMessage" className="text-xs text-[hsl(var(--support-error))]">
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
      className={`w-full rounded-[var(--radius-sm)] border border-[hsl(var(--border-strong))] bg-[hsl(var(--ui-01))] px-3.5 py-2.5 text-sm text-[hsl(var(--text-01))] placeholder:text-[hsl(var(--text-03))] focus:outline-none focus:border-[hsl(var(--focus))] focus:ring-1 focus:ring-[hsl(var(--focus))] data-[invalid]:border-[hsl(var(--support-error))] ${className}`}
      {...props}
    />
  )
}

export { Label, TextField }
export type { DOSFieldProps }
export type { ReactNode }
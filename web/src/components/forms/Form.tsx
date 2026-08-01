import { useForm, type FieldValues, type SubmitHandler, type UseFormRegister } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { type ZodType } from 'zod'
import { TextField, Label, Input } from '../ui/Field'
import { Button } from '../ui/Button'
import { type ReactNode } from 'react'

interface FormProps<T extends FieldValues> {
  schema: ZodType
  onSubmit: SubmitHandler<T>
  children?: ReactNode
  submitLabel?: string
}

export function Form<T extends FieldValues>({ schema, onSubmit, children, submitLabel = 'Submit' }: FormProps<T>) {
  const {
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<T>({ resolver: zodResolver(schema as never) as never })

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4" noValidate>
      {children}
      {Object.entries(errors).length > 0 && (
        <div role="alert" className="rounded-lg border border-[hsl(var(--status-error))] bg-[hsl(var(--status-error))]/10 p-3 text-sm text-[hsl(var(--status-error))]">
          Please correct the errors below and try again.
        </div>
      )}
      <Button type="submit" isDisabled={isSubmitting}>
        {isSubmitting ? 'Working…' : submitLabel}
      </Button>
    </form>
  )
}

interface FormFieldProps {
  label: string
  name: string
  type?: string
  placeholder?: string
  isRequired?: boolean
  description?: string
  autoComplete?: string
  defaultValue?: string
  inputMode?: 'numeric' | 'text' | 'email' | 'tel' | 'url'
  error?: string
  register: UseFormRegister<FieldValues>
}

export function FormField({
  label,
  name,
  type = 'text',
  placeholder,
  isRequired,
  description,
  autoComplete,
  defaultValue,
  inputMode,
  error,
  register,
}: FormFieldProps) {
  return (
    <TextField name={name} isInvalid={!!error} className="flex flex-col gap-1.5">
      <Label className="text-sm font-medium text-[hsl(var(--content))]">
        {label}
        {isRequired && <span className="ml-0.5 text-[hsl(var(--status-error))]" aria-hidden="true">*</span>}
      </Label>
      <Input
        type={type}
        placeholder={placeholder}
        autoComplete={autoComplete}
        inputMode={inputMode}
        defaultValue={defaultValue}
        {...register(name)}
      />
      {description && !error && (
        <span className="text-xs text-[hsl(var(--content-muted))]">{description}</span>
      )}
      {error && (
        <span role="alert" className="text-xs text-[hsl(var(--status-error))]">
          {error}
        </span>
      )}
    </TextField>
  )
}
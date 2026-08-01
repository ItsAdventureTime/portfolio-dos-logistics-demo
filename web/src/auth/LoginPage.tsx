import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { api } from '../api/client'
import { Button } from '../components/ui/Button'
import { TextField, Label, Input } from '../components/ui/Field'

const loginSchema = z.object({
  identifier: z.string().min(1, 'Enter your username or email.'),
  password: z.string().min(1, 'Enter your password.'),
})

const otpSchema = z.object({
  identifier: z.string().min(1, 'Enter your username or email.'),
  code: z.string().length(6, 'Enter the 6-digit verification code.'),
})

type LoginValues = z.infer<typeof loginSchema>
type OTPValues = z.infer<typeof otpSchema>

interface FieldProps {
  label: string
  name: string
  type?: string
  placeholder?: string
  isRequired?: boolean
  description?: string
  autoComplete?: string
  inputMode?: 'numeric' | 'text' | 'email'
  defaultValue?: string
  error?: string
}

function FieldInput({ label, name, type = 'text', placeholder, isRequired, description, autoComplete, inputMode, defaultValue, error }: FieldProps) {
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
      />
      {description && !error && <span className="text-xs text-[hsl(var(--content-muted))]">{description}</span>}
      {error && <span role="alert" className="text-xs text-[hsl(var(--status-error))]">{error}</span>}
    </TextField>
  )
}

export function LoginPage({ onSuccess }: { onSuccess: () => void }) {
  const [mode, setMode] = useState<'login' | 'otp'>('login')
  const [identifier, setIdentifier] = useState('')
  const [error, setError] = useState('')
  const [displayName, setDisplayName] = useState('')

  const loginForm = useForm<LoginValues>({ resolver: zodResolver(loginSchema) })
  const otpForm = useForm<OTPValues>({ resolver: zodResolver(otpSchema) })

  const handleLogin = async (values: LoginValues) => {
    setError('')
    setIdentifier(values.identifier)
    try {
      const res = await api.login(values.identifier, values.password)
      if (res.status === 'otp_required') {
        setDisplayName(res.display_name || '')
        setMode('otp')
      } else {
        onSuccess()
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Sign in failed. Please try again.')
    }
  }

  const handleVerify = async (values: OTPValues) => {
    setError('')
    try {
      await api.verifyEmail(values.identifier, values.code)
      onSuccess()
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Verification failed. Please try again.')
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-[hsl(var(--surface))] px-4">
      <div className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-semibold text-[hsl(var(--content))]">DOS FreightFlow Control</h1>
          <p className="mt-1 text-sm text-[hsl(var(--content-muted))]">Sign in to your workspace</p>
        </div>

        {mode === 'login' ? (
          <div className="rounded-xl border border-[hsl(var(--border))] bg-[hsl(var(--surface-elevated))] p-6">
            <form onSubmit={loginForm.handleSubmit(handleLogin)} className="flex flex-col gap-4" noValidate>
              <FieldInput
                label="Username or email"
                name="identifier"
                isRequired
                autoComplete="username"
                placeholder="admin or admin@company.example"
              />
              <FieldInput
                label="Password"
                name="password"
                type="password"
                isRequired
                autoComplete="current-password"
                error={error}
              />
              {Object.keys(loginForm.formState.errors).length > 0 && !error && (
                <div role="alert" className="rounded-lg bg-[hsl(var(--status-error))]/10 p-3 text-sm text-[hsl(var(--status-error))]">
                  Please correct the errors and try again.
                </div>
              )}
              <Button type="submit" isDisabled={loginForm.formState.isSubmitting}>
                {loginForm.formState.isSubmitting ? 'Working…' : 'Sign in'}
              </Button>
            </form>
          </div>
        ) : (
          <div className="rounded-xl border border-[hsl(var(--border))] bg-[hsl(var(--surface-elevated))] p-6">
            <h2 className="mb-2 text-lg font-medium text-[hsl(var(--content))]">Verify your email</h2>
            <p className="mb-4 text-sm text-[hsl(var(--content-muted))]">
              A verification code was sent to the email on file{displayName ? `, ${displayName}` : ''}. Enter it below.
            </p>
            <form onSubmit={otpForm.handleSubmit(handleVerify)} className="flex flex-col gap-4" noValidate>
              <FieldInput
                label="Username or email"
                name="identifier"
                isRequired
                autoComplete="username"
                defaultValue={identifier}
              />
              <FieldInput
                label="Verification code"
                name="code"
                isRequired
                placeholder="000000"
                description="6-digit code"
                inputMode="numeric"
                error={error}
              />
              <Button type="submit" isDisabled={otpForm.formState.isSubmitting}>
                {otpForm.formState.isSubmitting ? 'Working…' : 'Verify and sign in'}
              </Button>
            </form>
            <Button variant="ghost" size="sm" className="mt-4" onPress={() => setMode('login')}>
              Back to sign in
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}
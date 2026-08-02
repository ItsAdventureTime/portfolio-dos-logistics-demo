import { useState } from 'react'
import { useForm, type UseFormRegister } from 'react-hook-form'
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
  register: UseFormRegister<Record<string, unknown>>
}

function FieldInput({ label, name, type = 'text', placeholder, isRequired, description, autoComplete, inputMode, defaultValue, error, register }: FieldProps) {
  return (
    <TextField name={name} isInvalid={!!error} className="flex flex-col gap-2">
      <Label className="text-sm font-medium text-[hsl(var(--text-01))]">
        {label}
        {isRequired && <span className="ml-0.5 text-[hsl(var(--support-error))]" aria-hidden="true">*</span>}
      </Label>
      <Input
        type={type}
        placeholder={placeholder}
        autoComplete={autoComplete}
        inputMode={inputMode}
        defaultValue={defaultValue}
        {...register(name)}
      />
      {description && !error && <span className="text-xs text-[hsl(var(--text-03))]">{description}</span>}
      {error && <span role="alert" className="text-xs text-[hsl(var(--support-error))]">{error}</span>}
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
    <div className="flex min-h-screen items-center justify-center bg-[hsl(var(--ui-background))] px-4">
      <div className="w-full max-w-sm">
        {/* DOS brand lockup */}
        <div className="mb-8 flex flex-col items-center gap-4">
          <img
            src="/dos-mark.png"
            alt="DOS"
            className="h-12 w-12 rounded-[var(--radius-sm)]"
          />
          <div className="text-center">
            <h1 className="text-lg font-semibold text-[hsl(var(--text-01))]">DOS FreightFlow Control</h1>
            <p className="mt-0.5 text-sm text-[hsl(var(--text-03))]">DelegateOps Business Support Services</p>
          </div>
        </div>

        {mode === 'login' ? (
          <div className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-6">
            <p className="mb-5 text-sm text-[hsl(var(--text-02))]">Sign in to your workspace</p>
            <form onSubmit={loginForm.handleSubmit(handleLogin)} className="flex flex-col gap-4" noValidate>
              <FieldInput
                label="Username or email"
                name="identifier"
                isRequired
                autoComplete="username"
                placeholder="admin or admin@company.example"
                register={loginForm.register as UseFormRegister<Record<string, unknown>>}
              />
              <FieldInput
                label="Password"
                name="password"
                type="password"
                isRequired
                autoComplete="current-password"
                error={error}
                register={loginForm.register as UseFormRegister<Record<string, unknown>>}
              />
              {Object.keys(loginForm.formState.errors).length > 0 && !error && (
                <div role="alert" className="rounded-[var(--radius-sm)] border border-[hsl(var(--support-error))]/30 bg-[hsl(var(--support-error))]/10 p-3 text-xs text-[hsl(var(--support-error))]">
                  Please correct the errors and try again.
                </div>
              )}
              <Button type="submit" isDisabled={loginForm.formState.isSubmitting}>
                {loginForm.formState.isSubmitting ? 'Working...' : 'Sign in'}
              </Button>
            </form>
          </div>
        ) : (
          <div className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-6">
            <h2 className="mb-1.5 text-sm font-semibold text-[hsl(var(--text-01))]">Verify your email</h2>
            <p className="mb-5 text-sm text-[hsl(var(--text-03))]">
              A verification code was sent to the email on file{displayName ? `, ${displayName}` : ''}. Enter it below.
            </p>
            <form onSubmit={otpForm.handleSubmit(handleVerify)} className="flex flex-col gap-4" noValidate>
              <FieldInput
                label="Username or email"
                name="identifier"
                isRequired
                autoComplete="username"
                defaultValue={identifier}
                register={otpForm.register as UseFormRegister<Record<string, unknown>>}
              />
              <FieldInput
                label="Verification code"
                name="code"
                isRequired
                placeholder="000000"
                description="6-digit code"
                inputMode="numeric"
                error={error}
                register={otpForm.register as UseFormRegister<Record<string, unknown>>}
              />
              <Button type="submit" isDisabled={otpForm.formState.isSubmitting}>
                {otpForm.formState.isSubmitting ? 'Working...' : 'Verify and sign in'}
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
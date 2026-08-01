import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import App from './App'

describe('App', () => {
  it('renders a loading state while checking auth', () => {
    render(<App />)
    // While auth is being checked, we show a loading message
    expect(screen.getByText(/loading/i)).toBeInTheDocument()
  })
})
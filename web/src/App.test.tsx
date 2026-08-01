import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import App from './App'

describe('App', () => {
  it('renders the product title', () => {
    render(<App />)
    expect(screen.getByText(/DOS FreightFlow Control/i)).toBeInTheDocument()
  })

  it('renders a primary action button', () => {
    render(<App />)
    expect(screen.getByRole('button', { name: /primary action/i })).toBeInTheDocument()
  })
})
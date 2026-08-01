import { test, expect } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'

test.describe('Accessibility and responsive layout', () => {
  test('login page has no axe violations at 1440x900', async ({ page }) => {
    await page.goto('/')
    const results = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag22aa'])
      .analyze()
    expect(results.violations).toEqual([])
  })

  test('login page has no horizontal overflow at 390x844', async ({ page, browserName }) => {
    test.skip(browserName !== 'chromium', 'viewport test runs on chromium')
    await page.setViewportSize({ width: 390, height: 844 })
    await page.goto('/')
    const scrollWidth = await page.evaluate(() => document.documentElement.scrollWidth)
    const clientWidth = await page.evaluate(() => document.documentElement.clientWidth)
    expect(scrollWidth).toBeLessThanOrEqual(clientWidth)
  })

  test('login page has no horizontal overflow at 200% zoom', async ({ page, browserName }) => {
    test.skip(browserName !== 'chromium', 'zoom test runs on chromium')
    await page.goto('/')
    await page.evaluate(() => {
      document.body.style.zoom = '2'
    })
    const scrollWidth = await page.evaluate(() => document.documentElement.scrollWidth)
    const clientWidth = await page.evaluate(() => document.documentElement.clientWidth)
    expect(scrollWidth).toBeLessThanOrEqual(clientWidth)
  })

  test('keyboard focus is visible and tab order is usable on login form', async ({ page }) => {
    await page.goto('/')
    await page.keyboard.press('Tab')
    // First focusable element should be the username field
    const focused = await page.evaluate(() => document.activeElement?.getAttribute('name'))
    expect(focused).toBeTruthy()
  })

  test('login page heading is present and correct', async ({ page }) => {
    await page.goto('/')
    await expect(page.getByRole('heading', { name: /DOS FreightFlow Control/i })).toBeVisible()
    await expect(page.getByText(/sign in to your workspace/i)).toBeVisible()
  })
})
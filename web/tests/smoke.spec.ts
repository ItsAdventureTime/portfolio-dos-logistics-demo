import { test, expect } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'

test.describe('Stage 2 smoke', () => {
  test('home page loads with no a11y violations at 1440x900', async ({ page }) => {
    await page.goto('/')
    await expect(page.getByRole('heading', { name: /DOS FreightFlow Control/i })).toBeVisible()

    const results = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag22aa'])
      .analyze()

    expect(results.violations).toEqual([])
  })

  test('no horizontal overflow at mobile 390x844', async ({ page, browserName }) => {
    test.skip(browserName !== 'chromium', 'viewport test runs on chromium')
    await page.goto('/')
    const scrollWidth = await page.evaluate(() => document.documentElement.scrollWidth)
    const clientWidth = await page.evaluate(() => document.documentElement.clientWidth)
    expect(scrollWidth).toBeLessThanOrEqual(clientWidth)
  })
})
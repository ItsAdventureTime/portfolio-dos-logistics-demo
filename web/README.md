# DOS FreightFlow Control — Web Client

The frontend for DOS FreightFlow Control, built with React 19, TypeScript,
Vite, and Tailwind CSS v4. Component primitives come from React Aria Components
(imported as packages, not vendored). Operational tables use TanStack Table v8.
Server state runs through TanStack Query. Forms use React Hook Form with Zod
validation.

## Scripts

```sh
npm run dev        # start the dev server
npm run build      # typecheck and build for production
npm run lint        # run Oxlint
npm run typecheck   # TypeScript strict mode check
npm run test        # run Vitest unit tests
npm run test:e2e    # run Playwright + axe-core e2e tests
```

## Testing

Unit tests use Vitest with Testing Library. End-to-end tests run through
Playwright at 1440x900 and 390x844 with axe-core WCAG 2.2 AA checks.
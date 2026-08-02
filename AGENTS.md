# DOS FreightFlow Control: Implementation Guide

This repository holds the implementation of DOS FreightFlow Control, a private
B2B workspace for freight movement and the related financial controls. The
build follows the specification package in [`docs/spec/`](docs/spec/), which
contains the product requirements, domain model, security boundaries,
architecture options, acceptance matrix, implementation plan, and approved
DOS design assets.

## Getting started

Before writing any code:

1. Read [`docs/spec/AGENTS.md`](docs/spec/AGENTS.md) for the build agent
   instructions.
2. Read [`docs/spec/00-build-contract.md`](docs/spec/00-build-contract.md) for
   the working rules and legal boundaries.
3. Read [`DECISIONS.md`](DECISIONS.md) for the architecture decision log. Every
   choice that diverges from the spec's suggestions is recorded there with
   tradeoffs, evidence, and a rollback plan.
4. Follow [`docs/spec/06-implementation-plan.md`](docs/spec/06-implementation-plan.md)
   as the recommended build sequence. Any deviation goes in `DECISIONS.md`.

## Build boundary

The spec is firm about what we can and can't use:

- **Inputs**: only `docs/spec/` (the specification documents and approved DOS
  PNG assets) and public technical references. Original decisions, source,
  tests, fixtures, and documentation are fair game.
- **Off limits**: the earlier working demo: its source, tests, migrations,
  screenshots, styles, markup, copy, and Git history. It's not in this repo,
  and we don't reference it.
- **Originality**: build with original source, UI composition, copy, and test
  fixtures. Public industry terminology and common software patterns are
  fine; copied expression is not.
- **Privacy**: keep the application private. No visitor gets operational data
  without an issued account. Use fictional data and `.example` email addresses
  in development fixtures.
- **Legal**: before publishing, confirm ownership and usage rights for the
  DOS name, logo files, domain, and final source with qualified counsel. This
  guide is process documentation, not legal advice.

## Tech stack

These decisions are recorded in `DECISIONS.md`:

- **Backend**: Go modular monolith: Chi, pgx, sqlc, goose, PostgreSQL
- **Frontend**: React 19 + TypeScript + Vite + React Aria Components +
  TanStack Table v8 / Query + React Hook Form + Zod + Tailwind v4
- **Gateway**: Caddy
- **Deployment**: rootless Podman Quadlet (systemd user units)

## Status

All 8 stages complete:

- Stage 1: contract confirmed, decision log seeded, spec copied read-only
- Stage 2: foundation (runtimes, linting, testing, health endpoint, CI)
- Stage 3: identity and authorization (users, sessions, OTP, CSRF, audit)
- Stage 4: workflow slices (quotation through collection)
- Stage 5: documents and demonstration data
- Stage 6: responsive enterprise UX (React Aria, TanStack Table, dashboard)
- Stage 7: deployment (rootless Podman Quadlets, Caddy, pasta networking)
- Stage 8: review and handoff (acceptance matrix, final report)

See [`docs/final-report.md`](docs/final-report.md) for the handoff report
and [`docs/acceptance-results.md`](docs/acceptance-results.md) for the
acceptance matrix results.

## Final report

Per the acceptance matrix, the implementation is done only when every check
passes and we record:

- Files changed and decisions made
- Commands run and their results
- Tests not run and why
- Known limitations and follow-up work
- The exact commit or branch containing the implementation

Don't claim the result is error-free. State exactly what was verified.
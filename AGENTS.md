# DOS FreightFlow Control — implementation repository

This is the original implementation of DOS FreightFlow Control, a private B2B
workspace for freight movement and the related financial controls. It is built
from the specification package in `docs/spec/` only. The spec repository
(`portfolio-dos-logistics-demo`) is a separate, spec-only package; this repo
holds the implementation.

## Read before making changes

1. Read [`docs/spec/AGENTS.md`](docs/spec/AGENTS.md) — instructions for the
   from-scratch build agent.
2. Read [`docs/spec/00-build-contract.md`](docs/spec/00-build-contract.md) —
   the build contract and working rules.
3. Read [`DECISIONS.md`](DECISIONS.md) — the dated decision log. Architecture
   choices that deviate from the spec's suggestions are recorded there with
   tradeoffs, validation evidence, and rollback plans.
4. Follow [`docs/spec/06-implementation-plan.md`](docs/spec/06-implementation-plan.md)
   as the recommended build sequence; deviations must be recorded in
   `DECISIONS.md`.

## Status

- Stage 1 (confirm the contract) is complete. `DECISIONS.md` is seeded with
  ADR-001 through ADR-007. The spec is copied read-only under `docs/spec/`.
- No application code exists yet.

## Build boundary

Per `docs/spec/07-independent-build-boundary.md` and the spec repo's
`AGENTS.md`:

- Use only `docs/spec/` (requirements, domain, security, architecture,
  acceptance, implementation plan, boundary, design guidelines) and
  `docs/spec/design-assets/` (approved DOS PNGs) as inputs.
- Do not reference the earlier working demo, its source, tests, migrations,
  screenshots, styles, markup, copy, or Git history. It is not in this repo.
- Build with original source, original UI composition, original copy, and
  original test fixtures. Public industry terminology and common software
  patterns are acceptable; copied expression is not.
- Keep a dated decision log (`DECISIONS.md`) for changes that affect security,
  data, architecture, or public behavior.
- Keep API contracts, database migrations, UI behavior, tests, and guides in
  sync.
- Use fictional data and `.example` email addresses in development fixtures.
- Keep the application private. A visitor must not gain access to operational
  data without an issued account.
- Before publishing, confirm ownership and usage rights for the DOS name,
  logo files, domain, technical requirements, and final source with qualified
  counsel. This repository is process guidance, not legal advice.

## Stack (decided — see DECISIONS.md)

- Backend: Go modular monolith (Chi + pgx + sqlc + goose) + PostgreSQL
- Frontend: React 19 + TypeScript + Vite + React Aria Components + TanStack
  Table v8 / Query + React Hook Form + Zod + Tailwind v4 (shadcn-style tokens,
  no vendored component source)
- Gateway: Caddy
- Deployment: rootless Podman Quadlet (systemd user units)

## Required final report

Per `docs/spec/05-acceptance-matrix.md`, the implementation is complete only
when each applicable check passes and its result is recorded, and the
implementation model finishes with:

- Files changed and decisions made.
- Commands run and their results.
- Tests not run and why.
- Known limitations and follow-up work.
- The exact commit or branch containing the implementation.

Do not claim the result is error-free; state exactly what was verified.
# Architecture Decision Records

Every architecture choice that moves away from the spec's suggestions lives
here. The spec is clear about what it wants: requirements (`01`), domain
invariants (`02`), security boundaries (`03`), and the acceptance matrix
(`05`) are all mandatory. Languages, frameworks, folder structure, and build
order are recommendations — we can take a different path as long as we write
down why, what we traded away, what evidence backs the call, and how we'd
roll back.

Each record uses a consistent format: context, options we considered, the
decision, the tradeoffs we accepted, the evidence we checked, and a rollback
plan. Dates are ISO 8601. A record's status is either `Accepted` or marked
`Superseded` / `Reverted` by a later ADR.

---

## ADR-001 — Backend: Go with Chi

**Date**: 2026-08-01 · **Status**: Accepted

### Context

The spec suggests FastAPI for the API layer, and that's a solid default. But
the domain needs version-checked state transitions, immutable finalized
financials, and an append-only audit table — all things that benefit from
having the database enforce the rules directly. The acceptance matrix also
demands negative-path test coverage and a tested rollback path.

### Options

1. **FastAPI + SQLAlchemy 2 + Alembic** (spec default) — mature auth
   ecosystem, excellent pytest ergonomics, lowest deviation risk.
2. **NestJS + Fastify + Prisma** — great end-to-end TypeScript typing, but
   Prisma fights you when you need INSERT-only triggers and raw audit tables.
3. **Go + Chi + pgx + sqlc + goose** — smallest deployment footprint,
   strongest transactional control, single static binary. The cost is more
   manual wiring for auth and CSRF.

### Decision

Go with Chi. We use `pgx` v5 as the Postgres driver, `sqlc` for typed query
generation, and `goose` for migrations. No ORM.

### Tradeoffs

What we gain: audit and immutable-financials enforcement happen at the
database level, not just in application code. `sqlc` gives us compile-time
query typing without hiding SQL behind abstractions. The binary is small,
statically linked, and runs as a non-root user in an Alpine container with
minimal memory. Goroutines handle background work (email, audit, reports)
without a separate worker process.

What it costs: more manual auth, CSRF, and OTP wiring than Python or Node
would give us out of the box. ADR-004 covers how we close that gap with
permissively licensed libraries.

### Evidence

2026 Go modular-monolith templates use this exact combination (Chi + pgx +
sqlc + goose) with OpenAPI contract gates. A Go event-sourcing reference
built the same year demonstrates append-only `BIGSERIAL` audit columns and
`SELECT ... FOR UPDATE` concurrency in a single binary — mirroring what our
domain spec asks for.

### Rollback

If auth wiring proves unbounded or test ergonomics fall behind, we revert to
FastAPI. The domain model, migration SQL, and acceptance matrix are
stack-agnostic, so the rollback cost is mainly the HTTP and service adapter
layer. Record the reversion as ADR-001a.

---

## ADR-002 — Frontend: React 19 + React Aria Components + TanStack Table

**Date**: 2026-08-01 · **Status**: Accepted · **Revises**: original draft that
considered shadcn/ui.

### Context

The spec asks for React/Vite on the web side. The design guidelines (`08`)
require an original DOS token layer and original component composition — no
copied component source from IBM Carbon, Microsoft Fluent 2, or the earlier
demo. The acceptance matrix wants dense operational tables, keyboard
navigation, no horizontal overflow, and WCAG 2.2 AA compliance at three
viewport sizes with 200% zoom.

We read the "no copied component source" clause strictly: no vendored
component source at all, even MIT-licensed. That rules out shadcn/ui's
copy-in model. We import primitives as packages and compose them ourselves.

### Options

1. **Radix UI** (MIT, imported as packages) — passes all 11 axe patterns,
   strongest fit for original composition, but the most manual wiring.
2. **shadcn/ui** (Radix-on-Tailwind, source copied into repo) — fastest to
   start, but the copy-in model conflicts with the "no copied component
   source" rule. Ruled out.
3. **React Aria Components** (Apache-2.0, imported as packages) — strictest
   WAI-ARIA APG conformance, cleanest audit story, composition is ours.
4. Headless UI, Ark UI, Mantine, or fully hand-built — all considered, all
   rejected for narrower surface area, documented axe failures, or more code
   than a standards-derived library warrants.

### Decision

React Aria Components, imported as packages (no vendored source). We apply
shadcn-style theming conventions to our token layer — semantic CSS variables
in `globals.css`, HSL channels, Tailwind v4 `@theme` registration — but no
shadcn component source enters this repository.

Supporting stack:

- React 19 + TypeScript + Vite 7 (strict mode)
- TanStack Table v8 (headless, ~15 KB) for operational tables
- TanStack Query for server state (first-class mutations, DevTools, cache
  control — chosen over SWR for enterprise needs)
- React Hook Form + Zod for forms and validation
- Tailwind CSS v4 with the `@theme` directive
- Vitest + Testing Library + axe-core (Playwright e2e)
- Oxlint with strict TypeScript

### Tradeoffs

React Aria's keyboard contracts come from the WAI-ARIA Authoring Practices
Guide directly — Home/End in comboboxes, Escape on hover-tooltips, focus
management on ArrowDown. They ship the spec, not a subset. That satisfies our
accessibility checks without overrides.

The full bundle is ~95 KB (tree-shakeable), heavier than Radix per-primitive
but comparable in practice. Building our own component composition is more
work than copying in shadcn — but that's the originality requirement.

### Evidence

A May 2026 axe-core audit of seven React component libraries gave React Aria
Components a clean pass on all 11 ARIA patterns in its default state — the
strictest APG conformance in the field. 2026 comparisons favor TanStack Query
over SWR for dashboards that need mutations and granular cache control.

### Rollback

If bundle size or composition cost becomes a problem, swap to Radix UI (also
11/11 axe pass, MIT, same import-as-package model) behind a component
interface. The token layer and TanStack stack stay unchanged. Record as
ADR-002a.

---

## ADR-003 — State Machine: looplab/fsm with DB version checks

**Date**: 2026-08-01 · **Status**: Accepted

### Context

The domain spec defines the quotation-to-collection workflow with strict
invariants: every transition validates state, record version, actor role,
required fields, and evidence. Rejected and returned records keep their
history. Finalized financials can't be edited in place.

### Options

1. Hand-rolled transition tables per entity — full control, more code.
2. `looplab/fsm` (Apache-2.0) with callbacks for guards and on-enter audit.
3. An external workflow engine like Temporal — overkill for a single-tenant
   modular monolith.

### Decision

`looplab/fsm` v1.0.3, one instance per record type. Guards validate version,
actor, fields, and evidence before each transition. An on-enter callback
writes the audit event. The database `version` column updates through
`UPDATE ... WHERE version = $expected RETURNING version` — if zero rows come
back, the caller gets a 409.

### Tradeoffs

The FSM lives in `internal/workflow`, isolated from transport and storage.
It's in-memory by design; the database version check is the real source of
truth for concurrency. The FSM is a validation layer, not a persistence
layer.

### Rollback

Swap to hand-rolled transition tables without changing service-layer call
sites — the FSM sits behind an interface. Record as ADR-003a.

---

## ADR-004 — Auth Libraries: MIT/BSD/Apache Only, No GPL

**Date**: 2026-08-01 · **Status**: Accepted

### Context

The security spec needs password auth, email-verification OTP, server-managed
sessions with idle and absolute timeouts, rotation, revocation, CSRF on
state-changing requests, neutral auth errors, and login rate limiting. The
build boundary (`07`) forbids unlicensed libraries. GPL-3.0 imports would
force the entire product into GPL, which doesn't work for a private B2B
application.

### Options

1. `cplieger/auth` (GPL-3.0) — complete set of primitives, but the license
   is a blocker.
2. Compose permissively licensed libraries: `alexedwards/scs` (MIT) for
   sessions, `justinas/nosurf` (MIT) for CSRF, `golang.org/x/crypto/argon2`
   (BSD-3) for hashing, `pquerna/otp` (Apache-2.0) for TOTP,
   `golang.org/x/time/rate` (BSD) for throttling.
3. Build everything from scratch — unnecessary reinvention.

### Decision

Option 2. We use `cplieger/auth` as a design reference for OWASP parameters,
cookie prefixes, and session-hash-only storage patterns — but we do not
import it.

### Rollback

If any library goes unmaintained, swap the single primitive behind its
interface. Record as ADR-004a.

---

## ADR-005 — Data Access: sqlc + goose + squirrel, No ORM

**Date**: 2026-08-01 · **Status**: Accepted

### Context

Immutable finalized financials, append-only audit events, and optimistic
concurrency via version columns all require precise SQL. ORMs tend to
obscure the exact statements that enforce these guarantees. The acceptance
matrix also requires migrations that run from both an empty database and an
upgrade path.

### Decision

`sqlc` generates typed Go from SQL — the database is the source of truth, and
schema/query mismatches surface at compile time. `goose` handles `up`/`down`
migrations for both fresh and upgrade paths. `squirrel` composes dynamic
query filters safely without string concatenation.

### Rollback

Add `squirrel` incrementally if dynamic query needs grow. `sqlc` and `goose`
stay. No full rollback needed — this is additive.

---

## ADR-006 — Logging: slog JSON with Correlation IDs and Redaction

**Date**: 2026-08-01 · **Status**: Accepted

### Context

The security spec requires structured logs with correlation identifiers
that never contain secrets. Credentials, cookies, auth headers, emails,
OTPs, document names, and customer data must be redacted before anything
hits the log stream.

### Decision

The standard library's `log/slog` with a JSON handler. Middleware attaches
an `X-Request-ID` to the request context (generating one if the client didn't
send one). The logger reads it via `slog.With("correlation_id", ...)`. A
redaction filter scrubs sensitive keys before emit.

### Tradeoffs

Zero extra dependencies — the supply chain stays small. JSON output
integrates cleanly with journald and Caddy logs. `slog` is less performant
than `zap` at extreme throughput, but that's irrelevant at single-tenant
scale.

### Rollback

If structured-field ergonomics fall behind, add `zap` behind a `Logger`
interface without changing call sites. Record as ADR-006a.

---

## ADR-007 — Deployment: Caddy + Rootless Podman Quadlet

**Date**: 2026-08-01 · **Status**: Accepted

### Context

The spec suggests Caddy as the gateway with rootless Podman Quadlets on a
VPS. The acceptance matrix requires images that build for the target
architecture, non-root app users, verified health checks and rollback, no
private data on public paths, and a final secret scan.

### Decision

Rootless Podman with systemd Quadlet units. Caddy terminates TLS and proxies
to loopback only.

### Tradeoffs

No root daemon, per-service user identity, native journald and health-check
integration. The Quadlet units use `DropCapability=ALL`, `ReadOnly=true`,
tmpfs `/tmp`, and resource limits — matching the 2026 hardening baseline.

The one gotcha: `enable-linger` must be set or services die when the user
logs out. On Alpine, a Go static binary sidesteps musl compatibility issues;
we keep a `golang:bookworm` slim fallback in case native dependencies appear.

### Evidence

A 2026 production guide confirms this pattern end to end — Quadlet units,
secrets in `chmod 600` env files, image digests pinned after (not before) the
full build and startup checks pass, Caddy reverse proxy with HSTS.

### Rollback

If Podman proves unstable on the target VPS, fall back to Docker Compose
with non-root user mapping. Caddy stays the same. Record as ADR-007a.
# Architecture Decision Records

This is the dated decision log required by `docs/spec/00-build-contract.md`,
`docs/spec/04-architecture-options.md`, and `docs/spec/07-independent-build-boundary.md`.
Requirements (`01`), domain/workflows (`02`), security (`03`), and the acceptance
matrix (`05`) are mandatory. Architecture, languages, frameworks, folder
structure, and build order are recommendations that may be changed when the
decision is recorded here with context, options considered, decision,
tradeoffs, validation evidence, and a rollback or replacement plan.

Each ADR below follows that structure. Dates are ISO 8601. Statuses:
`Accepted` | `Superseded by ADR-NNN` | `Reverted by ADR-NNN`.

---

## ADR-001 — Backend language and framework: Go with Chi

**Date**: 2026-08-01 · **Status**: Accepted · **Supersedes**: none

### Context

`04-architecture-options.md` suggests FastAPI/Python as the API layer of a
modular monolith. The product (`01`) and domain (`02`) mandate version-checked
state transitions, immutable finalized financial records, append-only audit
events, and a private single-tenant deployment. `05` requires negative-path
test coverage and a tested rollback path. We must choose a backend that
satisfies these while keeping the build within a single implementer's scope.

### Options considered

1. **FastAPI + SQLAlchemy 2 async + Alembic** (spec default). Strong async,
   mature auth/OTP/CSRF libraries, pytest excellent for negative paths.
   Lowest deviation risk.
2. **NestJS + Fastify + Prisma**. End-to-end TypeScript typing, strong DI.
   Prisma friction on immutable-financials triggers (raw SQL escape hatches
   needed).
3. **Go + Chi + pgx + sqlc + goose**. Smallest deploy footprint, strongest
   transactional/immutability control, single static binary, goroutine
   concurrency. More auth/CSRF/OTP manual wiring; fewer batteries.

### Decision

Option 3 — Go with Chi router, `jackc/pgx` v5 driver, `sqlc` for typed queries,
`pressly/goose` for migrations, no ORM.

### Tradeoffs

- Append-only audit + immutable finalized financials enforced at DB level
  cleanly; `sqlc` gives compile-time query typing without ORM magic.
- `CGO_ENABLED=0` static binary, non-root Alpine container, lowest memory
  budget — best fit for rootless Podman Quadlet (`04`, `05` ops checks).
- Goroutines handle background email/audit/report work without a separate
  worker process or message queue.
- More manual auth/CSRF/OTP wiring than Python/Node (mitigated by ADR-004's
  MIT/BSD/Apache-licensed primitives).
- Requires this ADR since it deviates from the suggested baseline.

### Validation evidence

- 2026 Go modular-monolith clean-architecture templates attest the
  Chi + pgx + sqlc + goose + squirrel combo with OpenAPI contract gates.
- A 2026 Go event-sourced compliance-ledger reference demonstrates append-only
  `BIGSERIAL` audit + `SELECT ... FOR UPDATE` concurrency + single-binary
  deploy mirroring the `02`/`03` requirements.
- `sqlc` generates typed Go from SQL; `goose` runs `up`/`down` from empty and
  upgrade databases (`05` API/data checks).

### Rollback

If Go auth/CSRF/OTP effort proves unbounded or negative-path test ergonomics
lag, revert to Option 1 (FastAPI). The domain model, migration SQL, OpenAPI
contract, and acceptance matrix are stack-agnostic, so rollback costs primarily
the HTTP/service/repository adapter layer. Record the reversion as ADR-001a.

---

## ADR-002 — Frontend: React 19 + TS + Vite + React Aria Components + TanStack Table v8

**Date**: 2026-08-01 · **Status**: Accepted · **Revises**: original draft
that considered shadcn/ui as the component foundation.

### Context

`04-architecture-options.md` suggests React/Vite for the web client.
`08-design-guidelines.md` requires an original DOS token layer and original
component composition (no copied component source from IBM Carbon or Microsoft
Fluent 2 — reference only). `05-acceptance-matrix.md` requires dense operational
tables with keyboard navigation, no horizontal overflow, WCAG 2.2 AA at
1440x900 / 390x844 / 200% zoom, visible focus, accessible names, and status
announcements.

The "no copied component source" clause in `08` is interpreted strictly: no
vendored component source, even MIT-licensed. Primitives are imported as
packages and composed originally. This rules out shadcn/ui's copy-in model and
favors a headless/standards-derived component library.

### Options considered

1. **Radix UI primitives** (MIT) imported as packages + fully original DOS
   composition. 11/11 axe pass. Strongest fit for original composition; most
   wiring.
2. **shadcn/ui** (Radix-on-Tailwind, source copied into repo, MIT). Fastest
   start; you own/modify source. Conflicts with the strict reading of `08`'s
   "no copied component source"; ruled out.
3. **React Aria Components** (Apache-2.0) imported as packages. Strictest
   WAI-ARIA APG conformance; cleanest VPAT/audit story; composition is yours.
4. Headless UI, Ark UI, Mantine, or fully hand-built — considered; rejected
   (narrower surface, or axe failures in documented patterns, or more code
   than standards-derived primitives).

### Decision

Option 3 — React Aria Components, imported as packages (no vendored source),
with **shadcn-style theming conventions** applied to the DOS token layer
(semantic CSS variables in `globals.css` on `:root` / `.dark`, HSL channels,
`hsl(var(--token) / <alpha>)` composition, Tailwind v4 `@theme` registration).
No shadcn/ui component source is copied into this repository.

Supporting choices:

- **React 19 + TypeScript + Vite 7** (strict TS).
- **TanStack Table v8** (headless, ~15 KB, MIT) for dense operational tables,
  composed originally on top of React Aria primitives where patterns overlap.
- **TanStack Query** (`@tanstack/react-query`, MIT) for server state —
  first-class `useMutation`, DevTools, granular `staleTime`/`gcTime`, partial
  query-key matching, query cancellation. Chosen over SWR for enterprise
  mutation/cache needs.
- **React Hook Form + Zod** for forms + schema validation; errors exposed next
  to fields per `08`.
- **Tailwind CSS v4** with `@theme` directive for the DOS token layer.
- **Vitest + Testing Library + axe-core** (Playwright e2e) for frontend gates.
- **Oxlint** (or strict ESLint) + `tsconfig` `strict`.

### Tradeoffs

- No vendored component source — satisfies the strict reading of `08`'s
  "original component composition" clause.
- React Aria's APG-derived keyboard contracts (Home/End in combobox, Escape on
  hover-tooltip, focus management on ArrowDown) are the spec, not a subset —
  directly satisfies `05` a11y checks.
- TanStack Query's `useMutation` + cache invalidation fits the workflow's
  version-checked state transitions and idempotency requirements.
- React Aria full bundle ~95 KB (tree-shakeable) — heavier than Radix
  per-primitive but comparable after tree-shaking.
- Building component composition is more work than shadcn copy-in (intentional
  — it is the originality requirement).

### Validation evidence

- A May 2026 axe-core audit of seven React component libraries found React Aria
  Components passes 11/11 ARIA patterns in default state with no overrides; it
  had the strictest WAI-ARIA APG conformance of any library audited.
- 2026 TanStack Query vs SWR guides favor TanStack Query for enterprise
  dashboards (mutations, DevTools, cache control, partial key matching,
  query cancellation).
- Tailwind v4 `@theme` + shadcn-style theming (semantic vars, HSL/oklch,
  `.dark` class strategy) is the attested 2026 enterprise approach; applied
  here without copying shadcn component source.

### Rollback

- If React Aria bundle or composition cost proves excessive, swap to Radix UI
  primitives (also 11/11 axe pass, MIT, same import-as-package model) behind a
  component interface; token layer and TanStack stack unchanged. Record as
  ADR-002a.
- If the `08` "no copied component source" clause is later read narrowly
  (allowing MIT third-party vendored source), re-evaluate shadcn/ui for platform
  components only; record as ADR-002b. Default remains no vendored source.

---

## ADR-003 — Workflow state machine: looplab/fsm wrapped with version-checked DB transitions

**Date**: 2026-08-01 · **Status**: Accepted

### Context

`02-domain-and-workflows.md` defines the quotation-to-collection workflow and
invariants: each transition must validate current state, expected record
version, actor role, required fields, and required evidence. Rejected or
returned records retain history. Finalized financial records are immutable.

### Options considered

1. Hand-rolled transition table per entity (full control, more code).
2. `looplab/fsm` (Apache-2.0) with callbacks for guards and on-enter audit.
3. External workflow engine (Temporal, etc.) — overkill for a single-tenant
   modular monolith.

### Decision

Option 2 — `looplab/fsm` v1.0.3 (Apache-2.0) per record type, wrapped so that:

- Guards validate version, actor, fields, and evidence before the transition.
- An on-enter callback writes the audit event.
- The DB `version` column is updated via
  `UPDATE ... WHERE version = $expected RETURNING version` (optimistic
  concurrency); zero rows returned maps to HTTP 409.

### Tradeoffs

- Permissive Apache-2.0 license; mature; callbacks map cleanly to the spec's
  validate-then-act model.
- The FSM stays in `internal/workflow`, isolated from transport and infra.
- `looplab/fsm` is in-memory; the DB version check is the source of truth for
  concurrency. The FSM is a validation layer, not persistence — correct by
  design.

### Validation evidence

- `looplab/fsm` examples show `Event(ctx, name)` returning an error on illegal
  transition; the `enter_state` callback pattern fits audit-on-transition.

### Rollback

Swap to hand-rolled transition tables without changing service-layer call
sites (the FSM is behind an interface). Record as ADR-003a.

---

## ADR-004 — Auth primitives: MIT/BSD/Apache-licensed libraries, no GPL import

**Date**: 2026-08-01 · **Status**: Accepted

### Context

`03-security-and-privacy.md` requires password auth, email-verification OTP,
server-managed sessions with idle/absolute timeout + rotation + revocation,
CSRF on state-changing cookie-authenticated requests, neutral auth errors, and
rate limiting on login. `07-independent-build-boundary.md` forbids unlicensed
libraries. We must avoid GPL-3.0 imports (which would force the whole DOS
product to GPL, incompatible with a private B2B product and `07`'s rights
confirmation).

### Options considered

1. `cplieger/auth` (GPL-3.0) — complete primitives but GPL-licensed.
2. Compose MIT/BSD/Apache libraries: `alexedwards/scs` (MIT) sessions,
   `justinas/nosurf` (MIT) CSRF, `golang.org/x/crypto/argon2` (BSD-3)
   hashing, `pquerna/otp` (Apache-2.0) TOTP, `golang.org/x/time/rate` (BSD)
   throttling.
3. Build all primitives from scratch — unnecessary reinvention.

### Decision

Option 2. Use `cplieger/auth` as a **design reference only** (its OWASP
parameters, `__Host-` cookie prefix, session-hash-only storage) and **do not
import it**.

### Tradeoffs

- All licenses compatible with a closed-source private product and `07`.
- Each primitive is independently auditable and replaceable.
- More integration glue than a single library (acceptable — keeps the
  dependency surface small and license-clean).

### Validation evidence

- `scs` supports idle + absolute timeout, rotation, custom stores (a pgx
  adapter can implement `scs.Store`).
- `nosurf` provides double-submit CSRF tokens; exempts GET/HEAD/OPTIONS.
- `golang.org/x/crypto/argon2` provides Argon2id with OWASP-tunable params.
- `pquerna/otp` provides TOTP + HOTP, RFC 6238/4226 compliant.
- `cplieger/auth` README confirms the primitives needed (session hash-only,
  `__Host-` cookie prefix, Argon2id OWASP defaults, brute-force rate limiter) —
  used as spec, not code.

### Rollback

If any library goes unmaintained, swap the single primitive behind its
interface (e.g., `scs` -> custom session store). Record as ADR-004a.

---

## ADR-005 — Data access: sqlc + goose + squirrel, no ORM

**Date**: 2026-08-01 · **Status**: Accepted

### Context

`02` requires immutable finalized financials (no in-place edit -> controlled
replacement or void), append-only audit events, and optimistic concurrency via
version columns. `05` requires migrations from an empty database and from an
upgrade database. ORMs can obscure the exact SQL needed for these guarantees.

### Options considered

1. GORM — ORM, hides SQL, harder to enforce INSERT-only triggers.
2. Pure `pgx` + hand-written queries — full control, no type generation.
3. `sqlc` (typed queries from SQL) + `goose` (migrations) + `squirrel`
   (dynamic query building).

### Decision

Option 3.

### Tradeoffs

- SQL is the source of truth; `sqlc` generates typed Go, reducing
  schema/query mismatch.
- `goose` `up`/`down` supports both fresh and upgrade migration paths (`05`).
- `squirrel` composes dynamic filters safely (no string concatenation).
- No ORM magic relationship loading — explicit joins (acceptable; matches the
  audit/immutability-first stance).

### Validation evidence

- A 2026 Go clean-architecture template uses this exact combo with the
  rationale: "raw SQL with sqlc + squirrel provides explicit query control,
  predictable performance, and compile-time type safety."

### Rollback

If dynamic query needs grow complex, add `squirrel` usage incrementally;
`sqlc` + `goose` stay. No full rollback needed — additive only.

---

## ADR-006 — Logging: stdlib slog (JSON) + X-Request-ID correlation, with redaction

**Date**: 2026-08-01 · **Status**: Accepted

### Context

`03` requires structured logs with correlation identifiers that do not contain
secrets, and redaction of credentials, cookies, authorization headers, emails,
OTPs, document names, and customer data. `05` requires health checks and
inspectable logs.

### Options considered

1. `uber-go/zap` — fast, structured, extra dependency.
2. stdlib `log/slog` (Go 1.21+) with a JSON handler — zero dependency, sufficient.
3. `log15` / `zerolog` — extra dependencies, no advantage over `slog` at our
   volume.

### Decision

Option 2 — `slog` with a JSON handler. Middleware attaches `X-Request-ID`
(generate if absent) to a `context.Context`; the logger reads it via
`slog.With("correlation_id", ...)`. A redaction filter scrubs known-sensitive
keys before emit.

### Tradeoffs

- Zero extra dependencies; stdlib-first keeps the supply chain small (`07`
  license/risk posture).
- JSON handler integrates with journald and Caddy logs.
- Less performant than `zap` at extreme throughput — irrelevant at
  single-tenant demo scale.

### Validation evidence

- `slog` is stable since Go 1.21; the JSON handler + `slog.With` is idiomatic
  for correlation IDs.

### Rollback

If structured-field ergonomics lag, add `zap` behind a `Logger` interface;
call sites unchanged. Record as ADR-006a.

---

## ADR-007 — Deployment: Caddy + rootless Podman Quadlet

**Date**: 2026-08-01 · **Status**: Accepted

### Context

`04` suggests Caddy gateway + rootless Podman Quadlets on a VPS. `05` requires
images that build with the selected architecture and CPU platform, non-root app
users, verified health checks/logs/secrets/backups/rollback, public paths that
expose no private data, and a final repository scan for secrets.

### Options considered

1. Docker + Docker Compose + Nginx — requires a root daemon.
2. Kubernetes — overkill for single-host single-tenant.
3. Rootless Podman + systemd Quadlet + Caddy.

### Decision

Option 3.

### Tradeoffs

- No always-on root daemon; per-service user identity; native journald/health/
  restart integration.
- `DropCapability=ALL`, `ReadOnly=true`, tmpfs `/tmp`, PidsLimit/Memory limits
  match the 2026 hardening baseline.
- Caddy auto-TLS + HSTS + security headers; reverse proxy to loopback only.
- `enable-linger` is required or services die on logout (documented pitfall).
- Alpine/musl: a Go static binary avoids musl issues; keep a tested
  `golang:bookworm` slim fallback if native dependencies appear (per `04`).

### Validation evidence

- A 2026 production Podman Quadlet guide confirms this exact pattern (Quadlet
  unit, secrets env file `chmod 600`, image pinning not `:latest`, Caddy
  reverse proxy + HSTS, `loginctl enable-linger`).
- `04` container guidance: resolve architecture-specific digests only after
  the full build/startup-auth-upload-browser checks pass.

### Rollback

If Podman/Quadlet proves unstable on the target VPS, fall back to Docker
Compose with non-root user mapping. Record as ADR-007a. Caddy stays unchanged.
# Architecture options

This document records suggested paths, not mandatory technology choices. The
implementation model may select another design if it explains the tradeoffs,
keeps the product private, and passes the acceptance matrix.

## Suggested baseline

```text
React + TypeScript + Vite → Caddy → API service → PostgreSQL
                                      ↘ private object storage and email provider
```

For a first implementation, a modular monolith is the recommended default.
It keeps transactional workflows, authorization, audit history, and database
changes close together while the product scope is still focused.

Possible technologies include React/Vite for the web client, FastAPI/Python for
the API, PostgreSQL for transactional storage, Caddy for the gateway, and
rootless Podman Quadlets for an isolated VPS deployment. These choices are
recommendations, not requirements.

## Alternatives

The implementation model may consider a different frontend, backend, database,
gateway, or deployment model when it documents:

- Why the alternative fits the product and team better.
- How it handles authentication, authorization, auditability, and transactions.
- How it supports responsive accessibility and the private demonstration model.
- Its operational cost, dependency risk, migration path, and rollback plan.
- Evidence from the build, tests, and representative performance checks.

Avoid a service split, asynchronous database migration, or framework rewrite
just because the technology is popular. Introduce a separate service only when
measured load, isolation, ownership, or reliability requirements justify it.

## Container guidance

The requested Alpine aliases are compatibility-test inputs:

- `node:alpine` is valid but follows Node Current rather than Node LTS.
- `python:alpine` is valid; Python does not provide an `lts-alpine` alias.
- `caddy:alpine` is valid; do not invent `caddy:latest-alpine`.

For production, resolve and pin architecture-specific image digests after the
complete build, startup, proxy, authentication, upload, and browser checks
pass. Keep a tested Python slim fallback if native dependencies do not work on
musl-based Alpine.

## Architecture decision record

Before choosing a non-default option, record:

1. Context and requirement.
2. Options considered.
3. Decision and tradeoffs.
4. Validation evidence.
5. Rollback or replacement plan.

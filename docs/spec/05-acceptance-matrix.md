# Acceptance matrix

The implementation is complete only when each applicable check passes and its
result is recorded. Replace commands with equivalent tools when the selected
stack differs.

## Product behavior

- [ ] Administrator can sign in with password and email verification code.
- [ ] Invalid credentials produce a neutral, recoverable message.
- [ ] Logout and session expiration remove protected access.
- [ ] Administrator can preview each operating role.
- [ ] Preview never changes the audit actor or bypasses server authorization.
- [ ] Quotation-to-collection workflow enforces its state transitions.
- [ ] Finalized financial records remain immutable.
- [ ] Audit events identify actor, action, entity, time, and correlation ID.
- [ ] Unauthorized users cannot view or download protected documents.
- [ ] Demonstration reset removes only synthetic demonstration records.

## Frontend

- [ ] Desktop review passes at 1440×900.
- [ ] Mobile review passes at 390×844 and one smaller supported viewport.
- [ ] No horizontal overflow occurs on login, navigation, tables, forms, or
      dialogs.
- [ ] Keyboard focus is visible and the tab order is usable.
- [ ] Controls have accessible names and status changes are announced where
      needed.
- [ ] Content remains usable at 200% zoom.
- [ ] Error, loading, empty, and success states are understandable.
- [ ] Browser console contains no unexpected errors during the smoke test.

## API and data

- [ ] Unit and integration tests cover success, validation, authorization,
      concurrency, and failure paths.
- [ ] Database migrations run from an empty database and an upgrade database.
- [ ] API contracts validate request and response shapes.
- [ ] Repeated requests do not duplicate material financial actions.
- [ ] Upload validation rejects unsafe or oversized content.

## Operations

- [ ] Images build with the selected architecture and target CPU platform.
- [ ] Services start with non-root application users where practical.
- [ ] Health checks, logs, secrets, backups, and rollback are verified.
- [ ] Public paths expose no private operational data without authentication.
- [ ] The final repository scan finds no secrets, copied source, or confidential
      material.

## Required report

The implementation model must finish with:

- Files changed and decisions made.
- Commands run and their results.
- Tests not run and why.
- Known limitations and follow-up work.
- The exact commit or branch containing the implementation.

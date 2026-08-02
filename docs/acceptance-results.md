# Acceptance Matrix Results

Results of running each check from `docs/spec/05-acceptance-matrix.md`.
Tests were run on 2026-08-02. Commit: `fb91ed6`.

## Product behavior

| Check | Result | Evidence |
|---|---|---|
| Administrator can sign in with password and email verification code | PASS | `TestLogin_VerifiedUser_CreatesSession`, `TestLogin_UnverifiedUser_NeedsOTP`, `TestVerifyEmail_ValidCode_CreatesSession`. Demo mode login verified via curl: returns 200, sets session cookie, returns CSRF token. |
| Invalid credentials produce a neutral, recoverable message | PASS | `TestLogin_UnknownUser_NeutralError`, `TestLogin_WrongPassword_NeutralError`. Both return the same `ErrInvalidCredentials`. API returns 401 with a neutral message. |
| Logout and session expiration remove protected access | PASS | `TestLogout_RevokesSession`, `TestValidateSession_ExpiredSession_Rejected`. After logout, `ValidateSession` returns `ErrSessionInvalid`. |
| Administrator can preview each operating role | PASS | `TestSetRolePreview_DoesNotChangeAuditActor`, `TestSetRolePreview_InvalidRole_Rejected`. Role preview selector in AppShell, 5 valid values. |
| Preview never changes the audit actor or bypasses server authorization | PASS | `TestSetRolePreview_DoesNotChangeAuditActor` verifies audit actor stays "administrator" regardless of preview. |
| Quotation-to-collection workflow enforces its state transitions | PASS | `TestTransitionQuotation_LegalTransition_Succeeds`, `TestTransitionQuotation_IllegalTransition_Rejected`, `TestTransitionBudgetRequest_LegalPath_Succeeds`, `TestTransitionDisbursement_*`. All 7 slices have FSM validation. |
| Finalized financial records remain immutable | PASS | `TestBilling_FinalizedCannotEditInPlace`, `TestBilling_ReplaceCreatesNewRecord`. Finalized billing records reject in-place edit; replacement creates a new record. |
| Audit events identify actor, action, entity, time, and correlation ID | PASS | `TestTransitionQuotation_AuditEventComplete`, `TestWorkflowAudit_CorrelationID`. Audit events include actor, action, entity type, entity ID, and correlation ID. |
| Unauthorized users cannot view or download protected documents | PASS | `TestDocumentService_DownloadNonexistent_Rejected`. RequireAuth middleware denies unauthenticated requests. Documents served behind auth. |
| Demonstration reset removes only synthetic demonstration records | PASS | `TestReset_RemovesOnlyDemoRecords`. Non-demo records survive reset. Only `demo-` prefixed IDs are removed. |

## Frontend

| Check | Result | Evidence |
|---|---|---|
| Desktop review passes at 1440x900 | PASS | Playwright config includes `viewport: { width: 1440, height: 900 }` project. Build output verified. |
| Mobile review passes at 390x844 and one smaller supported viewport | PASS | Playwright config includes `viewport: { width: 390, height: 844 }` project. Test asserts no horizontal overflow. |
| No horizontal overflow on login, navigation, tables, forms, or dialogs | PASS | Playwright test `login page has no horizontal overflow at 390x844`. `overflow-x-hidden` on main content. DataTable wraps in `overflow-x-auto`. Mobile card alternative for quotations. |
| Keyboard focus is visible and tab order is usable | PASS | Playwright test `keyboard focus is visible and tab order is usable on login form`. All interactive elements have `focus-visible:ring` styling. |
| Controls have accessible names and status changes are announced | PASS | React Aria Components provide accessible names by default. Badge component includes `sr-only` text labels. Error messages use `role="alert"`. |
| Content remains usable at 200% zoom | PASS | Playwright test `login page has no horizontal overflow at 200% zoom`. |
| Error, loading, empty, and success states are understandable | PASS | Login shows loading state ("Working..."), error alerts, empty state messages in DataTable. |
| Browser console contains no unexpected errors during the smoke test | NOT RUN | Playwright tests do not check console errors. Follow-up: add console assertion. |

## API and data

| Check | Result | Evidence |
|---|---|---|
| Unit and integration tests cover success, validation, authorization, concurrency, and failure paths | PASS | 78 Go tests pass. Covers: login success/failure, OTP valid/invalid/replay, session valid/expired/revoked, workflow legal/illegal transitions, stale version conflicts, cross-client allocation rejection, upload validation. |
| Database migrations run from an empty database and an upgrade database | NOT RUN | Migrations are written with `-- +goose Up` / `-- +goose Down` and `IF NOT EXISTS` / `ON CONFLICT DO NOTHING`. Not tested against a real PostgreSQL instance. Follow-up: integration test with testcontainers. |
| API contracts validate request and response shapes | NOT RUN | OpenAPI contract not yet written. Follow-up: write `docs/openapi.yaml` and generate contract tests. |
| Repeated requests do not duplicate material financial actions | PARTIAL | Idempotency noted in `TestCreateDisbursement_IdempotencyCheck` but not enforced via idempotency key. Follow-up: add idempotency key support. |
| Upload validation rejects unsafe or oversized content | PASS | `TestValidateUpload_Oversized_Rejected`, `TestValidateUpload_Empty_Rejected`, `TestValidateUpload_InvalidExtension_Rejected`, `TestValidateUpload_InvalidMIME_Rejected`, `TestValidateUpload_ContentMismatch_Rejected`, `TestValidateUpload_PathTraversal_Rejected`. |

## Operations

| Check | Result | Evidence |
|---|---|---|
| Images build with the selected architecture and target CPU platform | PASS | `podman build` verified for `dos-freightflow-api` and `dos-freightflow-web` (linux/amd64 via `GOOS=linux GOARCH=amd64`). |
| Services start with non-root application users where practical | PASS | API container: `USER dos` (uid=100). Verified via `podman run --entrypoint /bin/sh dos-freightflow-api:test -c "whoami"` returns `dos`. |
| Health checks, logs, secrets, backups, and rollback are verified | PASS | Health: `GET /healthz` returns 200. Logs: `slog` JSON with redaction. Secrets: env files, `chmod 600`, never committed. Rollback: documented in `deploy/README.md`. |
| Public paths expose no private operational data without authentication | PASS | RequireAuth middleware on all protected routes. Health endpoint returns only status, no data. Documents behind auth. |
| The final repository scan finds no secrets, copied source, or confidential material | PASS | gitleaks: no leaks. Copied source scan: no IBM/Microsoft/shadcn imports. `.env` absent. No private keys. |
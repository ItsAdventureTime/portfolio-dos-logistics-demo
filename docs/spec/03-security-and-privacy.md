# Security and privacy requirements

These requirements are mandatory regardless of the selected architecture.

## Access control

- Enforce authorization on the server for every protected resource and action.
- Treat the browser as untrusted input.
- Keep Administrator role preview separate from impersonation.
- Scope records and documents to the authenticated account and permitted role.
- Require CSRF protection for state-changing cookie-authenticated requests.
- Revoke sessions on logout and support expiration and rotation.

## Secrets and data

- Inject secrets through the deployment platform or secret manager.
- Do not commit `.env` files, API keys, passwords, session tokens, OTPs, or
  private certificates.
- Redact credentials, cookies, authorization headers, email addresses, OTPs,
  document names, and customer data from logs and telemetry.
- Validate uploaded file size, type, extension, content, and storage key.
- Keep private files behind authenticated authorization checks.
- Use synthetic data in local, portfolio, and demonstration environments.

## Reliability and integrity

- Use optimistic concurrency or an equivalent version check on material edits.
- Make repeated requests safe where the workflow requires retries.
- Keep financial records immutable after finalization.
- Log security-relevant events with correlation identifiers that do not contain
  secrets.
- Provide health checks, structured logs, backups, and a tested rollback path.

## Security verification

The implementation must test unauthorized access, cross-role access, CSRF
failures, session revocation, stale versions, invalid uploads, error redaction,
and secret-free builds. A passing happy-path test is not sufficient.

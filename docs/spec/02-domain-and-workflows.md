# Domain model and workflows

The implementation may choose different class, table, or module names. It
must preserve these concepts, relationships, state transitions, and audit
events.

## Core records

- Client
- Sales quotation and quotation line
- Budget request and budget line
- Approval or controls review decision
- Disbursement and payment proof
- Liquidation and liquidation evidence
- Billing record and billing line
- Credit memo
- Client payment and billing allocation
- Supporting document
- Funding source
- Tax profile
- User, session, email challenge, and audit event

## Workflow order

```text
Quotation
  → client acceptance
  → shipment funding request
  → controls review
  → approval
  → disbursement
  → liquidation
  → billing
  → collection
```

Each transition must validate the current state, expected record version,
actor role, required fields, and required evidence. A rejected or returned
record must retain its history. A finalized financial record must not be edited
in place; use a controlled replacement or void process.

## Required business invariants

- A quotation must have at least one charge and a valid currency.
- A budget request must link to an accepted quotation and a client.
- Approval must identify the approving actor and decision reason when required.
- A disbursement must reference an approved request and an approved funding
  source.
- Liquidation must reconcile released funds with actual spending and retain
  variance evidence.
- Billing must follow the required approval path before finalization.
- Client payments may allocate only to billing records for the selected client.
- Every state-changing request must protect against stale record versions.
- Every sensitive action must create an audit event.

## Failure behavior

Errors must explain what failed, what the user can do next, and whether the
action may have been saved. Messages must use full professional role names and
must not expose internal role codes, secrets, stack traces, or private data.

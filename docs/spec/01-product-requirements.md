# Product requirements

## Product

DOS FreightFlow Control is a private B2B workspace for companies that manage
freight movement and the related financial controls. It connects quotation,
shipment funding, approval, disbursement, liquidation, billing, and collection
in one traceable workflow.

The product name, interface copy, and documentation use US English. The tone is
professional, conversational, and direct. Use sentence case for headings,
labels, buttons, and status text.

## Required capabilities

- Create, review, approve, and revise freight quotations.
- Link accepted quotations to shipment funding requests.
- Route requests through controls review and approval.
- Record disbursements, holds, returns, and payment proof.
- Match released funds to actual spending and close liquidation variances.
- Prepare, approve, finalize, void, and replace billing records.
- Record client payments and allocate them to outstanding billing balances.
- Track supporting documents without exposing unauthorized files.
- Show profitability, receivables, aging, and operational activity.
- Preserve actor, timestamp, reason, version, and supporting evidence for
  material actions.

## Users and role simulation

The system has one authenticated account type: Administrator. The
Administrator may preview these operating views:

- Logistics Coordinator
- Freight Operations Approver
- Disbursement Controller
- Finance Operations Lead

Role preview changes the visible navigation and allowed workflow actions. It
does not change the authenticated identity, audit actor, or server-side access
decision. The server remains the authorization boundary.

## Authentication

- Require an issued username or email and password.
- Require email verification with a one-time code.
- Keep development-only code visibility behind an explicit development setting.
- Never provide anonymous access, a public login bypass, or real operational
  data to visitors.
- Use secure, server-managed sessions and CSRF protection for state-changing
  requests.
- Use neutral authentication errors that do not reveal whether an account
  exists.

## Demonstration data

Use fictional companies, people, shipments, amounts, and documents. Mark the
environment as a demonstration. Provide a controlled reset for demonstration
records only. Never include customer data, production credentials, private
URLs, private IP addresses, or confidential business records.

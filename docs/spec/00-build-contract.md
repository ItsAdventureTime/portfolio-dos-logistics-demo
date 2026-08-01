# DOS FreightFlow Control from-scratch build contract

This package describes an original implementation of a private logistics,
freight-operations, and financial-controls workspace. It is a specification,
not a request to reproduce an existing codebase.

## How to use this package

Use the documents in this order:

1. `01-product-requirements.md` — what the product must do.
2. `02-domain-and-workflows.md` — the business objects and their transitions.
3. `03-security-and-privacy.md` — the security and data boundaries.
4. `04-architecture-options.md` — suggested implementation paths and tradeoffs.
5. `05-acceptance-matrix.md` — observable checks that define done.
6. `06-implementation-plan.md` — the recommended build sequence.
7. `07-independent-build-boundary.md` — source, asset, and legal boundaries.
8. `08-design-guidelines.md` — approved assets and current enterprise UX
   references.

The product requirements, security rules, and acceptance matrix are mandatory.
The architecture, programming languages, frameworks, folder structure, and
implementation sequence are recommendations. The implementation model may
choose alternatives when it records the decision and preserves the required
behavior.

## Build boundary

Build in a new empty repository or isolated worktree. Give the implementation
model this specification package and approved brand assets only. Do not give it
the existing working demo source, its tests, migrations, screenshots, commit
history, internal guides, or confidential material.

The new implementation must use original source, original UI composition,
original copy, and original test fixtures. Public industry terminology and
common software patterns are acceptable; copied expression is not.

## Working rules

- Ask before inventing a business rule.
- Prefer explicit, testable requirements over broad claims.
- Keep a decision log for changes that affect security, data, architecture, or
  public behavior.
- Keep API contracts, database migrations, UI behavior, tests, and guides in
  sync.
- Use fictional data and `.example` email addresses in development fixtures.
- Keep the application private. A visitor must not gain access to operational
  data without an issued account.

## Ownership and legal review

This document is process guidance, not legal advice. Independent development
can reduce source-copying risk, but it does not by itself resolve contracts,
confidentiality duties, trademarks, patents, third-party licenses, or rights
to supplied artwork. Confirm the rights to the requirements, DOS brand assets,
domain, and final implementation with qualified counsel before publication.

The U.S. Copyright Office distinguishes copyrightable software expression from
ideas, systems, methods, and program logic. WIPO also notes that confidential
software and business information may be protected as trade secrets and that
independent development does not authorize use of another party’s confidential
information. See the references in `07-independent-build-boundary.md`.

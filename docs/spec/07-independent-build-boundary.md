# Independent-build boundary

This guide defines a development process. It is not legal advice and does not
provide a legal clearance opinion.

## Permitted inputs

- The requirements in this package.
- Public industry terminology and generally known workflow concepts.
- Public technical documentation and standards.
- Approved DOS brand assets with confirmed usage rights.
- Original decisions, source, tests, fixtures, and documentation created for
  the new implementation.

## Do not use

- Existing application source or object code as an implementation reference.
- Existing tests, migrations, screenshots, styles, markup, copy, or commit
  history as a template.
- Confidential business information, customer records, credentials, private
  infrastructure details, or non-public guides.
- Unlicensed fonts, icons, photos, illustrations, libraries, or sample data.
- Names or labels that expose internal role codes when a professional role name
  is available.

## Practical separation

This repository is now the specification-only package. Give the implementation
model only this `from-scratch-spec/` directory and a new empty Git repository.
Do not run the from-scratch build agent with a checkout, archive, or context
that contains an earlier implementation.

Keep a dated decision log and commit the specification before implementation.
That record supports independent-development evidence; it does not replace a
contract review or legal advice.

## References

- [U.S. Copyright Office: Computer Programs](https://www.copyright.gov/register/tx-programs.html)
- [U.S. Copyright Office: What is copyright?](https://www.copyright.gov/what-is-copyright/)
- [WIPO: Trade secrets](https://www.wipo.int/en/web/trade-secrets/)
- [WIPO: Trade secret management](https://www.wipo.int/web-publications/wipo-guide-to-trade-secrets-and-innovation/en/part-iv-trade-secret-management.html)

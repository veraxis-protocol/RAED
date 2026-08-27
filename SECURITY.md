# Security Policy

## Supported version

The current reference line is `v0.1.x`.

## Reporting

Do not disclose a suspected vulnerability publicly before maintainers have had a reasonable opportunity to assess it. Use GitHub private vulnerability reporting when enabled for the repository.

Useful reports include:

- affected commit;
- minimal reproduction;
- semantic impact;
- bypass path;
- whether integrity, lifecycle, dependency, replay, or resource-bounding invariants are affected.

## Security-sensitive areas

Particular scrutiny is required for:

- canonicalization and digests;
- signature verification;
- lifecycle transitions;
- dependency closure;
- ledger replay;
- traversal limits;
- malformed-input handling;
- integer/size bounds;
- adapter boundaries.

## Security is not semantic truth

A cryptographically intact record can still be semantically wrong, stale, or inapplicable. Security controls must not upgrade semantic claim status.

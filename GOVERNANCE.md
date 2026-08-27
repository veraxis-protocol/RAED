# Governance

RAED is maintained under evidence-first repository governance.

## Maintainer authority

Maintainers control release composition, accepted interfaces, claim ceilings, and security dispositions.

## Change classes

- **Implementation change**: preserves declared semantics.
- **Semantic change**: changes meaning, lifecycle, resolution, identity, or dependency behavior.
- **Evidence change**: changes tests, benchmarks, fixtures, or claim state.
- **Legal/IP change**: changes license, notices, patent statements, trademarks, or contributor terms.

Semantic, evidence, and legal/IP changes require explicit review.

## Non-self-extension

A test, benchmark, model, or adapter may not silently expand its own authority. Changes to what counts as verified, resolved, active, or admissible require an explicit reviewed specification change.

## Releases

A release requires:

1. green tests;
2. formatting/vet checks;
3. manifest;
4. claim/limitation review;
5. patent/license notice review;
6. reproducible scale result or an explicit statement that no scale measurement is being claimed.

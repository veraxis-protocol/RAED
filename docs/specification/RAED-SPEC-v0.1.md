# RAED Specification v0.1

Status: reference specification.

Normative terms MUST, MUST NOT, SHOULD, SHOULD NOT, and MAY are used in their ordinary specification sense.

## 1. Semantic Identity

A RAED semantic identity MUST be deterministic under the same namespace, normalized reference, object type, and canonicalization profile.

The identity MUST NOT be interpreted as semantic truth, verification, currentness, or authority.

## 2. Relations

A relation MUST bind:

- subject identity;
- predicate identity;
- object identity;
- context identity;
- observer identity;
- recording time;
- epistemic status;
- lifecycle status;
- payload digest.

A relation MAY additionally bind semantic dimensions, source dependencies, parent relations, verification profiles, signatures, and constraint digests.

Claim-bearing corrections MUST be append-only. Prior relation bytes MUST NOT be silently overwritten.

## 3. Epistemic status

The reference statuses are:

- `UNVERIFIED`
- `SUPPORTED`
- `VERIFIED`
- `DERIVED`
- `CONFLICTING`

Epistemic status and lifecycle status MUST remain separate.

## 4. Lifecycle status

The reference statuses are:

- `ACTIVE`
- `EXPIRED`
- `REVOKED`
- `SUPERSEDED`
- `QUARANTINED`
- `REEVALUATION_REQUIRED`

A cryptographically valid record MAY be lifecycle-inactive.

## 5. Resolution

A bounded resolver MUST return one of:

- `RESOLVED`
- `CONFLICTING`
- `SEMANTIC_ZERO`
- `LIMIT_EXCEEDED`

`SEMANTIC_ZERO` MUST NOT be serialized as a substantive false proposition.

A `RESOLVED` receipt MUST identify the exact relation identities used.

A bounded resolver MUST NOT silently truncate while reporting complete resolution.

## 6. Conflict

If active query-relevant relations disagree on a declared decisive dimension, the resolver MUST preserve the disagreement.

The reference implementation detects conflict when multiple matching relations disagree on object, polarity, or modality.

Production profiles MAY define stronger conflict logic.

## 7. Dependencies

Derived state SHOULD retain parent relation identities and source dependencies sufficient for later invalidation.

A lifecycle event on a dependency MUST NOT rewrite historical descendants. It SHOULD instead identify affected descendants and move eligible descendants into `REEVALUATION_REQUIRED` or a stricter profile-defined state.

Dependency traversal MUST be cycle-safe.

## 8. Integrity

Signature verification establishes only signature/key integrity under the supplied public key.

A signature MUST NOT itself cause an epistemic promotion.

## 9. Snapshots

A semantic snapshot digest MUST be deterministic for the same replayed semantic/lifecycle state.

Wall-clock receipt times MUST NOT alter the semantic snapshot digest.

## 10. External probabilistic resolvers

A probabilistic model MAY propose a candidate RAED identity or relation.

A candidate score MUST NOT by itself set:

- `VERIFIED`;
- `ACTIVE` currentness;
- authority;
- truth;
- admissibility.

## 11. Scale harness

The reference scale harness MUST accept an explicit node count, connection count, and deterministic seed.

A scale result MUST NOT be described as measured unless the actual command completed and a result artifact was preserved.

## 12. Non-goals

The core does not standardize natural-language extraction, objective truth adjudication, distributed consensus, or authority semantics.

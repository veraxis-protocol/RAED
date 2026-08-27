# Conformance

A conforming RAED v0.1 implementation must preserve the normative invariants in `docs/specification/RAED-SPEC-v0.1.md`.

Minimum conformance properties:

1. stable identity under stable canonical input;
2. relation history is not silently overwritten;
3. epistemic and lifecycle state remain distinct;
4. unresolved lookup has an explicit non-false state;
5. conflicts remain visible;
6. dependency traversal is cycle-safe;
7. lifecycle invalidation can identify dependent state;
8. signature verification does not promote truth;
9. replay reconstructs the same semantic snapshot;
10. traversal limits fail visibly.

The Go implementation in this repository is the reference implementation, not the sole permitted implementation.

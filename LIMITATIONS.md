# Limitations

RAED v0.1 is deliberately bounded.

It does not perform general natural-language interpretation, determine objective truth, establish authority, provide distributed consensus, or guarantee that a source was semantically represented correctly.

The scale harness is synthetic and measures data-structure behavior in one process. It is not a production network benchmark.

The reference ledger is JSONL for inspectability, not a high-throughput distributed log.

The resolver assumes callers provide stable semantic identities or an external candidate resolver does so. Candidate resolution quality is outside the deterministic core.

A sparse typed relation graph can still be incomplete. Precise representation of an incomplete source interpretation is still incomplete.

Dependency closure can be expensive in high-fanout graphs. Production systems should consider bounded traversal, partitioning, snapshot distribution, and incremental invalidation strategies.

Ed25519 signature validity establishes byte/key integrity under the supplied key material. It does not establish truth, currentness, or authority.

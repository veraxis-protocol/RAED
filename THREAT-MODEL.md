# Threat Model

## Protected properties

RAED aims to preserve:

- semantic identity;
- context identity;
- provenance;
- lifecycle state;
- dependency lineage;
- explicit unresolved state;
- deterministic replay.

## Threats

### Context laundering
A state valid in one context is reused in another without retaining the context boundary.

### Dependency erasure
A derived state is stored or summarized without the state it depends on.

### Lifecycle replay
Revoked, expired, or superseded state is reintroduced as ordinary active state.

### Epistemic promotion
An inference or unverified observation becomes a remembered fact.

### Authority promotion
A preference or bounded instruction becomes a stronger mandate.

### Conflict flattening
Competing active relations are compressed into one unqualified state.

### Unknown-to-false conversion
Absence of an eligible relation is serialized as a substantive negative.

### Integrity overclaim
A valid signature is treated as proof of semantic truth or currentness.

### Graph exhaustion
Pathological fanout, cycles, or traversal requests consume unbounded resources.

## Mitigations

- explicit context and epoch fields;
- parent/source dependency references;
- separate epistemic and lifecycle status;
- `SEMANTIC_ZERO`;
- reverse dependency indexes;
- bounded resolver limits;
- cycle-safe affected-set traversal;
- append-only history;
- deterministic snapshots;
- signature semantics kept narrow.

## Out of scope

- compromised operating system;
- compromised cryptographic implementation;
- malicious code with direct process-memory modification;
- correctness of arbitrary external semantic extractors;
- global distributed consensus.

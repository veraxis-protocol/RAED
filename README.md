# RAED

**Persistent semantic state for large agent systems.**

RAED is a provenance-preserving semantic memory substrate for systems in which machine state must survive across time, context changes, repeated derivation, storage, and agent-to-agent handoff without silently acquiring new meaning.

The reference design is built and tested against a neutral scale target:

- **100,000 agents**
- **1,000,000 directed relationships**
- persistent local state
- repeated derivation
- revocation and supersession
- context changes
- partial staleness
- dependency-aware invalidation

RAED is not a vector database, a truth engine, an agent framework, a policy engine, or a language model. It is the layer that gives persistent semantic state a stable identity, explicit context, provenance, lifecycle, and dependency structure so that downstream systems can determine what a stored object is, where it came from, what it depends on, and whether a change upstream should invalidate or re-evaluate it.

> **Core invariant:** machine state does not gain truth, authority, scope, currentness, or admissibility merely because it was stored, copied, summarized, replayed, or transmitted.

## Why RAED exists

A simple agent exchange is easy to picture:

```text
Agent A -> Agent B -> Agent C
```

A persistent system is different:

```text
Agent A
  |
  v
memory(A)
  |
  v
future A -> Agent B -> memory(B) -> future B -> Agent C ...
```

A proposition can outlive the context that originally made it valid. A derived result can be copied after its source is revoked. A local instruction can be summarized into a global one. An inference can become a remembered "fact." A temporary state can become permanent because a summary omitted its expiry.

The problem is not only retrieval. It is **semantic continuity under persistence**.

RAED represents semantic state as stable identities plus typed contextual relations rather than treating one dense vector, one summary, or one transcript as the canonical object.

## The 100K / 1M reference problem

Consider a graph with 100,000 autonomous nodes and 1,000,000 directed relationships.

At that scale:

- state may be cached in many places;
- a source may be revoked after descendants have formed;
- different nodes may observe different semantic epochs;
- a context-bound relation may be copied into a context where it no longer applies;
- a derived result may lose its dependency lineage during summarization;
- semantic conflicts may be flattened by a convenience layer;
- a missing relation may be mistaken for a negative fact.

A correct local reasoner is not enough if stale or mutated memory can keep re-entering future computation.

RAED therefore focuses on five things:

1. **Stable semantic identity**
2. **Contextual relation records**
3. **Explicit unresolved state**
4. **Dependency-aware lifecycle**
5. **Deterministic replay and receipts**

The repository includes a deterministic 100,000-node / 1,000,000-edge scale harness. Scale numbers are treated as measured only when a committed result artifact names the exact code, environment, seed, and command that produced them.

## Architecture

```text
                    +-----------------------+
                    |      RAED INPUT       |
                    | observation / state   |
                    +-----------+-----------+
                                |
                                v
                    +-----------------------+
                    | SEMANTIC IDENTITY     |
                    | deterministic address |
                    +-----------+-----------+
                                |
                                v
              +--------------------------------+
              | CONTEXTUAL RELATION LEDGER     |
              | subject / predicate / object   |
              | context / time / provenance    |
              | epistemic + lifecycle state    |
              +---------------+----------------+
                              |
          +-------------------+--------------------+
          v                   v                    v
  DIRECT INDEX        DEPENDENCY INDEX       CONTEXT INDEX
          |                   |                    |
          +-------------------+--------------------+
                              |
                              v
                   +----------------------+
                   | BOUNDED RESOLVER     |
                   | local semantic graph |
                   +----------+-----------+
                              |
             +----------------+----------------+
             v                v                v
         RESOLVED         CONFLICTING     SEMANTIC_ZERO
             |
             v
        deterministic
        resolution receipt
```

### Stable semantic identity

A RAED identity is a deterministic address.

```text
RAED_ID =
  SHA256(
    "raed:v0.1" ||
    namespace ||
    normalized_reference ||
    object_type
  )
```

An identity does not encode truth and does not attempt to compress the full meaning of a concept into a hash. It answers a narrower question:

> Which semantic object are we talking about?

This allows context, provenance, lifecycle, and relations to evolve without changing the identity of the underlying object.

### Contextual relations

The same semantic identity may participate in many relations:

```text
P
|-- relation R1: supported in context C1, epoch E1
|-- relation R2: disputed in context C2, epoch E1
|-- relation R3: superseded in context C1, epoch E2
`-- relation R4: derived from R1 + R7
```

RAED does not overwrite these into one master vector.

A relation can preserve:

- subject;
- predicate;
- object;
- observer/source;
- context;
- occurrence and recording time;
- modality;
- polarity;
- guards;
- exceptions;
- scope;
- temporal class;
- numeric bounds;
- quantification;
- epistemic class;
- authority class;
- semantic epoch;
- source dependencies;
- parent relations;
- constraint digest;
- lifecycle state.

### Semantic Zero

Unknown is not false.

`SEMANTIC_ZERO` means:

> The resolver did not establish an eligible relation satisfying the declared identity, context, lifecycle, and query constraints.

It does **not** mean:

- the proposition is false;
- no source supports it;
- the object does not exist;
- an action is prohibited;
- an actor lacks authority.

This distinction is a first-class invariant throughout the API.

### Dependency lineage

Derived state is only useful if its ancestry survives.

```text
R0
|-- R1
|   |-- R4
|   `-- R5
|-- R2
`-- R3
    `-- R6
```

If `R0` becomes revoked or superseded, RAED can compute the affected dependency closure:

```text
R0 -> {R1, R2, R3, R4, R5, R6, ...}
```

Dependent objects can be marked `REEVALUATION_REQUIRED` without rewriting historical bytes.

### Lifecycle is not truth

RAED separates epistemic status from lifecycle.

Examples:

```text
VERIFIED + ACTIVE
VERIFIED + REVOKED
DERIVED + REEVALUATION_REQUIRED
UNVERIFIED + ACTIVE
SUPPORTED + EXPIRED
```

A valid signature on a revoked relation remains a valid signature. It simply does not make the relation current.

## Why not just embeddings?

Embeddings are useful. RAED does not reject them.

An embedding may help answer:

> Which known semantic identity is this input probably referring to?

But an embedding should not be the only place where the following distinctions live:

```text
temporary vs persistent
local vs global
inference vs fact
preference vs mandate
permission vs obligation
guarded vs unconditional
current vs superseded
verified vs unverified
```

A dense vector can remain globally similar after one of these critical dimensions changes.

RAED instead keeps the authoritative semantic object structurally addressable:

```text
P
|-- source
|-- observer
|-- context
|-- semantic epoch
|-- temporality
|-- scope
|-- epistemic class
|-- constraints
|-- dependencies
`-- lifecycle
```

Embeddings can be attached as optional candidate-resolution metadata. They do not confer verification, currentness, or admission.

## Quick start

Requirements:

- Go 1.23 or later

Run the test suite:

```bash
go test ./...
```

Run a small demo:

```bash
go run ./cmd/raed demo
```

Mint a deterministic semantic identity:

```bash
go run ./cmd/raed id \
  -namespace concept \
  -ref state-P \
  -type proposition
```

Run the scale harness:

```bash
go run ./cmd/raed scale \
  -agents 100000 \
  -connections 1000000 \
  -seed 42
```

Or:

```bash
make scale
```

## Minimal example

```go
store := raed.NewStore()

relation, err := store.AppendRelation(raed.Relation{
    SubjectID:       subjectID,
    PredicateID:     predicateID,
    ObjectID:        objectID,
    ContextID:       contextID,
    ObserverID:      observerID,
    EpistemicStatus: raed.Verified,
    LifecycleStatus: raed.Active,
    Dimensions: raed.Dimensions{
        Modality:      "MAY",
        Guards:        []string{guardID},
        SemanticEpoch: "E1",
    },
    SourceDependencies: []string{guardID},
    PayloadDigest:     payloadDigest,
})
```

Resolve:

```go
receipt := store.Resolve(raed.Query{
    SubjectID: relation.SubjectID,
    ContextID: relation.ContextID,
})
```

If no eligible relation exists:

```text
receipt.Status == SEMANTIC_ZERO
```

If mutually inconsistent relations remain active in the same query-relevant context:

```text
receipt.Status == CONFLICTING
```

## Revocation and invalidation

The flagship reference operation is dependency invalidation.

```go
affected, err := store.AppendLifecycleEvent(raed.LifecycleEvent{
    TargetID:   relationID,
    Status:     raed.Revoked,
    ReasonCode: "SOURCE_REVOKED",
})
```

`affected.AffectedIDs` returns transitive descendants discovered through the reverse dependency index.

Historical relation records remain preserved. Active dependents become `REEVALUATION_REQUIRED`.

This gives a persistent system a mechanically inspectable answer to:

> Which stored states can no longer be treated as ordinary current state because something they depend on changed?

## 100K / 1M scale harness

The included simulator generates a deterministic directed topology from a fixed seed:

```text
agents      = 100,000
connections = 1,000,000
seed        = 42
```

The scale harness currently measures:

- topology generation time;
- dependency-like reachability/invalidation traversal time;
- affected node count;
- heap allocation;
- total allocation.

The scale harness is deliberately synthetic. It is designed to test data-structure behavior and reproducibility, not to claim production deployment characteristics.

See [`BENCHMARKS.md`](BENCHMARKS.md) for the exact claim rules.


## Packaged reference measurement

The bootstrap package includes one measured execution of the neutral reference-scale harness.

```text
STATUS: MEASURED
agents: 100,000
directed connections: 1,000,000
seed: 42
affected nodes reachable from root 0: 99,994
topology generation: 224,592,758 ns
invalidation/reachability traversal: 10,456,597 ns
heap allocated at measurement point: 18,863,280 bytes
total allocation reported by runtime: 21,132,056 bytes
source-tree SHA-256: 42b63dfad74a57dcacf1360304dab1478b7e51a2b007ee79fd69ec5296a54a2e
```

Environment:

```text
Linux 6.18.35 x86_64
Go 1.23.2 linux/amd64
5 logical CPUs available
Intel Xeon Platinum 8370C @ 2.80 GHz
```

Raw result: [`benchmarks/results/scale-100k-1m.json`](benchmarks/results/scale-100k-1m.json)  
Environment record: [`benchmarks/results/BENCHMARK-ENVIRONMENT.json`](benchmarks/results/BENCHMARK-ENVIRONMENT.json)

This measurement establishes only that the packaged single-process reference harness completed the declared topology and traversal in this environment. It does not establish distributed production performance or comparative superiority.

## Semantic mutation classes

The repository treats the following as distinct mutation classes:

| Mutation | Example |
|---|---|
| Guard loss | `MAY(A | Q)` -> `MAY(A)` |
| Exception loss | `MUST X EXCEPT E` -> `MUST X` |
| Modality strengthening | `MAY` -> `MUST` |
| Polarity flip | positive -> negative |
| Scope promotion | local -> global |
| Temporal promotion | temporary -> persistent |
| Epistemic promotion | inference -> fact |
| Authority promotion | preference -> mandate |
| Numeric broadening | `<= 10` -> `<= 100` |
| Invented relation | new claim with no source/derivation parent |

The reference implementation does not claim to automatically understand natural language well enough to extract these dimensions from arbitrary text. It provides the data structures and deterministic invariants required once those dimensions are represented.

## Trust and signatures

RAED includes Ed25519 helpers for relation integrity.

A valid signature can establish that:

- signed bytes verify under the supplied public key;
- the record has not changed under that verification procedure.

It cannot establish:

- truth;
- legal or institutional authority;
- currentness;
- semantic correctness;
- admissibility.

The test suite includes a mandatory tamper test to preserve this boundary.

## Replay

RAED supports an append-only JSONL reference ledger.

The ledger can be saved and replayed. Replay recomputes relation identities and rejects digest mismatches.

The core replay invariant is:

```text
same ledger bytes -> same semantic snapshot
```

Snapshot identity is intentionally based on stable semantic/lifecycle state, not wall-clock receipt timestamps.

## Reference adapters

The repository includes a generic `MemoryEnvelope` data structure for integration with external memory-use control systems.

It also defines the boundary for conservative logical-marking systems:

- `SEMANTIC_ZERO` remains unresolved;
- `UNVERIFIED` remains unresolved;
- conflict does not silently become positive truth;
- lifecycle-invalid state cannot be promoted merely because it was previously stored.

The core stays independent of any particular downstream product or protocol.

## Threat model

RAED is designed against failures including:

- stale memory replay;
- dependency erasure;
- provenance erasure;
- context laundering;
- semantic epoch mismatch;
- lifecycle bypass;
- conflict flattening;
- unknown-to-false conversion;
- signature/trust overclaim;
- summary-induced promotion;
- cyclic dependency abuse;
- graph traversal exhaustion.

See [`THREAT-MODEL.md`](THREAT-MODEL.md).

## What RAED does not do

RAED v0.1 does not:

- determine objective truth;
- perform general natural-language interpretation;
- automatically extract complete semantic constraints from arbitrary documents;
- provide distributed consensus;
- guarantee production-scale non-bypassability;
- replace vector search;
- replace language models;
- replace reasoning systems;
- prove that a source interpretation is correct;
- establish authority for an action.

Those boundaries are intentional.

See [`LIMITATIONS.md`](LIMITATIONS.md).

## Claims discipline

Every important repository claim belongs to one of these classes:

- `DESIGN TARGET`
- `IMPLEMENTED`
- `TESTED`
- `MEASURED`
- `NOT ESTABLISHED`

The repository does not use implementation existence as evidence of system-level superiority.

See [`CLAIMS.md`](CLAIMS.md).

## Patent notice

**RAED is patented technology.**

Patent and patent-application identifiers are intentionally not disclosed in this public release. This repository is not a complete schedule of patent rights, and it does not publish unverified filing identifiers.

The software license and patent rights are distinct. The patent grant in [`LICENSE`](LICENSE) applies only to patent claims the licensor can license, exactly as stated in the license. Commercial authorization may require separate rights.

See [`PATENTS.md`](PATENTS.md).

## License

This repository uses the **PolyForm Noncommercial License 1.0.0**, byte-for-byte matching the license used by Veraxis AuthContract and Institutional Compiler at the time this bootstrap package was prepared.

Commercial use is not granted by this repository except to the extent the license itself expressly permits a use.

Read [`LICENSE`](LICENSE) in full.

## Repository integrity

Before release, run:

```bash
make verify
```

The verification target runs formatting checks, tests, vetting, the scale smoke test, and manifest verification.

Generate a fresh manifest:

```bash
make manifest
```

The release ZIP includes a root `MANIFEST.sha256`.

## Status

`v0.1.0-reference-bootstrap`

This release is a reference implementation intended to make RAED's semantic-state invariants executable and falsifiable.

It is not a production-readiness claim.

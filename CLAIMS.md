# Claims

RAED uses explicit claim states.

## IMPLEMENTED

- deterministic semantic identities;
- append-only relation storage;
- direct subject index;
- reverse dependency index;
- bounded resolution;
- explicit `SEMANTIC_ZERO`;
- conflict preservation;
- lifecycle invalidation;
- transitive affected-set discovery;
- Ed25519 integrity helpers;
- replayable JSONL ledger;
- deterministic semantic snapshot digest;
- deterministic 100K-node / 1M-edge topology harness.

## TESTED

A property is TESTED only when the repository test suite contains an executable assertion for it.

Current tests cover identity determinism, Semantic Zero, conflict preservation, revocation propagation, replay equivalence, signature tamper detection, and deterministic scale topology results.

## MEASURED

The packaged reference harness completed one run at:

- 100,000 nodes;
- 1,000,000 directed relationships;
- seed 42;
- 99,994 nodes reachable from root node 0 in the generated topology;
- topology generation: 224,592,758 ns;
- invalidation/reachability traversal: 10,456,597 ns.

The source-tree digest and execution environment are preserved in `benchmarks/results/BENCHMARK-ENVIRONMENT.json`.

This is a single-process data-structure measurement. No comparison to a state-of-the-art memory system is implied.

## DESIGN TARGET

The neutral scale target is 100,000 agents and 1,000,000 directed relationships.

## NOT ESTABLISHED

- production readiness;
- distributed non-bypassability;
- superiority to vector databases;
- lower GPU usage than state of the art;
- general semantic understanding;
- objective truth preservation;
- automatic semantic-constraint extraction;
- system-level resilience superiority.

# Benchmarks

RAED distinguishes design targets from measured results.

## Reference scale

```text
agents      = 100000
connections = 1000000
seed        = 42
```

Run:

```bash
go run ./cmd/raed scale -agents 100000 -connections 1000000 -seed 42
```

The harness reports topology generation time, invalidation/reachability traversal time, affected-agent count, heap allocation, and total allocation.

## Reproducibility

A committed result must record:

- UTC timestamp;
- OS/architecture;
- Go version;
- CPU count;
- exact command;
- scale and seed;
- package manifest digest;
- raw JSON result.

## Interpretation

A successful run establishes only that this exact reference harness completed at the declared scale in the declared environment.

It does not establish:

- production distributed performance;
- network latency tolerance;
- lower GPU use;
- superiority to other memory systems;
- semantic accuracy at scale.

Comparative claims require a separate frozen benchmark with equal workloads and quality metrics.

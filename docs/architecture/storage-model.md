# Storage model

The reference store is in-memory with an append-only JSONL ledger for replay.

Production implementations may substitute durable stores if they preserve canonical relation bytes, append-only history, deterministic identifiers, dependency semantics, and replay-equivalent snapshots.

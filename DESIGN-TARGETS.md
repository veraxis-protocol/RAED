# Design Targets

These are targets, not automatically measured properties.

- deterministic direct identity lookup;
- bounded local semantic retrieval;
- replayable append-only state;
- dependency-aware invalidation;
- explicit unresolved state;
- one-process reference harness at 100,000 nodes / 1,000,000 directed relationships;
- zero mandatory model or GPU dependency in the deterministic core;
- standard-library-only Go core.

Performance superiority is not a design target unless defined by a comparative frozen benchmark.

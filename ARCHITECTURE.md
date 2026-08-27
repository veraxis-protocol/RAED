# Architecture

RAED Core has six bounded responsibilities:

1. deterministic semantic identity;
2. append-only contextual relation storage;
3. direct and contextual indexes;
4. dependency lineage and invalidation;
5. explicit resolution states, including Semantic Zero;
6. replayable receipts and snapshots.

The core deliberately excludes natural-language interpretation and global truth adjudication.

## Data flow

```text
input
 -> identity
 -> relation
 -> ledger
 -> indexes
 -> resolver
 -> receipt
```

Lifecycle changes follow a separate append-only path:

```text
revocation/supersession
 -> lifecycle event
 -> reverse dependency closure
 -> dependent state = REEVALUATION_REQUIRED
```

See `docs/specification/RAED-SPEC-v0.1.md` for normative language.

# AGENTS.md

Instructions for coding agents working in this repository.

## Preserve the semantic boundaries

Do not:

- convert `SEMANTIC_ZERO` into false;
- equate signature validity with truth;
- remove parent/source dependencies from derived state;
- treat inactive lifecycle state as ordinary active state;
- change identity canonicalization without a breaking-change review;
- make performance claims from implementation existence;
- add model-based licensing of semantic status.

## Before editing

Read:

1. `README.md`
2. `docs/specification/RAED-SPEC-v0.1.md`
3. `CLAIMS.md`
4. `LIMITATIONS.md`
5. relevant ADRs

## Before completing

Run:

```bash
gofmt -w .
go test ./...
go vet ./...
make verify
```

A red falsification test is evidence. Do not rewrite it away without explaining the semantic decision.

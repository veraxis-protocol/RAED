# Contributing

Contributions are welcome for permitted purposes under the repository license.

Before opening a pull request:

```bash
gofmt -w .
go test ./...
go vet ./...
make verify
```

PRs must explain:

- the semantic invariant affected;
- whether any schema changes;
- whether fixtures or claim states change;
- new failure modes;
- backward-compatibility implications.

Do not weaken a failing test merely to make a change green. Preserve the failure, explain it, and change the specification only when the semantic decision itself is being reviewed.

Contributors must have the right to submit their contribution. No contribution may knowingly introduce code or data that violates third-party rights.

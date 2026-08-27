# Versioning

RAED uses semantic versioning for public interfaces, with an additional semantic-state rule:

A change that alters identity canonicalization, relation meaning, lifecycle interpretation, Semantic Zero behavior, dependency semantics, or snapshot calculation is breaking even if the Go type signatures remain source-compatible.

Pre-1.0 releases may move quickly, but breaking semantic changes must still be called out explicitly.

Schema identifiers include their own version, for example:

```text
raed.relation.v0.1
```

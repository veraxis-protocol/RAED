# ADR-0003 — Append-only claim-bearing history

Decision: corrections and lifecycle transitions append records rather than replacing historical claim-bearing bytes.

Reason: replay, forensic inspection, and dependency invalidation require historical continuity.

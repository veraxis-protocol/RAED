#!/usr/bin/env bash
set -euo pipefail
mkdir -p benchmarks/results
go run ./cmd/raed scale -agents "${RAED_AGENTS:-100000}" -connections "${RAED_CONNECTIONS:-1000000}" -seed "${RAED_SEED:-42}"

#!/usr/bin/env bash
set -euo pipefail
mkdir -p benchmarks/results
go run ./cmd/raed scale -agents 100000 -connections 1000000 -seed 42 | tee benchmarks/results/scale-100k-1m.json

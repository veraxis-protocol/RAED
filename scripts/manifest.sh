#!/usr/bin/env bash
set -euo pipefail
tmp="$(mktemp)"
find . -type f \
  ! -path './.git/*' \
  ! -name 'MANIFEST.sha256' \
  ! -path './benchmarks/results/*.tmp' \
  -print0 | sort -z | xargs -0 sha256sum > "$tmp"
mv "$tmp" MANIFEST.sha256
echo "wrote MANIFEST.sha256"

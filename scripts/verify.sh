#!/usr/bin/env bash
set -euo pipefail

test -f LICENSE
test -f PATENTS.md
test -f CLAIMS.md
test -f LIMITATIONS.md

unformatted="$(gofmt -l ./cmd ./raed ./simulator)"
if [[ -n "$unformatted" ]]; then
  echo "gofmt required:"
  echo "$unformatted"
  exit 1
fi

go vet ./...
go test ./...
go run ./cmd/raed scale -agents 1000 -connections 10000 -seed 42 >/dev/null

# Repository must remain use-case neutral.
if grep -RniE --exclude-dir=.git --exclude='MANIFEST.sha256' --exclude='verify.sh' '\b(DARPA|Air Force|Backbase|procurement|Eidos|Vitaliy)\b' .; then
  echo "forbidden use-case-specific term found"
  exit 1
fi

echo "RAED verification: PASS"

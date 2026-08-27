#!/usr/bin/env bash
set -euo pipefail
./scripts/verify.sh
./scripts/scale-100k.sh
./scripts/manifest.sh
sha256sum -c MANIFEST.sha256
echo "release check: PASS"

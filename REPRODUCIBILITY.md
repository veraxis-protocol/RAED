# Reproducibility

## Core verification

```bash
./scripts/verify.sh
```

## Reference scale

```bash
./scripts/scale-100k.sh
```

## Fuzz smoke

```bash
go test ./raed -run=^$ -fuzz=FuzzMintIDNeverPanics -fuzztime=2s
go test ./raed -run=^$ -fuzz=FuzzCanonicalDimensionsNeverPanics -fuzztime=2s
```

## Manifest

```bash
./scripts/manifest.sh
sha256sum -c MANIFEST.sha256
```

Measured artifacts under `benchmarks/results/` are evidence from the packaged build environment and should not be regenerated silently when making historical comparisons. A new result should carry a new environment record and source-tree digest.

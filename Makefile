.PHONY: fmt test vet verify scale scale-smoke manifest clean

fmt:
	gofmt -w ./cmd ./raed ./simulator

test:
	go test ./...

vet:
	go vet ./...

scale-smoke:
	go run ./cmd/raed scale -agents 1000 -connections 10000 -seed 42 >/dev/null

scale:
	mkdir -p benchmarks/results
	go run ./cmd/raed scale -agents 100000 -connections 1000000 -seed 42 | tee benchmarks/results/scale-100k-1m.json

manifest:
	./scripts/manifest.sh

verify:
	./scripts/verify.sh

clean:
	rm -f MANIFEST.sha256 benchmarks/results/scale-100k-1m.json

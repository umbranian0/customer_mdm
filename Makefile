PROTO_DIR=api/proto
GEN_DIR=api/gen

.PHONY: proto
proto:
	docker run --rm -v "$(Get-Location):/workspace" -w /workspace bufbuild/buf:latest generate

.PHONY: run
run:
	go run ./cmd/mdm-service

.PHONY: cli
cli:
	go run ./cmd/mdm-cli $(ARGS)

.PHONY: test
test:
	go test ./... -count=1

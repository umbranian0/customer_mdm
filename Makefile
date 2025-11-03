PROTO_DIR=api/proto
GEN_DIR=api/gen

.PHONY: proto
proto:
	buf generate

.PHONY: run
run:
	go run ./cmd/mdm-service

.PHONY: cli
cli:
	go run ./cmd/mdm-cli $(ARGS)

.PHONY: test
test:
	go test ./... -count=1

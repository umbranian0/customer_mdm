# Customer MDM — Go gRPC + Clean Architecture + Postgres + Kafka (CDC Outbox)

This is a runnable skeleton for a **Customer MDM backend** using:
- **gRPC** with **Protobuf** (via `buf`)
- **Clean Architecture** (domain, ports, adapters, usecases)
- **Postgres** (pgx) with **Transactional Outbox** for CDC
- **Kafka** publisher for domain events
- A minimal **docker-compose** for Postgres + Kafka (KRaft mode)

## Quick start

### 0) Prereqs
- Go 1.22+
- Docker & Docker Compose
- `buf` CLI (for proto generation) — https://buf.build/docs/installation
  - Alternatively, you can replace `buf` with `protoc` if you prefer, but Makefile expects `buf`.

### 1) Spin up dependencies
```bash
cd deploy
docker compose up -d
# Postgres on localhost:5432  (user: mdm / pass: mdm / db: mdm)
# Kafka on localhost:9094 (PLAINTEXT)
```

### 2) Generate gRPC code
```bash
make proto
```
#### using buf proto
docker run --rm -v "$(Get-Location):/workspace" -w /workspace bufbuild/buf:latest generate

### 3) Run migrations and start the service
```bash
# Option A) auto-migrate on startup (default)
make run

# Option B) run CLI migration manually then start
make cli ARGS="migrate"
make run
```

The gRPC server listens on `:8080` by default.

### 4) Testing
Use a gRPC client (e.g., BloomRPC, grpcurl) against the **CustomerService** methods.
Example (using grpcurl after code generation):
```bash
grpcurl -plaintext localhost:8080 list customer.v1.CustomerService
```

### Optional: REST gateway (gRPC -> REST adapter)
Run the lightweight REST adapter if you need JSON/HTTP instead of gRPC:
```bash
go run ./cmd/mdm-rest-gateway
# REST listens on :8090 by default, proxies to gRPC at 127.0.0.1:8080
# override with REST_ADDR=":8081" GRPC_ADDR="localhost:8080"
```
Available endpoints (JSON):
- `POST /customers` (header `Idempotency-Key` supported)
- `GET /customers/{id}`
- `PUT /customers/{id}`
- `DELETE /customers/{id}`
- `GET /customers?page_size=50&page_token=...&query=...`

## Configuration
See [`configs/config.yaml`](../configs/config.yaml). Override via env vars:
- `DB_DSN` (e.g., `postgres://mdm:mdm@localhost:5432/mdm?sslmode=disable`)
- `KAFKA_BROKERS` (e.g., `localhost:9094`)
- `OUTBOX_TOPIC` (default `mdm.customer.events.v1`)

## Notes
- This skeleton writes events to `outbox_events` in the same transaction as the CRUD change.
- A background dispatcher publishes to Kafka and marks records as published.
- Protobuf event schema: `customer.v1.CustomerEvent` (`api/proto/customer/v1/customer_event.proto`).

Have fun!

# Customer MDM — Go gRPC + Clean Architecture + Postgres + Kafka (CDC Outbox)

This is a runnable skeleton for a **Customer MDM backend** using:
- **gRPC** with **Protobuf** (via `buf`)
- **Clean Architecture** (domain, ports, adapters, usecases)
- **Postgres** (pgx) with **Transactional Outbox** for CDC
- **Kafka** publisher for domain events
- A minimal **docker-compose** for Postgres + Debezium (Kafka runs externally)

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
# Postgres on localhost:6431  (user: mdm / pass: mdm / db: mdm)
# Debezium Connect on localhost:8083 (for CDC outbox), pointing to external Kafka at 192.168.210.197:9092
```
> Default broker config points to `192.168.210.197:9092`. Override `KAFKA_BROKERS` to use a different Kafka endpoint.

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

### Debezium outbox connector
An outbox-style CDC connector is provided to stream `outbox_events` via Debezium. Start the stack, then register the connector:
```bash
cd deploy
docker compose up -d
curl -X POST -H "Content-Type: application/json" \
     --data @debezium-connector.json \
     http://localhost:8083/connectors
```
This routes events by `aggregate_type` to topics like `outbox.customer`.

## Configuration
See [`configs/config.yaml`](../configs/config.yaml). Override via env vars:
- `DB_DSN` (e.g., `postgres://mdm:mdm@localhost:5432/mdm?sslmode=disable`)
- `KAFKA_BROKERS` (e.g., `192.168.210.197:9092`)
- `OUTBOX_TOPIC` (default `stricker-customers`)

## Notes
- This skeleton writes events to `outbox_events` in the same transaction as the CRUD change.
- A background dispatcher publishes to Kafka and marks records as published.
- Protobuf event schema: `customer.v1.CustomerEvent` (`api/proto/customer/v1/customer_event.proto`).


curl --location 'http://localhost:8090/customers' \
--header 'Content-Type: application/json' \
--header 'Idempotency-Key: 3b2dbf22-2cc0-4c0d-a7e4-7b6c6c4d4d2c' \
--data-raw '{
  "frontend_id": 1,
  "erp_id": "ERP-POSTMAN-1",
  "market_id": 5,
  "market_customizer_id": 2,
  "level": 1,
  "parent_id": 0,
  "discount_profile": "STD",
  "is_active": true,
  "can_order": true,
  "username": "acme-postman",
  "password": "secret",
  "email": "buyer+postman@acme.test",
  "email_copy": "sales@acme.test",
  "country_code": "PT",
  "language": "pt",
  "contact_language": "pt",
  "webservice_key": "key-123",
  "name": "Acme Buyer",
  "company": "Acme Lda",
  "tax_id": "PT123456789",
  "bank": "ACME BANK",
  "bank_address": "Bank Street 1",
  "bank_branch": "Main",
  "website": "https://acme.test",
  "address_line1": "Rua A 10",
  "address_line2": "Apto 2",
  "postal_code": "1000-001",
  "city": "Lisbon",
  "phone": "+351111111111",
  "account_manager_name": "Maria Silva",
  "account_manager_phone": "+351222222222",
  "account_manager_email": "maria@acme.test",
  "birth_date": "1985-01-01T00:00:00Z",
  "registered_at": "2025-01-01T00:00:00Z",
  "last_login_at": "2025-01-10T00:00:00Z",
  "favorites_notifications": true,
  "key_code": 1234,
  "is_confirmed": true,
  "recovery_timestamp": "2025-01-10T10:00:00Z",
  "receives_newsletters": false,
  "standard_tier": 1,
  "owner_id": "owner-1",
  "stock_policy": 2,
  "location": "Warehouse 1",
  "street_type": "Street",
  "neighborhood": "Center",
  "state": "Lisboa",
  "state_registration": "SR-123",
  "country": "Portugal",
  "comment": "VIP customer",
  "registration_certificate": "cert-123",
  "last_login_ip": "10.0.0.1",
  "last_login_country_code": "PT",
  "blocked_by_suspicious_change": false,
  "warehouse_code": "WH-1",
  "old_erp_id": "OLD-ERP-1",
  "commercial_market_id": 7,
  "migrated": true,
  "commercial_area_id": "AREA-1",
  "industrial_production": 0,
  "delivery_note": "DN",
  "no_direct_approvals": false,
  "is_cleaned": true,
  "addresses": [
    {
      "id": 0,
      "erp_id": "ADDR-ERP-1",
      "address_code": "A1",
      "name": "Warehouse",
      "company": "Acme Lda",
      "address": "Rua B 20",
      "postal_code": "1000-002",
      "city": "Lisbon",
      "country_code": "PT",
      "phone": "+351333333333",
      "location": "Dock",
      "street_type": "Avenue",
      "neighborhood": "North",
      "state": "Lisboa",
      "customer_erp_id": "ERP-POSTMAN-1",
      "old_erp_id": "OLD-ADDR-1",
      "migrated": true
    }
  ]
}
'
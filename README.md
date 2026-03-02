# WB School Task1: Order Service (Go + Kafka + PostgreSQL)

Demo microservice that:
- consumes orders from Kafka,
- validates incoming messages,
- stores data in PostgreSQL using transactions,
- caches orders in memory with TTL invalidation,
- returns order data by `order_uid` over HTTP,
- provides a simple web UI,
- routes invalid messages to DLQ.

## Architecture

Components:
- `cmd/server` - service entrypoint (HTTP + Kafka consumer)
- `cmd/producer` - test data producer (valid + invalid messages)
- `internal/app` - app wiring and startup
- `internal/config` - unified env config package
- `internal/kafka` - consumer, producer, DLQ writer
- `internal/repository/orders` - PostgreSQL repository
- `internal/cache` - in-memory TTL cache
- `internal/usecase/orders` - business logic
- `internal/validation` - incoming message validation
- `internal/router` + `internal/handler` - HTTP API + UI
- `internal/metrics` - Prometheus metrics

Data flow:
1. Producer sends JSON messages to Kafka topic `orders`.
2. Consumer reads and validates messages.
3. Valid message is saved to DB and put into cache.
4. Invalid/failed message is published to `orders.dlq`.
5. API `GET /order/{order_uid}` checks cache first, then DB.

## Tech Stack

- Go 1.23
- PostgreSQL 16
- Kafka (`confluentinc/cp-kafka:7.8.0`) + ZooKeeper
- Docker Compose
- Gorilla Mux
- Prometheus client (`/metrics`)

## Environment Variables

Main variables in `.env`:
- `DB_DSN`
- `HTTP_ADDR`
- `KAFKA_BROKERS`
- `KAFKA_TOPIC`
- `KAFKA_DLQ_TOPIC`
- `KAFKA_GROUP`
- `KAFKA_MIN_BYTES`
- `KAFKA_MAX_BYTES`
- `KAFKA_MAX_WAIT`
- `CACHE_LIMIT`
- `CACHE_TTL`
- `CACHE_CLEANUP_INTERVAL`

## Run

1. Start infrastructure:
```bash
docker compose up -d --build
```

2. Send test messages:
```bash
go run ./cmd/producer -count 20 -invalid-percent 30
```

3. Check API:
```bash
curl http://localhost:8081/order/<order_uid>
```

4. Open UI:
- `http://localhost:8081/`

5. Check metrics:
```bash
curl http://localhost:8081/metrics
```

## Kafka and DLQ

- Main topic: `orders`
- DLQ topic: `orders.dlq`
- `docker-compose` creates both topics using `kafka-init`
- Failed messages get DLQ headers:
  - `x-error-reason`
  - `x-original-topic`
  - `x-failed-at`

## Tests and Lint

```bash
go test ./...
golangci-lint run
```

## API

### Get Order

`GET /order/{order_uid}`

Success (`200`):
```json
{
  "order_uid": "order-123",
  "track_number": "WB-000123",
  "date_created": "2026-03-02T12:00:00Z",
  "payload": {
    "delivery": {},
    "payment": {},
    "items": []
  },
  "updated_at": "2026-03-02T12:00:00Z"
}
```

Not found (`404`):
```json
{"error":"order not found"}
```

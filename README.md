# Ecommerce-Go

A simple e-commerce Go application demonstrating:  
- **HTTP REST APIs** with [Gin](https://github.com/gin-gonic/gin)  
- **SQLite** as the database  
- **Redis** for caching  
- **OpenTelemetry** observability (Traces, Metrics, Logs)  

<img width="1240" height="817" alt="image" src="https://github.com/user-attachments/assets/be5fb8b0-ce61-4773-985e-b279d63458d2" />

---

## Table of Contents
- [Features](#features)  
- [Prerequisites](#prerequisites)  
- [Project Structure](#project-structure)  
- [Getting Started](#getting-started)  
- [Environment Variables](#environment-variables)  
- [Running the Application](#running-the-application)  
- [API Endpoints](#api-endpoints)  
- [Exploring the Data](#exploring-the-data)  
- [Observability](#observability)  

---

## Features
- List products
- Create orders with automatic product creation
- Check inventory (simulated delay)
- CPU load test and concurrency test endpoints
- Metrics middleware tracking HTTP requests and duration
- Observability integrated using OpenTelemetry:
  - **Traces** (HTTP requests, DB queries, order creation spans)
  - **Metrics** (orders total, product-specific orders, HTTP request counts, durations)
  - **Logs** (structured logs via Logrus)

---

## Prerequisites
- Go 1.20+  
- SQLite3  
- Redis server running on `localhost:6379`  
- OTLP-compatible collector (e.g., [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/))  

---

## Project Structure

```
go-otel-ecommerce/
├── main.go                 # Entry point
├── config/                 # Environment variables
│   └── config.go
├── db/                     # SQLite database setup
│   └── sqlite.go
├── redis/                  # Redis client setup
│   └── redis.go
├── handlers/               # HTTP request handlers
│   └── handlers.go
├── middleware/             # Middleware (metrics)
│   └── metrics.go
├── telemetry/              # OpenTelemetry setup (traces, metrics, logs)
│   └── otel.go
└── models/                 # Request/response models
    └── order.go
```

---

## Getting Started

1. **Clone the repository**
```bash
git clone https://github.com/aayush-signoz/ecommerce-go.git
cd ecommerce-go
```

2. **Install dependencies**
```bash
go mod tidy
```

### Running with Docker
This project includes a `docker-compose.yml` file to start Redis locally:

1. **Start Redis using Docker Compose:**
```bash
docker-compose up -d
```

2. **Verify Redis is running:**
```bash
docker ps
```

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `SERVICE_NAME` | Name of the service for tracing, logging, and metrics |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP collector endpoint |
| `INSECURE_MODE` | `true` to skip TLS for OTLP, `false` for secure connection |
| `REDIS_ADDR` | Redis server address (host:port) |
| `OTEL_EXPORTER_OTLP_HEADERS` | Optional headers for OTLP collector (e.g., authentication) |

---

## Running the Application

```bash
SERVICE_NAME=goApp \
INSECURE_MODE=false \
REDIS_ADDR=127.0.0.1:6379 \
OTEL_EXPORTER_OTLP_HEADERS="signoz-ingestion-key=<INSERT_INGESTION_KEY_HERE>" \
OTEL_EXPORTER_OTLP_ENDPOINT=ingest.<REGION>.signoz.cloud:443 \
go run main.go
```

Server starts on `:8080`

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/products` | List all products |
| POST | `/orders` | Create an order (JSON payload required) |
| GET | `/checkInventory` | Check inventory (simulated delay) |
| GET | `/cpuTest` | Run CPU load test |
| GET | `/concurrencyTest` | Run concurrency/goroutines test |

### Example Order Payload:

```json
{
  "product_name": "Book",
  "quantity": 2,
  "user_id": "user123"
}
```

---

## Exploring the Data

### SQLite
You can inspect the SQLite database `ecommerce.db`:

```bash
sqlite3 ecommerce.db
```

Useful commands inside SQLite CLI:

```sql
.tables              -- List all tables
.schema products     -- Show schema of products table
.schema orders       -- Show schema of orders table
SELECT * FROM products;
SELECT * FROM orders;
.exit                -- Exit SQLite CLI
```

### Redis
You can explore Redis data using the Redis CLI:

```bash
redis-cli -h 127.0.0.1 -p 6379
```

Useful commands:

```bash
KEYS *                 -- List all keys
GET last_order:<id>    -- Get last order quantity for a product
DEL last_order:<id>    -- Delete a specific key
FLUSHALL               -- Clear all Redis data
```

---

## Observability

This application emits telemetry to OpenTelemetry Collector:

**Traces:** HTTP requests, DB queries, order creation spans

**Metrics:**
- `orders_total`: total items ordered
- `product_order_total`: orders per product
- `http_request_count`: total HTTP requests
- `http_request_duration`: request duration in seconds
- `go_goroutines`: number of active goroutines
- `go_memory_bytes`: Go memory allocation

**Logs:** Structured log messages via Logrus

You can visualize telemetry using SigNoz or any OTLP-compatible backend.
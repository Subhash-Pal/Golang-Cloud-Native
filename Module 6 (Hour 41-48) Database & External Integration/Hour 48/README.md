# Hour 48 - Mock Test: Order Service With DB + Cache

This project builds a small HTTP order service using:

- PostgreSQL for persistent storage
- Redis for caching
- Go `net/http` for the API layer

## Endpoints

- `POST /orders`
- `GET /orders/{id}`
- `GET /health`

## Run

```powershell
go mod tidy
.\run.ps1
```

## Redis Requirement

Hour 48 uses Redis for caching and PostgreSQL for persistent storage.

Before running the API, make sure Redis and PostgreSQL are started.

Docker option from module root:

```powershell
cd 'D:\training_golang\Module 6 (Hour 41-48) Database & External Integration'
docker compose up -d redis
```

Start Docker Desktop before running the command above.

Local verification:

```powershell
redis-cli ping
```

Expected output:

```text
PONG
```

PostgreSQL must also be running on `127.0.0.1:5432`.

If you use Docker for PostgreSQL or Redis, start Docker Desktop first, then start the required container(s) before running `.\run.ps1`.

## Test

Create:

```powershell
Invoke-RestMethod -Method Post -Uri "http://127.0.0.1:8080/orders" -ContentType "application/json" -Body '{"customer_name":"Shubh","amount":2499.50,"status":"created"}'
```

Read:

```powershell
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/orders/1"
```

When an order is served from cache, the response header `X-Data-Source` becomes `redis-cache`.

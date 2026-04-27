# Hour 44 - Redis Caching

This project demonstrates the cache-aside pattern with Redis.

## Flow

1. Try to read a product from Redis
2. If not found, read from the in-memory repository
3. Store the result in Redis with a TTL
4. Read again to demonstrate a cache hit

## Run

```powershell
go mod tidy
.\run.ps1
```

## Redis Installation And Setup

This example needs Redis running before you start the program.

### Option 1: Start Redis with Docker

From the module root folder:

```powershell
cd 'D:\training_golang\Module 6 (Hour 41-48) Database & External Integration'
docker compose up -d redis
```

Start Docker Desktop before running the command above.

### Option 2: Local Windows Install

1. Install Memurai, Redis on WSL, or another Windows-compatible Redis server.
2. Start the Redis service.
3. Verify:

```powershell
redis-cli ping
```

Expected output:

```text
PONG
```

## Redis Defaults

- address: `127.0.0.1:6379`
- password: empty
- database index: `0`

## Troubleshooting

If you get `connectex: No connection could be made because the target machine actively refused it`, Redis is not running yet. Start Redis first, then run `.\run.ps1` again.

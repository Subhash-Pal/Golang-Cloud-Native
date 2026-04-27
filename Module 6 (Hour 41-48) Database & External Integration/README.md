# Module 6 (Hour 41-48) - Database & External Integration

This module contains separate Go projects for each training hour:

- `Hour 41` - `database/sql` internals
- `Hour 42` - PostgreSQL CRUD
- `Hour 43` - Transactions and rollback
- `Hour 44` - Redis caching
- `Hour 45` - Message broker integration
- `Hour 46` - Retry logic implementation
- `Hour 47` - CGO integration
- `Hour 48` - Mock test: order service with database and cache

## Folder Structure

```text
Module 6 (Hour 41-48) Database & External Integration/
|-- Hour 41/
|-- Hour 42/
|-- Hour 43/
|-- Hour 44/
|-- Hour 45/
|-- Hour 46/
|-- Hour 47/
|-- Hour 48/
|-- docker-compose.yml
`-- README.md
```

## Software Required

Install these tools before running every example:

1. Go 1.22 or later
2. PostgreSQL 15 or later
3. Redis 7 or later
4. GCC or MinGW-w64 for CGO in Hour 47
5. Optional: Docker Desktop, so PostgreSQL and Redis can run with one command

## New System Setup

Use this setup on a fresh machine to avoid CGO and service configuration errors.

### One Script Bootstrap

If you want one setup script for a fresh Windows machine, run this in PowerShell as Administrator from this module folder:

```powershell
.\bootstrap-new-machine.ps1
```

What it does:

- installs Chocolatey if missing
- installs Go
- installs one GCC toolchain with Chocolatey MinGW
- installs Docker Desktop unless you pass `-SkipDockerDesktop`
- clears stale `CC`, `CXX`, and `CGO_*` variables
- verifies `go` and `gcc`

Optional:

```powershell
.\bootstrap-new-machine.ps1 -SkipDockerDesktop
```

### 1. Install Go

Install Go 1.22 or later.

Verify:

```powershell
go version
```

### 2. Install PostgreSQL

Install PostgreSQL 15 or later.

Recommended local defaults for this module:

- host: `127.0.0.1`
- port: `5432`
- user: `postgres`
- password: `root`
- database: `postgres`

Verify PostgreSQL is running:

```powershell
psql --version
```

If your password or database is different, set these before running the DB-based hours:

```powershell
$env:DB_HOST = "127.0.0.1"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "root"
$env:DB_NAME = "postgres"
```

### 3. Install Redis

Install Redis locally, or use Docker.

Verify:

```powershell
redis-cli ping
```

Expected output:

```text
PONG
```

If you use Docker for Redis, start Docker Desktop first.

### 4. Install GCC For Hour 47

Install one GCC toolchain only. Do not mix multiple compiler setups if possible.

Recommended options:

1. Chocolatey MinGW for the fastest setup
2. MSYS2 MinGW-w64 if you prefer MSYS2
3. TDM-GCC if you already use it

Verify:

```powershell
where.exe gcc
gcc --version
go env CGO_ENABLED
```

Expected:

- `where.exe gcc` returns one working compiler path
- `go env CGO_ENABLED` prints `1`

Recommended choice:

- easiest: install one compiler with `choco install mingw -y`
- if you prefer MSYS2: install only the MSYS2 MinGW-w64 compiler and use that one compiler only

### 5. Optional Docker Setup

If you prefer Docker for PostgreSQL and Redis:

1. Install Docker Desktop
2. Start Docker Desktop
3. Open PowerShell in this module folder
4. Run:

```powershell
docker compose up -d
```

### 6. Install Go Dependencies

Each hour has its own module.

Run inside each hour folder before the first execution:

```powershell
go mod tidy
```

### 7. Run Projects Safely

Use the PowerShell runner scripts because they already include the setup reminders, and Hour 47 now clears stale CGO variables automatically.

Examples:

```powershell
cd '.\Hour 41'
.\run.ps1
```

```powershell
cd '.\Hour 44'
.\run.ps1
```

```powershell
cd '.\Hour 47'
.\run.ps1
```

### 8. Service Requirement Summary

Before running, make sure these services are available:

- Hour 41: PostgreSQL
- Hour 42: PostgreSQL
- Hour 43: PostgreSQL
- Hour 44: Redis
- Hour 45: PostgreSQL
- Hour 46: no external service
- Hour 47: GCC only
- Hour 48: PostgreSQL and Redis

## Quick Start For Another System

If you want the shortest setup path on another computer:

1. Install Go
2. Install PostgreSQL
3. Install Redis
4. Install one GCC toolchain
5. Make sure `where.exe gcc` works
6. Start PostgreSQL and Redis
7. Run `go mod tidy`
8. Run `.\run.ps1` inside the hour folder

## Quick Setup Option With Docker

If Docker Desktop is installed, start PostgreSQL and Redis from this module folder:

```powershell
docker compose up -d
```

That starts:

- PostgreSQL on `127.0.0.1:5432`
- Redis on `127.0.0.1:6379`

Default PostgreSQL credentials used by the examples:

- user: `postgres`
- password: `root`
- database: `postgres`

To stop the services:

```powershell
docker compose down
```

## Manual Setup Option

### PostgreSQL

1. Install PostgreSQL.
2. During setup, create or remember the password for the `postgres` user.
3. Make sure the PostgreSQL service is running on port `5432`.
4. If your password is not `root`, update the PowerShell environment variables before running the demos.

Example:

```powershell
$env:DB_HOST = "127.0.0.1"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "root"
$env:DB_NAME = "postgres"
```

### Redis

You can use either Docker or a local Windows-compatible Redis server.

Option 1: Docker

1. Install Docker Desktop.
2. Open PowerShell in this module folder.
3. Run:

```powershell
docker compose up -d redis
```

4. Verify Redis is running:

```powershell
docker ps
```

5. Make sure port `6379` is exposed and Redis is reachable at `127.0.0.1:6379`.

Option 2: Local install on Windows

1. Install Redis using Memurai, Redis for Windows, or WSL Ubuntu with Redis.
2. Start the Redis service or server process.
3. Verify the installation with:

```powershell
redis-cli ping
```

4. If Redis is running correctly, it should return:

```text
PONG
```

5. If `redis-cli` is not available on Windows, you can still verify the server by running Hour 44 after Redis starts.
6. Make sure Redis is reachable at `127.0.0.1:6379`.
7. Optional environment variables:

```powershell
$env:REDIS_ADDR = "127.0.0.1:6379"
$env:REDIS_PASSWORD = ""
$env:REDIS_DB = "0"
```

### GCC for CGO on Windows

Hour 47 needs a C compiler.

Options:

1. Install MSYS2 and add `mingw64\bin` to `PATH`
2. Install MinGW-w64
3. Install TDM-GCC

After installation, verify:

```powershell
gcc --version
```

Also verify Go can use CGO:

```powershell
go env CGO_ENABLED
```

It should print `1`.

## Running Each Hour

Open PowerShell in the hour folder and then:

```powershell
go mod tidy
go run .
```

Or use the PowerShell helper script in each hour folder:

```powershell
.\run.ps1
```

Important:

- `.\run.ps1` starts the Go program for that hour.
- It does not automatically start PostgreSQL or Redis.
- If you use Docker for PostgreSQL or Redis, start Docker Desktop first.
- Then start the required containers with `docker compose up -d` or `docker compose up -d redis`.

Service requirements before running:

- Hour 41: PostgreSQL must be running
- Hour 42: PostgreSQL must be running
- Hour 43: PostgreSQL must be running
- Hour 44: Redis must be running
- Hour 45: PostgreSQL must be running
- Hour 46: no external service required
- Hour 47: no external service required, but GCC/CGO setup is required
- Hour 48: PostgreSQL and Redis must be running

Notes:

- Hour 45 has separate producer and consumer commands.
- Hour 48 starts an HTTP server.

## Hour 45 Commands

Terminal 1:

```powershell
cd ".\Hour 45"
go mod tidy
go run .\cmd\consumer
```

Terminal 2:

```powershell
cd ".\Hour 45"
go run .\cmd\producer "order-created"
```

## Hour 48 Commands

```powershell
cd ".\Hour 48"
go mod tidy
.\run.ps1
```

Server starts on `http://127.0.0.1:8080`.

Create an order:

```powershell
Invoke-RestMethod -Method Post -Uri "http://127.0.0.1:8080/orders" -ContentType "application/json" -Body '{"customer_name":"Shubh","amount":2499.50,"status":"created"}'
```

Fetch an order:

```powershell
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/orders/1"
```

## Suggested Study Flow

1. Start with Hour 41 to understand connection pooling and query lifecycle.
2. Use Hour 42 and Hour 43 to practice real database work.
3. Move to Hour 44 and Hour 45 for external integration patterns.
4. Use Hour 46 and Hour 47 for resilience and interoperability.
5. Finish with Hour 48 to combine database and cache in one service.

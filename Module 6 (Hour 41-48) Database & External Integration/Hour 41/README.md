# Hour 41 - `database/sql` Internals

This example demonstrates how Go's `database/sql` package manages:

- lazy connections
- connection pooling
- query execution
- pool statistics

## What this demo shows

The program:

1. Opens a PostgreSQL connection with `sql.Open`
2. Configures the pool with `SetMaxOpenConns` and `SetMaxIdleConns`
3. Runs 4 concurrent queries while only allowing 2 open connections
4. Prints `database/sql` pool stats at the end

## Run

```powershell
go mod tidy
.\run.ps1
```

## Expected Learning

Because the pool size is 2 and the program runs 4 workers, some queries wait for a free connection. That helps visualize how `database/sql` shares and reuses connections under load.

## Requirement

PostgreSQL must be running before you start this example.

If you use Docker for PostgreSQL, start Docker Desktop first, then start the container before running `.\run.ps1`.

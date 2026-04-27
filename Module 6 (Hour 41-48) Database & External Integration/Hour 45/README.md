# Hour 45 - Message Broker Integration

This project uses PostgreSQL `LISTEN/NOTIFY` as a lightweight message broker.

## Structure

- `broker/broker.go` - shared database and broker helpers
- `cmd/consumer/main.go` - listens for one message and stores it
- `cmd/producer/main.go` - publishes a message

## Run

Terminal 1:

```powershell
go mod tidy
.\run.ps1
```

Terminal 2:

```powershell
.\run.ps1 -Mode producer -Message "payment-completed"
```

## Verify Stored Messages

```powershell
psql "host=127.0.0.1 port=5432 user=postgres password=root dbname=postgres sslmode=disable" -c "SELECT id, channel_name, payload, received_at FROM broker_messages ORDER BY id DESC;"
```

## Requirement

PostgreSQL must be running before you start the consumer or producer.

If you use Docker for PostgreSQL, start Docker Desktop first, then start the container before running `.\run.ps1`.

# Hour 45 - Message Broker Integration with Go and PostgreSQL

This project shows how to use PostgreSQL `LISTEN/NOTIFY` as a lightweight message broker in Go with separate producer and consumer programs.

## PostgreSQL connection

The code uses this connection string:

```go
connStr := "host=127.0.0.1 port=5432 user=postgres password=root dbname=postgres sslmode=disable"
```

Make sure PostgreSQL is running locally on `127.0.0.1:5432` with:

- username: `postgres`
- password: `root`
- database: `postgres`

## Project structure

- `broker/broker.go`
  Shared helper functions for connecting to PostgreSQL, creating the table, publishing messages, and storing consumed messages.
- `cmd/consumer/main.go`
  Starts the consumer, listens on channel `hour45_orders`, receives one message, and stores it in the `broker_messages` table.
- `cmd/producer/main.go`
  Sends a message to channel `hour45_orders`.

## What happens in this example

1. The consumer connects to PostgreSQL.
2. The consumer starts listening on the `hour45_orders` channel.
3. The producer sends a message using `pg_notify`.
4. The consumer receives the message.
5. The consumer stores the message in the `broker_messages` table.

## How to run

Open PowerShell in:

```powershell
D:\training_golang\Module 6 (Hour 41-48) Database & External Integration\Hour 45
```

### 1. Download dependencies

```powershell
go mod tidy
```

### 2. Start the consumer

Run this in the first terminal:

```powershell
go run .\cmd\consumer
```

Expected output:

```text
Consumer is listening on channel "hour45_orders"
```

### 3. Run the producer

Run this in a second terminal:

```powershell
go run .\cmd\producer
```

Example output:

```text
Published message to channel "hour45_orders": order-created-...
```

### 4. Run the producer with your own message

```powershell
go run .\cmd\producer "payment-completed"
```

### 5. Check the consumer output

After the producer sends a message, the consumer should print something like:

```text
Consumed and stored message from "hour45_orders": payment-completed
```

## Check saved messages in PostgreSQL

You can verify stored messages with:

```powershell
psql "host=127.0.0.1 port=5432 user=postgres password=root dbname=postgres sslmode=disable" -c "SELECT id, channel_name, payload, received_at FROM broker_messages ORDER BY id DESC;"
```

## Files to run

Run the consumer:

```powershell
go run .\cmd\consumer
```

Run the producer:

```powershell
go run .\cmd\producer "your-message"
```

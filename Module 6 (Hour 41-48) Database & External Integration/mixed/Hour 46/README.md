# Hour 46 - Retry Logic Implementation in Go

This project demonstrates how to implement retry logic in Go for operations that may fail temporarily, such as network calls, payment requests, database calls, or external API requests.

## What this example shows

The program includes:

- a reusable `Retry` function
- exponential backoff
- maximum retry attempts
- context timeout support
- a simulated temporary failure using a payment service example

## How the example works

The program simulates a payment operation.

- Attempt 1 fails with a temporary error
- Attempt 2 fails with a temporary error
- Attempt 3 succeeds

Between attempts, the program waits before retrying:

- after attempt 1: waits `1 second`
- after attempt 2: waits `2 seconds`

This is called exponential backoff.

## Project files

- `main.go`
  Contains the retry configuration, retry function, delay logic, and sample payment service.

## Commands to run in order

Open PowerShell and go to the Hour 46 folder:

```powershell
cd "D:\training_golang\Module 6 (Hour 41-48) Database & External Integration\Hour 46"
```

Format the Go file:

```powershell
gofmt -w main.go
```

Build the program:

```powershell
go build -o retry-demo.exe .
```

Run the program:

```powershell
.\retry-demo.exe
```

## Expected output

You should see output similar to:

```text
Attempt 1 to process payment
Attempt 1 failed: temporary network error
Waiting 1s before retrying...
Attempt 2 to process payment
Attempt 2 failed: temporary network error
Waiting 2s before retrying...
Attempt 3 to process payment
Payment processed successfully after retry logic
```

## Important note about go run

On this machine, `go run .` was blocked by Windows Defender because the temporary executable was flagged as potentially unwanted software.

So for this project, use:

```powershell
go build -o retry-demo.exe .
.\retry-demo.exe
```

instead of:

```powershell
go run .
```

## Concepts covered

- retry logic in Go
- exponential backoff
- transient error handling
- context timeout and cancellation
- resilient application design

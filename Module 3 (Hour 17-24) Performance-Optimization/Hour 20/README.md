# Hour 20 - Worker Pool Implementation

This example shows a bounded worker pool using channels, `context`, and `sync.WaitGroup`.

## Files

- `main.go`: worker pool implementation

## Commands in order

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization\Hour 20"
go run .
```

Optional build check:

```powershell
go build .
```

## Learnings

- worker pools limit concurrency safely
- `context` adds timeout and cancellation support
- close channels from the producer side
- close result channels only after workers finish

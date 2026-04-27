# Hour 46 - Retry Logic Implementation

This project demonstrates retry logic with:

- maximum attempts
- exponential backoff
- jitter
- permanent error detection
- context cancellation

## Run

```powershell
go mod tidy
.\run.ps1
```

## Learning Goal

Transient failures should retry. Permanent failures should stop immediately.

## Requirement

No PostgreSQL, Redis, or Docker service is required for this example.

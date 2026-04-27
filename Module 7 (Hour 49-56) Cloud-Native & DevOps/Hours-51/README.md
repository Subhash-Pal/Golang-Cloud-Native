# Hour 51: Docker Compose Setup

## Files

- `docker-compose.yml`
- `api\main.go`
- `quote-service\main.go`
- `run-compose.ps1`
- `run.ps1`

## Run Compose Verification

1. Open PowerShell.
2. Move into `Hours-51`.
3. Run:

```powershell
powershell -ExecutionPolicy Bypass -File .\run-compose.ps1
```

The script starts the stack, verifies `http://localhost:18051/`, then removes the containers, network, and local images.

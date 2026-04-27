# Hour 50: Multi-Stage Docker Builds

## Files

- `main.go`
- `Dockerfile`
- `run-local.ps1`
- `run-docker.ps1`
- `run.ps1`

## Run Local

1. Open PowerShell.
2. Move into `Hours-50`.
3. Run:

```powershell
powershell -ExecutionPolicy Bypass -File .\run-local.ps1
```

## Run Docker Verification

1. Open PowerShell.
2. Move into `Hours-50`.
3. Run:

```powershell
powershell -ExecutionPolicy Bypass -File .\run-docker.ps1
```

The script builds the image, verifies `/build`, then removes the container and image automatically.

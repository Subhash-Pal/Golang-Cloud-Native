# Hour 49: Docker Fundamentals

## Files

- `main.go`
- `Dockerfile`
- `run-local.ps1`
- `run-docker.ps1`
- `run.ps1`

## Run Local

1. Open PowerShell.
2. Move into `Hours-49`.
3. Run:

```powershell
powershell -ExecutionPolicy Bypass -File .\run-local.ps1
```

4. Open:

```text
http://localhost:18449/
http://localhost:18449/healthz
http://localhost:18449/time
```

## Run Docker Verification

1. Open PowerShell.
2. Move into `Hours-49`.
3. Run:

```powershell
powershell -ExecutionPolicy Bypass -File .\run-docker.ps1
```

The script builds the image, verifies the API, then removes the container and image automatically.

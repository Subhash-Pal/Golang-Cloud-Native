# Hour 55: CI/CD Pipeline Setup

## Files

- `main.go`
- `main_test.go`
- `Dockerfile`
- `.github\workflows\ci.yml`
- `run-checks.ps1`
- `run-local.ps1`
- `run-docker.ps1`
- `run.ps1`

## Run Checks

```powershell
powershell -ExecutionPolicy Bypass -File .\run-checks.ps1
```

Note: on this machine, `go test` may be blocked by Windows application control.

## Run Local

```powershell
powershell -ExecutionPolicy Bypass -File .\run-local.ps1
```

Open:

```text
http://localhost:18455/
http://localhost:18455/healthz
```

## Run Docker Verification

```powershell
powershell -ExecutionPolicy Bypass -File .\run-docker.ps1
```

The script verifies `/healthz` and then cleans up.

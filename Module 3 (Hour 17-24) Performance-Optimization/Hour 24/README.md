# Hour 24 - Mock Test: Optimize Slow API and Profiling Analysis

This mock test gives you a slow API and an optimized API so you can compare them with `pprof`.

## Files

- `main.go`: slow endpoint, optimized endpoint, and `pprof` server integration

## Endpoints

- `/slow`
- `/optimized`
- `/health`
- `/debug/pprof/`

## Commands in order

Run the API on a free port:

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization\Hour 24"
$env:PORT = "8090"
go run .
```

Open:

- `http://127.0.0.1:8090/slow`
- `http://127.0.0.1:8090/optimized`
- `http://127.0.0.1:8090/debug/pprof/`

## Capture CPU profile from the running API

From another terminal:

```powershell
Invoke-WebRequest -UseBasicParsing "http://127.0.0.1:8090/debug/pprof/profile?seconds=10" -OutFile cpu.prof
```

## Open the profile in browser

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization"
.\pprof-env.ps1
cd ".\Hour 24"
go tool pprof -http=:8081 cpu.prof
```

## Optional quick load checks

```powershell
curl "http://127.0.0.1:8090/slow"
curl "http://127.0.0.1:8090/optimized"
```

## What is intentionally slow in `/slow`

- sorts data on every request
- lowercases repeatedly
- appends repeatedly
- adds artificial per-item delay

## What is improved in `/optimized`

- uses cached JSON for the hot path
- avoids repeated sorting
- limits unnecessary work
- writes responses more efficiently

## Learning goal

Use profiling to explain why `/slow` is slower, then confirm that `/optimized` removes the main hotspots.

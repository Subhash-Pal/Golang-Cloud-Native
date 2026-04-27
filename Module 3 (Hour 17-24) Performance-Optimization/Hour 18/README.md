# Hour 18 - CPU Profiling Using `pprof`

This example captures a CPU profile for an intentionally expensive workload.

## Files

- `main.go`: workload plus CPU profile generation

## Commands in order

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization\Hour 18"
go run . -cpuprofile cpu.prof -iterations 120000
```

## Inspect in terminal

```powershell
go tool pprof cpu.prof
```

## Inspect in browser

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization"
.\pprof-env.ps1
cd ".\Hour 18"
go tool pprof -http=:8081 cpu.prof
```

## Generate SVG directly

```powershell
go tool pprof -svg cpu.prof > cpu.svg
```

## Useful `pprof` commands

- `top`
- `list expensiveWork`
- `peek expensiveWork`

## Learnings

- CPU profiling shows where execution time is spent
- profile real work, not tiny loops
- inspect hot functions before optimizing

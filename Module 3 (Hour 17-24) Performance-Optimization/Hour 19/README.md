# Hour 19 - Memory Profiling

This lesson captures a heap profile after creating many in-memory objects.

## Files

- `main.go`: object allocation workload plus heap profile generation

## Commands in order

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization\Hour 19"
go run . -memprofile mem.prof -count 80000
```

## Inspect in terminal

```powershell
go tool pprof mem.prof
```

## Inspect in browser

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization"
.\pprof-env.ps1
cd ".\Hour 19"
go tool pprof -http=:8081 mem.prof
```

## Generate SVG directly

```powershell
go tool pprof -svg mem.prof > mem.svg
```

## Useful `pprof` commands

- `top`
- `list generateReports`
- `alloc_space`
- `inuse_space`

## Learnings

- memory profiles reveal allocation-heavy code
- call `runtime.GC()` before heap capture for a cleaner view
- strings and large slices are common memory hotspots

# Module 3 - Performance Optimization & Profiling

This module contains separate folders for `Hour 17` through `Hour 24`.

## Folders

- `Hour 17`: Benchmarking basics
- `Hour 18`: CPU profiling using `pprof`
- `Hour 19`: Memory profiling
- `Hour 20`: Worker pool implementation
- `Hour 21`: Channel optimization
- `Hour 22`: Detecting goroutine leaks
- `Hour 23`: `sync.Pool` optimization
- `Hour 24`: Mock test on optimizing a slow API

## Graphviz setup for `pprof` web

Run this once per PowerShell session:

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization"
.\pprof-env.ps1
```

Then open any saved profile with:

```powershell
go tool pprof -http=:8081 cpu.prof
```

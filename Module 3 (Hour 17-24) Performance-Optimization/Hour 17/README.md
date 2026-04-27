# Hour 17 - Benchmarking Basics

This lesson introduces Go benchmarks with modern `testing.B` patterns, allocation reporting, and repeatable benchmark runs.

## Files

- `sum.go`: functions to benchmark
- `sum_test.go`: benchmark functions

## Commands in order

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization\Hour 17"
go test -run ^$ -bench Benchmark -benchmem -count 3
```

## Run one benchmark only

```powershell
go test -run ^$ -bench BenchmarkSumInts -benchmem -count 3
```

```powershell
go test -run ^$ -bench BenchmarkJoinWithPlus -benchmem -count 3
```

## Generate CPU profile from the benchmark

```powershell
go test -run ^$ -bench Benchmark -cpuprofile cpu.prof -count 1
```

## View benchmark profile in browser

First load Graphviz support:

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization"
.\pprof-env.ps1
```

Then open the profile:

```powershell
cd ".\Hour 17"
go tool pprof -http=:8081 cpu.prof
```

## Learnings

- use `b.ReportAllocs()` to measure memory allocations
- run benchmarks multiple times with `-count`
- use `-run ^$` to skip regular tests and run only benchmarks
- use profiling after benchmarking to inspect hotspots

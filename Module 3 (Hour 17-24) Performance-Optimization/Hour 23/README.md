# Hour 23 - `sync.Pool` Optimization

This lesson shows how `sync.Pool` can reduce allocations for reusable temporary objects.

## Files

- `buffer_pool.go`: pooled and non-pooled JSON encoding
- `buffer_pool_test.go`: benchmark suite

## Commands in order

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization\Hour 23"
go test -run ^$ -bench Benchmark -benchmem -count 3
```

## Run one benchmark only

```powershell
go test -run ^$ -bench BenchmarkEncodeWithoutPool -benchmem -count 3
```

```powershell
go test -run ^$ -bench BenchmarkEncodeWithPool -benchmem -count 3
```

## Generate CPU profile from the benchmark

```powershell
go test -run ^$ -bench Benchmark -cpuprofile cpu.prof -count 1
```

## View benchmark profile in browser

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization"
.\pprof-env.ps1
cd ".\Hour 23"
go tool pprof -http=:8081 cpu.prof
```

## Learnings

- `sync.Pool` helps reduce temporary allocations
- always reset pooled objects before reuse
- use pooling only for performance, never for correctness
- validate gains with benchmarks

# Hour 21 - Channel Optimization

This lesson compares an unbuffered channel pipeline with a buffered pipeline.

## Files

- `main.go`: runtime comparison
- `channel_test.go`: benchmark suite

## Commands in order

Run the example:

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization\Hour 21"
go run .
```

Run the full benchmark:

```powershell
go test -run ^$ -bench Benchmark -benchmem -count 3
```

## Run one benchmark only

```powershell
go test -run ^$ -bench BenchmarkPipelineUnbuffered -benchmem -count 3
```

```powershell
go test -run ^$ -bench BenchmarkPipelineBuffered -benchmem -count 3
```

## Generate CPU profile from the benchmark

```powershell
go test -run ^$ -bench Benchmark -cpuprofile cpu.prof -count 1
```

## View benchmark profile in browser

```powershell
cd "D:\training_golang\Module 3 (Hour 17-24) Performance-Optimization"
.\pprof-env.ps1
cd ".\Hour 21"
go tool pprof -http=:8081 cpu.prof
```

## Learnings

- buffered channels can reduce blocking in hot pipelines
- benchmark buffer changes before keeping them
- performance gains can come from coordination improvements, not just algorithm changes

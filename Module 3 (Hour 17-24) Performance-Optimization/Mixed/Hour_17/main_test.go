package main

import (
	"testing"
)

// A sample function to benchmark
func Add(a, b int) int {
	return a + b
}
func TestAdd(t *testing.T) {
	result := Add(10, 20)
	expected := 30
	if result != expected {
		t.Errorf("Add(10, 20) = %d; want %d", result, expected)
	}
}

//go test -benchmem -run=^$ -bench .
// Benchmark function
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(10, 20)
	}
}
//what is b.N in the benchmark function?In the benchmark function, `b.N` is a variable provided by the testing framework that indicates the number of iterations to run the benchmark. The testing framework automatically determines how many times to execute the benchmarked code (in this case, the `Add` function) in order to get a reliable measurement of its performance. The loop runs `b.N` times, allowing the benchmark to gather enough data to calculate metrics such as time taken per operation and memory allocations.	
// To run the benchmark, you can use the command:
// go test -bench=. -benchmem	
// This command will execute the benchmark and provide output that includes the time taken per operation and memory allocations, allowing you to evaluate the performance of the `Add` function.
// Output example:
// BenchmarkAdd-8   	1000000000	         0.000 ns/op	       0 B/op	       0 allocs/op	
// In this output, `BenchmarkAdd-8` indicates the name of the benchmark and the number of CPU cores used. `1000000000` is the number of iterations (value of `b.N`), `0.000 ns/op` is the average time taken per operation, `0 B/op` is the average number of bytes allocated per operation, and `0 allocs/op` is the average number of memory allocations per operation.
// In summary, `b.N` is a crucial part of the benchmark function that allows the testing framework to determine how many times to execute the code being benchmarked in order to provide accurate performance metrics.
// In addition to measuring the time taken per operation, the `-benchmem` flag also provides insights into memory usage. It shows how many bytes are allocated per operation and how many memory allocations occur. This information can help you identify potential performance bottlenecks related to memory usage in your code.
// In the example output, `0 B/op` indicates that no bytes were allocated per operation, and `0 allocs/op` indicates that there were no memory allocations. This suggests that the `Add` function is efficient in terms of memory usage, as it does not require any additional memory allocation to perform its task.
//	
// Overall, using benchmarks in Go allows you to measure the performance of your code and identify areas for optimization. By analyzing the time taken per operation and memory usage, you can make informed decisions about how to improve the efficiency of your code.

//Explaination of the output:
	
/*Running tool: C:\Program Files\Go\bin\go.exe test -test.fullpath=true -benchmem -run=^$ -bench ^BenchmarkAdd$ bench-test

goos: windows
goarch: amd64
pkg: bench-test
cpu: 13th Gen Intel(R) Core(TM) i5-13420H
BenchmarkAdd-12    	1000000000	         0.2398 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	bench-test	1.025s
*/
// In this output:	
// - `BenchmarkAdd-12` indicates the name of the benchmark and the number of CPU cores used (12 in this case).
// - `1000000000` is the number of iterations (value of `b.N`), which means the benchmark ran the `Add` function one billion times to gather performance data.
// - `0.2398 ns/op` is the average time taken per operation, indicating that each call to the `Add` function took approximately 0.2398 nanoseconds on average.
// - `0 B/op` indicates that no bytes were allocated per operation, suggesting that the `Add` function does not require any additional memory allocation to perform its task.
// - `0 allocs/op` indicates that there were no memory allocations per operation, further confirming that the `Add` function is efficient in terms of memory usage.
// The output also includes information about the operating system (`goos: windows`), architecture (`goarch: amd64`), package name (`pkg: bench-test`), and CPU model (`cpu: 13th Gen Intel(R) Core(TM) i5-13420H`)
// Finally, the `PASS` indicates that the benchmark ran successfully, and `ok  	bench-test	1.025s` indicates that the tests and benchmarks for the `bench-test` package completed in approximately 1.025 seconds.
// In summary, the benchmark results show that the `Add` function is very efficient, with a very low time taken per operation and no memory allocations, making it a performant function for adding two integers.

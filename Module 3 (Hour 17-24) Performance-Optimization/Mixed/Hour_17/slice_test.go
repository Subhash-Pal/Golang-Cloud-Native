package main

import "testing"
// Bad: Keeps reallocating memory as it grows
func BenchmarkSliceAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []int
		for j := 0; j < 1000; j++ {
			s = append(s, j)
		}
	}
}
//go test -bench=BenchmarkSliceAppend -benchmem


// Good: Allocates memory once
func BenchmarkSlicePreAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, 1000)
		for j := 0; j < 1000; j++ {
			s = append(s, j)
		}
	}
}
//go test -bench=BenchmarkSlicePreAlloc -benchmem



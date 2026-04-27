package benchmark

import (
	"strconv"
	"testing"
)

func BenchmarkSumInts(b *testing.B) {
	values := make([]int, 10_000)
	for i := range values {
		values[i] = i
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_ = SumInts(values)
	}
}

func BenchmarkJoinWithPlus(b *testing.B) {
	parts := make([]string, 5_000)
	for i := range parts {
		parts[i] = strconv.Itoa(i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_ = JoinWithPlus(parts)
	}
}

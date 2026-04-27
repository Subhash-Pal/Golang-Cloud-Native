package main
import (
	"strings"
	"testing"
)

// Bad: Slow and creates many memory allocations
func BenchmarkStringPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := ""
		for j := 0; j < 100; j++ {
			s += "go"
		}
	}
}

// Good: Fast and uses very little memory
func BenchmarkStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		for j := 0; j < 100; j++ {
			sb.WriteString("go")
		}
		_ = sb.String()
	}
}

/*
1. String Concatenation (Memory Optimization)
This is the classic "Performance 101" in Go. It shows why strings.Builder is better than using + in a loop.
go
*/
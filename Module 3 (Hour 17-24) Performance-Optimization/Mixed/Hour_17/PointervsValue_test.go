
package main
import "testing"

type SmallStruct struct {
	A, B int
}

func BenchmarkPassValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		val := SmallStruct{A: 1, B: 2}
		_ = consumeValue(val)
	}
}

func BenchmarkPassPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		val := &SmallStruct{A: 1, B: 2}
		_ = consumePointer(val)
	}
}

func consumeValue(s SmallStruct) int { return s.A + s.B }
func consumePointer(s *SmallStruct) int { return s.A + s.B }



//go test -bench=BenchmarkPassValue -benchmem
//go test -bench=BenchmarkPassPointer -benchmem
//go test -bench=^BenchmarkPass* -benchmem
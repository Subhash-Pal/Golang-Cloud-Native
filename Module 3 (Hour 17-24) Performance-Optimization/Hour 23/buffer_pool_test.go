package pooldemo

import "testing"

func BenchmarkEncodeWithoutPool(b *testing.B) {
	orders := sampleOrders(1_000)
	b.ReportAllocs()

	for b.Loop() {
		_, err := EncodeWithoutPool(orders)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeWithPool(b *testing.B) {
	orders := sampleOrders(1_000)
	b.ReportAllocs()

	for b.Loop() {
		_, err := EncodeWithPool(orders)
		if err != nil {
			b.Fatal(err)
		}
	}
}

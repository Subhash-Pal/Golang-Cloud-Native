package main

import "testing"

func BenchmarkPipelineUnbuffered(b *testing.B) {
	for b.Loop() {
		_ = timePipeline(20_000, false)// Create unbuffered pipeline with 20,000 stages and time it
	}
}

func BenchmarkPipelineBuffered(b *testing.B) {
	for b.Loop() {
		_ = timePipeline(20_000, true)//Create buffered pipeline with 20,000 stages and time it
	}
}

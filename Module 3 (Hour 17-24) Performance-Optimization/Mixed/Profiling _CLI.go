package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func heavyComputation() {
	for i := 0; i < 1e8; i++ {
		_ = i * i
	}
}

func main() {
	// Create a file to store the CPU profile
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Error creating profile file:", err)
		return
	}
	defer f.Close()

	// Start CPU profiling
	fmt.Println("Starting CPU profiling...")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Simulate some heavy computation
	fmt.Println("Running heavy computation...")
	heavyComputation()

	// Simulate a delay
	time.Sleep(2 * time.Second)

	fmt.Println("Profiling complete. Check cpu.prof.")
}

/*
How to Use:
Run the program.
After execution, analyze the cpu.prof file using:
bash

go tool pprof cpu.prof
This will open an interactive terminal where you can explore the profiling data.
You can use commands like 'top', 'list', and 'web' to visualize the results.
*/
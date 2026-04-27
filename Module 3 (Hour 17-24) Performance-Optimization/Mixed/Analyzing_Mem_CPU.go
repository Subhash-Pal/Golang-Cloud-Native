package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)
//convert this to goroutine count and memory usage analysis code
/*func demo_goroutine_and_memory_analysis() {
	for i := 0; i < 5; i++ {
		go func() {
			for {
				time.Sleep(100 * time.Millisecond)
			}
		}()
}
}
*/
func heavyComputation() {
	data := make([]int, 1e6)
	for i := 0; i < len(data); i++ {
		data[i] = i * i
	}
}

func main() {
	// CPU profiling
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Error creating CPU profile file:", err)
		return
	}
	defer cpuFile.Close()
	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	// Memory profiling
	memFile, err := os.Create("mem.prof")
	if err != nil {
		fmt.Println("Error creating memory profile file:", err)
		return
	}
	defer memFile.Close()

	// Simulate heavy computation
	fmt.Println("Running heavy computation...")
	for i := 0; i < 5; i++ {
		//o demo_goroutine_and_memory_analysis()
		heavyComputation()
		time.Sleep(500 * time.Millisecond)
	}

	// Write memory profile
	if err := pprof.WriteHeapProfile(memFile); err != nil {
		fmt.Println("Error writing heap profile:", err)
	}

	fmt.Println("Profiling complete. Check cpu.prof and mem.prof.")
}
//How to run this code:
//1. Save the code in a file named main.go.
//2. Open a terminal and navigate to the directory containing main.go.
//3. Run the command: go run main.go	
//4. After the program finishes, you will find two files: cpu.prof and mem.prof in the same directory.
//5. You can analyze the CPU profile using the command: go tool pprof cpu.prof
//6. You can analyze the memory profile using the command: go tool pprof mem.prof	
//7. In the pprof interactive shell, you can use commands like "top", "list", and "web" to explore the profiles and identify performance bottlenecks or memory usage patterns.
//	Note: Make sure you have the Go toolchain installed and properly set up on your system to run this code and analyze the profiles.
package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // Import pprof for profiling
	"time"
)

func heavyComputation() {
	for i := 0; i < 1e7; i++ {
		_ = i * i
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Starting heavy computation...\n")
	heavyComputation()
	fmt.Fprintf(w, "Heavy computation completed.\n")
}

func main() {
	// Start a goroutine to simulate background work
	go func() {
		for {
			heavyComputation()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Register HTTP handlers
	http.HandleFunc("/", handler)

	// Start the HTTP server with pprof endpoints
	fmt.Println("Starting server on :8080...")
	fmt.Println("Access pprof at http://localhost:8080/debug/pprof/")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

/*
How to Use:
Run the program.
Access http://localhost:8080/debug/pprof/ in your browser.
Use go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30 to collect a 30-second CPU profile.



You can modify how top displays information by passing specific flags: 
top -cum: Sorts the list by cumulative time instead of flat time, helping you find high-level functions that drive overall usage.
top 10: Displays the top 10 nodes instead of the default 10 (or 5 in your case).

(pprof) list heavyComputation


---
(pprof) top -cum
Showing nodes accounting for 100% of total CPU time
	  flat  flat%   sum%        cum   cum%
	0.50s 50.00% 50.00%     0.50s 50.00%  main.heavyComputation
	0.30s 30.00% 80.00%     0.30s 30.00%  main.handler
	0.20s 20.00% 100.00%    0.20s 20.00%  runtime.main




	To help you optimize your main.heavyComputation function, here are detailed examples of how to use pprof commands and options in interactive mode.
1. Identify the Bottleneck Line with list 
The list command is your most powerful tool for seeing exactly which line of code is slow. 

Example Command: (pprof) list heavyComputation
Example Output:
text
ROUTINE ======================== main.heavyComputation
   10ms      1.84s (flat, cum) 95.83% of Total
      .          .     15: func heavyComputation() {
      .          .     16:     for i := 0; i < 1000000; i++ {
   10ms      1.83s     17:         result := doExpensiveTask(i) // Bottleneck found here!
      .          .     18:         fmt.Println(result)
      .          .     19:     }
      .          .     20: }
Use code with caution.

This shows that almost all the time (1.83s) is spent on line 17. 

2. Filter Noise with focus and ignore
If your profile is cluttered with system background tasks, use regex filters to clean it up. 

focus Example: (pprof) focus=main
Result: Only shows call stacks that include your code in the main package, discarding everything else.
ignore Example: (pprof) ignore=runtime|netpoll
Result: Hides Go's internal scheduler and networking overhead, allowing you to focus on application logic. 


3. Trace Callers with peek
If heavyComputation is called from many places and you want to know which specific caller is responsible for the load, use peek. 

Example Command: (pprof) peek heavyComputation
Example Output:
text
Context: main.heavyComputation
----------------------------------------------------------
      0      1.84s  main.main (caller)
----------------------------------------------------------
 1.83s      1.84s  main.heavyComputation
      0      1.83s  main.doExpensiveTask (callee)
----------------------------------------------------------
Use code with caution.

This tells you that main.main is the direct caller and doExpensiveTask is where the work is actually happening. 

4. Visual Analysis with web
If you have Graphviz installed, the web command provides a high-level visual graph. 

Example Command: (pprof) web
Result: It opens your default browser with a flowchart where:
Large boxes represent functions with high "flat" time (the function doing the actual work).
Thick arrows represent the most frequently executed code paths. 


5. Check Low-Level Detail with disasm
For extreme optimization, disasm shows the raw assembly instructions and how much time each instruction took. 

Example Command: (pprof) disasm heavyComputation
Why: Useful if you suspect the compiler isn't optimizing a loop properly or if there is unexpected memory allocation. 
Would you like to see how to compare two different profiles to see if your latest code changes actually improved performance?






*/
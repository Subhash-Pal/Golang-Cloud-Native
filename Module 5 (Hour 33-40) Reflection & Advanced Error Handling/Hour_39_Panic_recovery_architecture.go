/*

In Go, panic recovery is a critical part of building robust applications. When a panic occurs, it causes the program to terminate unless it is explicitly recovered using recover. A well-designed panic recovery architecture ensures that your application can gracefully handle unexpected errors without crashing.
Below is an example of how to implement a panic recovery architecture in Go.

Code: Panic Recovery Architecture
go

*/

package main

import (
	"fmt"
	"log"
	"runtime/debug"
)

// SafeRun wraps a function call with panic recovery.
func SafeRun(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			// Log the panic details
			log.Printf("Recovered from panic: %v\n", r)
			log.Printf("Stack Trace:\n%s", string(debug.Stack()))
		}
	}()

	// Execute the function
	fn()
}

// Example function that may panic.
func riskyOperation() {
	fmt.Println("Starting risky operation...")
	panic("Something went wrong!")
	fmt.Println("This will not be executed.")
}

// Another example function that calls a nested function which panics.
func nestedOperation() {
	fmt.Println("Starting nested operation...")
	SafeRun(func() {
		fmt.Println("Inside nested operation...")
		panic("Nested panic occurred!")
	})
	fmt.Println("Nested operation completed.")
}

func main() {
	// Example 1: Recovering from a direct panic
	fmt.Println("Example 1: Direct Panic Recovery")
	SafeRun(riskyOperation)
	fmt.Println("Continuing after panic recovery...\n")

	// Example 2: Recovering from a nested panic
	fmt.Println("Example 2: Nested Panic Recovery")
	nestedOperation()
	fmt.Println("Continuing after nested panic recovery...")
}


/*
Explanation of the Code
1. Panic Recovery with SafeRun
The SafeRun function wraps any function call (fn) with a defer block that uses recover to catch panics.
If a panic occurs:
The recover function captures the panic value.
The panic details and stack trace are logged using log.Printf.
This ensures that the program does not crash and continues execution.
2. Logging Stack Trace
The debug.Stack() function retrieves the current stack trace as a byte slice.
This is useful for debugging the source of the panic.
3. Risky Operations
The riskyOperation function simulates a panic by calling panic("Something went wrong!").
The nestedOperation function demonstrates recovering from a panic in a nested function.
4. Main Function
The main function demonstrates two scenarios:
Recovering from a direct panic in riskyOperation.
Recovering from a nested panic in nestedOperation.
Output Example
When you run the above code, the output will look something like this:

Example 1: Direct Panic Recovery
Starting risky operation...
2026/04/11 12:00:00 Recovered from panic: Something went wrong!
2026/04/11 12:00:00 Stack Trace:
goroutine 1 [running]:
runtime/debug.Stack(0xc0000a8008, 0x4d5f20, 0x52b9e0)
        /usr/local/go/src/runtime/debug/stack.go:24 +0x65
main.SafeRun.func1(0xc0000a8008)
        /path/to/main.go:14 +0x7f
panic(0x4d5f20, 0x52b9e0)
        /usr/local/go/src/runtime/panic.go:965 +0x1b9
main.riskyOperation()
        /path/to/main.go:28 +0x95
main.SafeRun(0x4efc68)
        /path/to/main.go:17 +0x4e
main.main()
        /path/to/main.go:41 +0x25
Continuing after panic recovery...

Example 2: Nested Panic Recovery
Starting nested operation...
Inside nested operation...
2026/04/11 12:00:00 Recovered from panic: Nested panic occurred!
2026/04/11 12:00:00 Stack Trace:
goroutine 1 [running]:
runtime/debug.Stack(0xc0000a8008, 0x4d5f20, 0x52b9e0)
        /usr/local/go/src/runtime/debug/stack.go:24 +0x65
main.SafeRun.func1(0xc0000a8008)
        /path/to/main.go:14 +0x7f
panic(0x4d5f20, 0x52b9e0)
        /usr/local/go/src/runtime/panic.go:965 +0x1b9
main.nestedOperation.func1()
        /path/to/main.go:35 +0x95
main.SafeRun(0xc0000a8008)
        /path/to/main.go:17 +0x4e
main.nestedOperation()
        /path/to/main.go:34 +0x65
main.main()
        /path/to/main.go:48 +0x5a
Nested operation completed.
Continuing after nested panic recovery...




Key Concepts Demonstrated
Panic Recovery:
The recover function allows you to catch and handle panics, preventing the program from terminating.
Graceful Degradation:
By recovering from panics, the program can continue executing other tasks, ensuring minimal disruption.
Stack Trace Logging:
Capturing and logging the stack trace helps identify the source of the panic during debugging.
Use Cases:
Web Servers: Recover from panics in HTTP handlers to prevent the server from crashing.
Background Workers: Handle unexpected errors in goroutines without affecting the main program.
Critical Systems: Ensure the system remains operational even in the face of unexpected failures.
Advanced Usage: Custom Recovery Handlers
You can extend the SafeRun function to support custom recovery handlers:
go

// SafeRunWithHandler wraps a function call with panic recovery and a custom handler.
func SafeRunWithHandler(fn func(), handler func(interface{})) {
	defer func() {
		if r := recover(); r != nil {
			// Call the custom handler if provided
			if handler != nil {
				handler(r)
			} else {
				// Default behavior: log the panic
				log.Printf("Recovered from panic: %v\n", r)
				log.Printf("Stack Trace:\n%s", string(debug.Stack()))
			}
		}
	}()

	// Execute the function
	fn()
}

func main() {
	// Example: Using a custom recovery handler
	SafeRunWithHandler(func() {
		panic("Custom panic occurred!")
	}, func(r interface{}) {
		fmt.Printf("Custom Handler: Recovered from panic: %v\n", r)
	})
}
	This will output:

	Notes
Best Practices:
Use panic recovery sparingly; prefer explicit error handling where possible.
Log detailed information (e.g., stack trace) to aid debugging.
Limitations:
recover only works within the same goroutine where the panic occurred.
Avoid using panics for normal control flow; they should be reserved for exceptional cases.
Extensibility:
You can integrate panic recovery with monitoring tools (e.g., Sentry, Prometheus) to track and alert on panics.



Key Observations from the Output
Recovered Panics:
The program successfully recovered from two panics:
"Something went wrong!" in the riskyOperation function.
"Nested panic occurred!" in the nestedOperation function.
These panics were caught by the recover function inside the SafeRun wrapper, preventing the program from crashing.
Stack Trace:
The stack trace provides detailed information about where the panic occurred.
For example:



This shows that the panic originated in the riskyOperation function at line 36 of your file.
It also shows the call stack leading up to the panic, including the SafeRun wrapper.
Graceful Continuation:
After recovering from each panic, the program continued executing subsequent code.
For example:
1



This demonstrates that the panic did not terminate the program.
Exit Code:
The program exited with code 0, indicating successful execution despite the panics.
This is expected because the panics were recovered, and no unhandled errors remained.
Detailed Breakdown of the Output
1. First Panic Recovery

2026/04/11 13:58:25 Recovered from panic: Something went wrong!
2026/04/11 13:58:25 Stack Trace:
goroutine 1 [running]:
runtime/debug.Stack()
    C:/Program Files/Go/src/runtime/debug/stack.go:26 +0x5e
main.SafeRun.func1()
    c:/Users/Shubh/OneDrive/Desktop/training_golang/Module 5 (Hour 33-40) Reflection & Advanced Error Handling/Hour_39_Panic_recovery_architecture.go:25 +0x8d
panic({0x7ff72fbca000?, 0x7ff72fbf3120?})
    C:/Program Files/Go/src/runtime/panic.go:860 +0x13a
main.riskyOperation()
    c:/Users/Shubh/OneDrive/Desktop/training_golang/Module 5 (Hour 33-40) Reflection & Advanced Error Handling/Hour_39_Panic_recovery_architecture.go:36 +0x59
main.SafeRun(0x7ff72fbf37c8?)
    c:/Users/Shubh/OneDrive/Desktop/training_golang/Module 5 (Hour 33-40) Reflection & Advanced Error Handling/Hour_39_Panic_recovery_architecture.go:30 +0x33
main.main()
    c:/Users/Shubh/OneDrive/Desktop/training_golang/Module 5 (Hour 33-40) Reflection & Advanced Error Handling/Hour_39_Panic_recovery_architecture.go:53 +0x56


Explanation:
The panic occurred in riskyOperation at line 36.
The SafeRun wrapper caught the panic and logged the details.
The stack trace shows the sequence of function calls leading to the panic.
2. Second Panic Recovery
2026/04/11 13:58:25 Recovered from panic: Nested panic occurred!
2026/04/11 13:58:25 Stack Trace:
goroutine 1 [running]:
runtime/debug.Stack()
    C:/Program Files/Go/src/runtime/debug/stack.go:26 +0x5e
main.SafeRun.func1()
    c:/Users/Shubh/OneDrive/Desktop/training_golang/Module 5 (Hour 33-40) Reflection & Advanced Error Handling/Hour_39_Panic_recovery_architecture.go:25 +0x8d
panic({0x7ff72fbca000?, 0x7ff72fbf31a0?})
    C:/Program Files/Go/src/runtime/panic.go:860 +0x13a
main.nestedOperation.func1()
    c:/Users/Shubh/OneDrive/Desktop/training_golang/Module 5 (Hour 33-40) Reflection & Advanced Error Handling/Hour_39_Panic_recovery_architecture.go:45 +0x59
main.SafeRun(0x7ff72fbf37c8?)
    c:/Users/Shubh/OneDrive/Desktop/training_golang/Module 5 (Hour 33-40) Reflection & Advanced Error Handling/Hour_39_Panic_recovery_architecture.go:30 +0x33
main.nestedOperation()
    c:/Users/Shubh/OneDrive/Desktop/training_golang/Module 5 (Hour 33-40) Reflection & Advanced Error Handling/Hour_39_Panic_recovery_architecture.go:43 +0x56
main.main()
    c:/Users/Shubh/OneDrive/Desktop/training_golang/Module 5 (Hour 33-40) Reflection & Advanced Error Handling/Hour_39_Panic_recovery_architecture.go:58 +0xcb


Explanation:
The panic occurred in the anonymous function inside nestedOperation at line 45.
The SafeRun wrapper caught the panic and logged the details.
The stack trace shows the nested structure of the function calls.
3. Program Completion

[Done] exited with code=0 in 1.824 seconds



Explanation:
The program completed successfully after recovering from both panics.
The exit code 0 indicates no unhandled errors.

Key Takeaways
Panic Recovery Works:
The SafeRun wrapper effectively caught and logged both panics, allowing the program to continue execution.
Stack Trace is Useful:
The stack trace provides critical debugging information, showing exactly where the panic occurred and the sequence of function calls leading to it.
Graceful Degradation:
By recovering from panics, the program avoided crashing and continued executing other tasks.
Best Practices:
Always log detailed information (e.g., stack trace) when recovering from a panic to aid debugging.
Use panic recovery sparingly; prefer explicit error handling for normal control flow.
Next Steps
If you'd like to enhance this architecture further, here are some ideas:
Integrate with Monitoring Tools:
Send panic details (e.g., stack trace) to monitoring tools like Sentry, Prometheus, or Datadog for centralized error tracking.
Custom Recovery Handlers:
Allow users to define custom recovery logic for different parts of the application.
Goroutine-Specific Recovery:
Extend the recovery mechanism to handle panics in goroutines, ensuring they don't crash the entire program.
Error Classification:
Classify panics into categories (e.g., critical vs. non-critical) and handle them differently based on their severity.
Let me know if you'd like to explore any of these enhancements or have further questions!

*/
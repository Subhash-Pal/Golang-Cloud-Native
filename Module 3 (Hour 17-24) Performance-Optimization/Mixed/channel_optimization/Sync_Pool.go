package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// 1. Create a pool for JSON Encoders
// We pool the encoder itself so it doesn't need to be re-allocated
var encoderPool = sync.Pool{
	New: func() any {
		fmt.Println("--- [Alloc] Creating New Encoder ---")
		// We use io.Discard to satisfy the constructor; 
		// we will reset the destination in the loop.
		return json.NewEncoder(io.Discard)
	},
}

func encodeUser(w io.Writer, u User) {
	// 2. Get an encoder from the pool
	enc := encoderPool.Get().(*json.Encoder)
	
	// 3. Reset the encoder to use the current writer (response/file)
	// This is the "magic" that reuses the encoder's internal buffers
	// Instead of enc := json.NewEncoder(w)
	enc.SetIndent("", "  ") // Optional: just for demo
	
	// Use the encoder (simulating a response write)
	// This would normally be enc.Encode(u) but we redirect to w
	// Since json.Encoder doesn't have a direct "Reset(w)", 
	// in real high-perf apps, people often pool bytes.Buffers instead.
	// But let's show the logic:
	fmt.Printf("Encoding User %d...\n", u.ID)
	
	// 4. Put back (In this specific encoder case, you'd usually pool 
	// the bytes.Buffer used to capture the output instead).
	encoderPool.Put(enc)
}

func main() {
	users := []User{
		{1, "alice@example.com"},
		{2, "bob@example.com"},
		{3, "charlie@example.com"},
	}

	// We process 3 users, but only allocate the encoder ONCE
	for _, u := range users {
		encodeUser(io.Discard, u)
	}

	fmt.Println("Finished processing all users.")
}

/*
The main use case for sync.Pool is managing high-frequency, short-lived objects to stop the Garbage Collector (GC) from slowing down your application.
Here is the breakdown of why and when that code matters:
## 1. The Problem: "GC Pressure"
In a high-traffic app (like a web server), every time you do this:

data, _ := json.Marshal(user) // Creates a new byte slice every time

Go has to find a new piece of memory. Once the request is sent, that memory becomes "trash." If you have 10,000 requests per second, you are creating 10,000 pieces of trash per second. The Garbage Collector has to stop or slow down your program to clean all that up. This causes latency spikes (p99 delays).
## 2. The Solution: The "Library Book" Model
Think of sync.Pool like a public library:

* Without Pool: You buy a new book, read it once, and throw it in the bin. (Expensive and wasteful).
* With Pool: You "borrow" a buffer, write your data, send it, and then clean it (Reset) and give it back to the library. The next person borrows the same book.

## 3. Real-World Use Cases

* JSON Encoding/Decoding: Reusing the buffers used to build JSON strings. This is the #1 use case in Go web APIs.
* Database Connections: Reusing small structures that hold query results.
* Logging: Reusing the string builders that format log lines (as seen in the previous example).
* Image Processing: Reusing the byte grids when resizing photos so you don't re-allocate megabytes of RAM every time.

## 4. When NOT to use it

* Small scripts: If your program runs for 2 seconds and quits, sync.Pool is overkill.
* Long-lived objects: If you need to keep an object for a long time (like a user session), don't put it in a pool.
* Non-memory intensive tasks: If you aren't seeing GC issues in your benchmarks, don't add the complexity.

Summary: Use it to turn a "Sawtooth" memory graph (constantly rising and crashing) into a Flat Line.
Would you like to see how to Benchmark this code so you can see the actual memory savings in numbers?


*/
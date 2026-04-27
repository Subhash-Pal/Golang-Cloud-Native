# 🧪 Golang Goroutine Leak vs Fix (No Restart Required)

## 🎯 Objective

Demonstrate a **controlled goroutine leak** and a **runtime fix** in a single Go application **without restarting the server**, with clear observable behavior.

---

## 🧠 Concept Overview

This project shows:

| Endpoint | Behavior                                                   |
| -------- | ---------------------------------------------------------- |
| `/leak`  | Creates goroutines that wait indefinitely (simulated leak) |
| `/fix`   | Cancels all leaked goroutines (cleanup)                    |
| `/stats` | Displays runtime metrics (goroutines, memory, GC)          |

---

## 🏗️ Architecture

```
Client → HTTP Server → Goroutines (tracked via context)
                              ↓
                        Cancel Functions
                              ↓
                           /fix endpoint
```

---

## 📦 Code (`main.go`)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	_ "net/http/pprof"
)

var (
	mu      sync.Mutex
	cancels []context.CancelFunc
)

func leakHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())

	mu.Lock()
	cancels = append(cancels, cancel)
	mu.Unlock()

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		}
	}(ctx)

	fmt.Fprintln(w, "leak created")
}

func fixHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	for _, cancel := range cancels {
		cancel()
	}
	cancels = nil
	mu.Unlock()

	time.Sleep(200 * time.Millisecond)

	fmt.Fprintln(w, "all leaks cleaned")
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Fprintf(w,
		"Goroutines: %d\nHeapAlloc: %d bytes\nHeapObjects: %d\nNumGC: %d\n",
		runtime.NumGoroutine(),
		m.HeapAlloc,
		m.HeapObjects,
		m.NumGC,
	)
}

func main() {
	http.HandleFunc("/leak", leakHandler)
	http.HandleFunc("/fix", fixHandler)
	http.HandleFunc("/stats", statsHandler)

	log.Println("Server running at :8080")
	log.Println("Endpoints: /leak | /fix | /stats")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## ⚙️ Setup & Run

### Step 1 — Initialize project

```bash
go mod init leak-demo
```

### Step 2 — Run application

```bash
go run main.go
```

---

## 🚀 Execution Steps

### 🔹 Step 1 — Check baseline

```powershell
curl http://localhost:8080/stats
```

Expected:

```
Goroutines: 2–4
```

---

### 🔹 Step 2 — Generate leak

```powershell
1..500 | ForEach-Object { curl.exe -s http://localhost:8080/leak > $null }
```

---

### 🔹 Step 3 — Observe leak

```powershell
curl http://localhost:8080/stats
```

Expected:

```
Goroutines: 500+
```

👉 This value **will not decrease automatically**

---

### 🔹 Step 4 — Fix leak (no restart)

```powershell
curl http://localhost:8080/fix
```

---

### 🔹 Step 5 — Verify recovery

```powershell
curl http://localhost:8080/stats
```

Expected:

```
Goroutines: back to ~2–10
```

---

## 🔍 Profiling (Optional)

Open in browser:

```
http://localhost:8080/debug/pprof/goroutine?debug=2
```

### Before `/fix`

* Thousands of blocked goroutines

### After `/fix`

* Only minimal system goroutines remain

---

## 📊 Expected Behavior Summary

| Phase         | Goroutines        |
| ------------- | ----------------- |
| Start         | ~2                |
| After `/leak` | High (e.g., 500+) |
| After `/fix`  | Back to baseline  |

---

## ⚠️ Key Learning Points

* Goroutine leaks occur when **no termination condition exists**
* Go runtime **cannot kill goroutines externally**
* Proper design requires:

  * `context.WithCancel`
  * lifecycle tracking
  * explicit cleanup

---

## ✅ Why This Demo Works

* No restart required
* Deterministic behavior
* Clear rise → drop pattern
* Realistic concurrency control pattern

---

## 🔧 Extensions

You can extend this demo to:

* HTTP client leaks
* Database connection leaks
* Channel blocking patterns
* Worker pool mismanagement

---

## 📌 Conclusion

This project demonstrates a **practical and observable approach** to:

* Identify goroutine leaks
* Control their lifecycle
* Implement deterministic cleanup

This pattern is directly applicable to **production debugging and system design**.

---

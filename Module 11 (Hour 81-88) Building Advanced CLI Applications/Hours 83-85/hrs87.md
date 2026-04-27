Good—this is where many “it works on my machine” setups break. We’ll keep it **practical and verifiable**.

---

# 🔹 Hour 87 – Cross Compilation & Packaging (Go CLI)

We’ll take a simple CLI and:

1. Build for **multiple OS/ARCH**
2. Verify binaries
3. Package for distribution

---

# 🎯 Objective

From one machine, generate binaries for:

* Windows (`.exe`)
* Linux (server deploy)
* macOS (optional)

---

# 📂 Minimal Project (Standalone)

```text
hour87_build/
│── main.go
│── go.mod
```

---

# 📄 `main.go`

```go
package main

import "fmt"

func main() {
	fmt.Println("CLI Running Successfully")
}
```

---

# 🧪 Step 1 — Initialize

```bash
go mod init hour87_build
go mod tidy
```

---

# 🧪 Step 2 — Native Build (baseline)

```bash
go build -o app
```

### ✔️ Verify

```bash
./app
```

Output:

```
CLI Running Successfully
```

---

# 🧪 Step 3 — Cross Compile

## 🔸 Windows (from Linux/Mac or same machine)

```bash
GOOS=windows GOARCH=amd64 go build -o app.exe
```

---

## 🔸 Linux (most important for servers)

```bash
GOOS=linux GOARCH=amd64 go build -o app-linux
```

---

## 🔸 macOS

```bash
GOOS=darwin GOARCH=amd64 go build -o app-mac
```

---

# 🧪 Step 4 — Verify Binary Type

## On Linux/Mac:

```bash
file app-linux
```

Expected:

```
ELF 64-bit executable
```

## On Windows (PowerShell):

```powershell
Get-Item app.exe
```

---

# 🧠 Key Concepts

| Variable | Meaning          |
| -------- | ---------------- |
| `GOOS`   | Target OS        |
| `GOARCH` | CPU architecture |
| `amd64`  | 64-bit systems   |

---

# 🔧 Step 5 — Static Binary (important for deployment)

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-linux
```

### Why?

* No external C dependencies
* Runs on minimal servers (Docker, Alpine)

---

# 📦 Step 6 — Packaging

## 🔸 Linux (.tar.gz)

```bash
tar -czvf app-linux.tar.gz app-linux
```

---

## 🔸 Windows (.zip)

```powershell
Compress-Archive -Path app.exe -DestinationPath app-windows.zip
```

---

# 📂 Final Distribution Example

```text
dist/
│── app-linux.tar.gz
│── app-windows.zip
│── app-mac.tar.gz
```

---

# ⚠️ Common Failures (Real Debugging)

---

## ❌ 1. “exec format error”

### Cause:

Wrong OS binary

### Example:

Running Linux binary on Windows

### Fix:

Match OS:

```bash
GOOS=windows
```

---

## ❌ 2. “permission denied” (Linux)

### Cause:

Executable bit missing

### Fix:

```bash
chmod +x app-linux
```

---

## ❌ 3. Binary runs locally but fails on server

### Cause:

CGO dependency

### Fix:

```bash
CGO_ENABLED=0
```

---

## ❌ 4. Wrong architecture

### Example:

```bash
GOARCH=arm64
```

Running on:

```
amd64 server
```

### Fix:

```bash
GOARCH=amd64
```

---

# 🚀 (Optional but Industry-Grade)

Use GoReleaser to automate:

* Multi-OS builds
* Packaging
* GitHub releases

---

# 🧪 Your Task (Do This Now)

Run:

```bash
GOOS=linux GOARCH=amd64 go build -o app-linux
```

Then tell me:

👉 What happens when you try to run `app-linux` on your current system?

---

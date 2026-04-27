# Hour 47 - CGO Integration in Go

This project demonstrates how to use CGO to call C functions from a Go program.

## What this example does

The program defines two C functions inside the Go file:

- `addNumbers`
- `multiplyNumbers`

Go then calls these C functions and prints the results.

## Files

- `main.go`
  Contains the embedded C code and the Go code that calls it through CGO.

## Commands to run in order

Open PowerShell and go to the folder:

```powershell
cd "D:\training_golang\Module 6 (Hour 41-48) Database & External Integration\Hour 47"
```

Set the CGO compiler environment for this PowerShell session:

```powershell
Remove-Item Env:CGO_CFLAGS -ErrorAction SilentlyContinue
Remove-Item Env:CGO_LDFLAGS -ErrorAction SilentlyContinue
$env:PATH = "C:\msys64\mingw64\bin;" + $env:PATH
$env:CC = "x86_64-w64-mingw32-gcc"
```

Format the Go file:

```powershell
gofmt -w main.go
```

Build the program:

```powershell
go build -o cgo-demo.exe .
```

Run the program:

```powershell
.\cgo-demo.exe
```

## Expected output

```text
CGO Integration Example
Addition using C function: 10 + 5 = 15
Multiplication using C function: 10 * 5 = 50
```

## Requirements

- Go installed
- CGO enabled
- C compiler available, such as `gcc`
- MinGW compiler available at `C:\msys64\mingw64\bin`

## Check your environment

```powershell
go env CGO_ENABLED CC
```

## Working build notes for this machine

The build worked successfully on this machine after:

- clearing `CGO_CFLAGS`
- clearing `CGO_LDFLAGS`
- adding `C:\msys64\mingw64\bin` to `PATH`
- setting `CC=x86_64-w64-mingw32-gcc`

After that, these commands worked:

```powershell
go build -o cgo-demo.exe .
.\cgo-demo.exe
```

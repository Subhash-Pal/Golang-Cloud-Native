# Hour 47 - CGO Integration

This project demonstrates how Go can call C code through CGO.

## Prerequisite

A C compiler must be installed and available in `PATH`.

Useful official/reference links:

- Go install: [https://go.dev/doc/install](https://go.dev/doc/install)
- Chocolatey install: [https://chocolatey.org/install](https://chocolatey.org/install)
- MSYS2 install: [https://www.msys2.org/](https://www.msys2.org/)
- TDM-GCC: [https://jmeubank.github.io/tdm-gcc/](https://jmeubank.github.io/tdm-gcc/)

Check:

```powershell
gcc --version
go env CGO_ENABLED
```

`go env CGO_ENABLED` should print `1`.

In PowerShell, prefer:

```powershell
where.exe gcc
```

instead of `where gcc`.

## Fresh Windows Setup

Use this setup if CGO is not working yet on your machine.

## One Script Setup For New Machine

If you want the easiest setup for Hour 47 on a fresh Windows machine, run this in PowerShell as Administrator:

```powershell
.\setup-hour47-new-machine.ps1
```

Optional if Go is already installed:

```powershell
.\setup-hour47-new-machine.ps1 -SkipGoInstall
```

What it does:

- installs Chocolatey if missing
- installs Go unless skipped
- installs one GCC toolchain with Chocolatey MinGW
- clears stale `CC`, `CXX`, and `CGO_*` variables
- verifies `go`, `gcc`, and `CGO_ENABLED`

After the setup script finishes:

1. Close the Administrator PowerShell window
2. Open a fresh normal PowerShell
3. Run:

```powershell
cd 'D:\training_golang\Module 6 (Hour 41-48) Database & External Integration\Hour 47'
.\run.ps1
```

### Fastest Admin Fix

Open PowerShell as Administrator and run:

```powershell
choco install mingw -y
```

Then close that window and open a new normal PowerShell.

Verify:

```powershell
gcc --version
where.exe gcc
```

Then run the project:

```powershell
cd 'D:\training_golang\Module 6 (Hour 41-48) Database & External Integration\Hour 47'
go run .\main.go
```

If Chocolatey install fails, use the MSYS2 repair command instead:

```powershell
C:\msys64\usr\bin\pacman.exe -S --noconfirm mingw-w64-x86_64-gcc mingw-w64-x86_64-binutils
```

This is the easiest option on a fresh Windows machine.

### Option 1: MSYS2 with MinGW-w64

Use this option if you prefer MSYS2.

1. Install MSYS2 from [https://www.msys2.org/](https://www.msys2.org/).
2. Open the `MSYS2 MINGW64` terminal.
3. Update packages:

```powershell
pacman -Syu
```

4. Reopen the MSYS2 terminal if it asks you to.
5. Install the compiler toolchain:

```powershell
pacman -S --needed mingw-w64-x86_64-gcc mingw-w64-x86_64-binutils
```

6. Add the compiler to your Windows `PATH`.

Typical path:

```text
C:\msys64\mingw64\bin
```

7. Close PowerShell and open a fresh PowerShell window.
8. Verify:

```powershell
gcc --version
where.exe gcc
go env CGO_ENABLED
```

9. Run the project from normal PowerShell:

```powershell
cd 'D:\training_golang\Module 6 (Hour 41-48) Database & External Integration\Hour 47'
.\run.ps1
```

### Option 2: MinGW-w64

1. Install MinGW-w64.
2. Add its `bin` folder to `PATH`.
3. Open a new PowerShell window.
4. Verify with:

```powershell
gcc --version
where.exe gcc
go env CGO_ENABLED
```

## Clean Go Environment

If CGO still fails, clear any saved custom Go compiler overrides:

```powershell
go env -u CC
go env -u CXX
go env -u CGO_CFLAGS
go env -u CGO_CPPFLAGS
go env -u CGO_CXXFLAGS
go env -u CGO_LDFLAGS
```

Then close PowerShell and open a new one.

## Run

```powershell
go mod tidy
.\run.ps1
```

The `run.ps1` script now clears stale `CC`, `CXX`, and `CGO_*` variables before running, so it is safer than calling `go run .\main.go` directly on a machine with older CGO settings.

This runner was verified successfully on this machine.

Or run directly:

```powershell
go run .\main.go
```

## Expected Output

When the setup is correct, the program should print:

```text
CGO Integration Example
Addition using C function: 10 + 5 = 15
Multiplication using C function: 10 * 5 = 50
```

## Verified Working Command

Verified command:

```powershell
.\run.ps1
```

Verified result:

- exits with code `0`
- prints the expected CGO output

## Requirement

Docker is not required for this example. GCC and CGO setup are required.

## Troubleshooting

### Error: `runtime/cgo: ... cgo.exe: exit status 2`

This means the local C compiler toolchain is not healthy yet.

On this machine, the investigation showed:

- `CGO_ENABLED=1`
- `gcc.exe` was found
- the internal GCC compiler step still failed

That means the Go code is fine, but the Windows GCC/MSYS2 setup still needs repair.

### Recommended recovery steps

1. Clear saved Go CGO overrides with the `go env -u ...` commands above.
2. Open a fresh PowerShell window.
3. Verify plain C compilation works first:

```powershell
gcc --version
```

4. If GCC still behaves incorrectly, reinstall or repair MSYS2 `mingw-w64-x86_64-gcc`.
5. After GCC is healthy, run:

```powershell
go run .\main.go
```

### Important note

If you previously set custom `TMP`, `TEMP`, `CC`, or `CGO_*` environment variables, they can also break CGO builds. Use a clean shell session after installation changes.

## Another System Setup

Use these minimal steps on another machine to avoid the same issue:

1. Install Go.
2. Install one GCC toolchain only.
3. Make sure `gcc` is on `PATH`.
4. Do not manually set `CC`, `CGO_CFLAGS`, or `CGO_LDFLAGS` unless required.
5. Open a fresh PowerShell window.
6. Verify:

```powershell
where.exe gcc
gcc --version
go env CGO_ENABLED
```

7. Run the project with:

```powershell
cd 'D:\training_golang\Module 6 (Hour 41-48) Database & External Integration\Hour 47'
.\run.ps1
```

## Short Version

On a clean machine:

1. Install Go
2. Install GCC
3. Confirm `where.exe gcc` works
4. Run `.\run.ps1`

## Suggested Choice

Use one of these:

1. Easiest: `choco install mingw -y`
2. If you prefer MSYS2:
   Install MSYS2
   Open `MSYS2 MINGW64`
   Run `pacman -Syu`
   Reopen terminal if asked
   Run `pacman -S --needed mingw-w64-x86_64-gcc mingw-w64-x86_64-binutils`
   Add `C:\msys64\mingw64\bin` to `PATH`
   Open a normal PowerShell
   Run `.\run.ps1`

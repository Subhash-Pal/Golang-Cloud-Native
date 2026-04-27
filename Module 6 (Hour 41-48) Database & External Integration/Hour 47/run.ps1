$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot

# Clean stale CGO/MSYS2 overrides that may be inherited by a new shell.
Remove-Item Env:CC -ErrorAction SilentlyContinue
Remove-Item Env:CXX -ErrorAction SilentlyContinue
Remove-Item Env:CGO_CFLAGS -ErrorAction SilentlyContinue
Remove-Item Env:CGO_CPPFLAGS -ErrorAction SilentlyContinue
Remove-Item Env:CGO_CXXFLAGS -ErrorAction SilentlyContinue
Remove-Item Env:CGO_LDFLAGS -ErrorAction SilentlyContinue

$env:CC = "gcc"
$env:CXX = "g++"

Write-Host "Requirement: GCC and CGO must be configured before this example will build." -ForegroundColor Yellow
Write-Host "Docker is not required for this example." -ForegroundColor Yellow
Write-Host "Runner note: this script clears stale CGO variables and uses gcc/g++ from PATH." -ForegroundColor Yellow
Write-Host "Running Hour 47 - CGO integration..."
go run .

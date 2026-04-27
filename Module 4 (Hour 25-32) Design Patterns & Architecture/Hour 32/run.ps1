$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot
Write-Host "Starting Hour 32 user API on http://127.0.0.1:8081"
go run .

$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot
Write-Host "Starting Hour 26 book API on http://127.0.0.1:8080"
go run .

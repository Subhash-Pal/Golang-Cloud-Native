Set-Location (Join-Path $PSScriptRoot "..")

Write-Host "Running: go run . bucket create" -ForegroundColor Cyan
go run . bucket create

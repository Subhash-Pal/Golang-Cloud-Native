Set-Location (Join-Path $PSScriptRoot "..")

Write-Host "Running: go run . list" -ForegroundColor Cyan
go run . list

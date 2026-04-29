Set-Location (Join-Path $PSScriptRoot "..")

Write-Host "Running: go run . put app.config.version v1" -ForegroundColor Cyan
go run . put app.config.version v1

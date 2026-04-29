Set-Location (Join-Path $PSScriptRoot "..")

Write-Host "Running: go run . get app.config.version" -ForegroundColor Cyan
go run . get app.config.version

Set-Location (Join-Path $PSScriptRoot "..")

Write-Host ""
Write-Host "Building distkv-cli..." -ForegroundColor Cyan
go build ./...

Write-Host ""
Write-Host "Build completed." -ForegroundColor Green

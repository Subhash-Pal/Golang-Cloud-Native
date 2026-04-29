Set-Location (Join-Path $PSScriptRoot "..")

Write-Host ""
Write-Host "Starting watcher for app.> keys..." -ForegroundColor Cyan
Write-Host "Open another console and run: go run . put app.config.version v3" -ForegroundColor Yellow
Write-Host ""

go run . watch "app.>"

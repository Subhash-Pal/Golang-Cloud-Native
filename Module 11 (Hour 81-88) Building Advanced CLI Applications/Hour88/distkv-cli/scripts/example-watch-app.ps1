Set-Location (Join-Path $PSScriptRoot "..")

Write-Host "Running: go run . watch ""app.>""" -ForegroundColor Cyan
Write-Host "Open another console and run a put command to see live updates." -ForegroundColor Yellow
go run . watch "app.>"

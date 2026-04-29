Set-Location (Join-Path $PSScriptRoot "..")

Write-Host ""
Write-Host "Starting NATS JetStream backend..." -ForegroundColor Cyan
docker compose up -d

Write-Host ""
Write-Host "Checking backend status..." -ForegroundColor Cyan
docker compose ps

param(
    [switch]$RemoveVolume
)

Set-Location (Join-Path $PSScriptRoot "..")

Write-Host ""
Write-Host "Stopping NATS JetStream backend..." -ForegroundColor Cyan

if ($RemoveVolume) {
    docker compose down -v
} else {
    docker compose down
}

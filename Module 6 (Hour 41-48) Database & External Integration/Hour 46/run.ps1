$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot

Write-Host "Requirement: No external service needed for this example." -ForegroundColor Yellow
Write-Host "Docker is not required for this example." -ForegroundColor Yellow
Write-Host "Running Hour 46 - Retry logic..."
go run .

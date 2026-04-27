$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot

Write-Host "Requirement: Redis must be running on 127.0.0.1:6379." -ForegroundColor Yellow
Write-Host "If you use Docker, start Docker Desktop and the Redis container first." -ForegroundColor Yellow
try {
    $tcpClient = New-Object System.Net.Sockets.TcpClient
    $tcpClient.Connect("127.0.0.1", 6379)
    $tcpClient.Close()
} catch {
    Write-Host "Redis is not running. Start it first, then run this script again." -ForegroundColor Red
    Write-Host "If you use Docker, make sure Docker Desktop is running." -ForegroundColor Yellow
    Write-Host "Example: docker compose up -d redis" -ForegroundColor Yellow
    exit 1
}

Write-Host "Running Hour 44 - Redis caching..."
go run .

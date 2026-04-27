$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot

Write-Host "Requirement: PostgreSQL must be running on 127.0.0.1:5432." -ForegroundColor Yellow
Write-Host "If you use Docker, start Docker Desktop and the PostgreSQL container first." -ForegroundColor Yellow
try {
    $tcpClient = New-Object System.Net.Sockets.TcpClient
    $tcpClient.Connect("127.0.0.1", 5432)
    $tcpClient.Close()
} catch {
    Write-Host "PostgreSQL is not running. Start it first, then run this script again." -ForegroundColor Red
    Write-Host "If you use Docker, make sure Docker Desktop is running." -ForegroundColor Yellow
    Write-Host "Example: docker compose up -d postgres" -ForegroundColor Yellow
    exit 1
}

Write-Host "Running Hour 42 - PostgreSQL CRUD..."
go run .

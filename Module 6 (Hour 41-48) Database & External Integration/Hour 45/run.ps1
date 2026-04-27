param(
    [ValidateSet("consumer", "producer")]
    [string]$Mode = "consumer",
    [string]$Message = ""
)

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

if ($Mode -eq "consumer") {
    Write-Host "Running Hour 45 consumer..."
    go run .\cmd\consumer
    exit $LASTEXITCODE
}

if ([string]::IsNullOrWhiteSpace($Message)) {
    Write-Host "Running Hour 45 producer with auto-generated message..."
    go run .\cmd\producer
    exit $LASTEXITCODE
}

Write-Host "Running Hour 45 producer with message: $Message"
go run .\cmd\producer $Message

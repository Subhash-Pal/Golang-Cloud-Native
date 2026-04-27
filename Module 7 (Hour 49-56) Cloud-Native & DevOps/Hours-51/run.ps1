param(
    [ValidateSet("compose")]
    [string]$Mode = "compose"
)

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

Write-Host "Starting Hour 51 Docker Compose stack on http://localhost:18051"
try {
    docker compose down --rmi local --remove-orphans 2>$null | Out-Null
    docker compose up -d --build
    if ($LASTEXITCODE -ne 0) { throw "Docker Compose startup failed." }
    Start-Sleep -Seconds 5
    $response = Invoke-RestMethod -Uri "http://localhost:18051/"
    $response | ConvertTo-Json -Depth 5
}
finally {
    Write-Host "Stopping stack and removing local images"
    docker compose down --rmi local --remove-orphans
}

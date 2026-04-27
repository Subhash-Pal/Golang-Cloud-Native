$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

$env:CLUSTER_NAME = "training-cluster"
$env:NODE_NAME = "worker-1"
$env:NAMESPACE = "module-7"
$env:POD_NAME = "hour52-architecture-demo"
$env:PORT = "18452"

Write-Host "Verifying Hour 52 locally on http://localhost:18452"
Write-Host "Endpoints: /  /components"
$job = Start-Job -ScriptBlock {
    param($scriptRoot)
    Set-Location $scriptRoot
    $env:CLUSTER_NAME = "training-cluster"
    $env:NODE_NAME = "worker-1"
    $env:NAMESPACE = "module-7"
    $env:POD_NAME = "hour52-architecture-demo"
    $env:PORT = "18452"
    go run .
} -ArgumentList $PSScriptRoot

try {
    Start-Sleep -Seconds 3
    $response = Invoke-RestMethod -Uri "http://localhost:18452/"
    $response | ConvertTo-Json -Depth 5
}
finally {
    Stop-Job $job -ErrorAction SilentlyContinue | Out-Null
    Remove-Job $job -ErrorAction SilentlyContinue | Out-Null
}

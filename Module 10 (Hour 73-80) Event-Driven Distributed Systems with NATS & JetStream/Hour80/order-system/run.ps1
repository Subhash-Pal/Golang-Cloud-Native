# ================================
# Distributed Order System Runner
# (go run mode — Defender-safe)
# ================================

$ErrorActionPreference = "Stop"

# Resolve script directory (not current terminal)
$projectPath = Split-Path -Parent $MyInvocation.MyCommand.Path

Write-Host "Project Path: $projectPath" -ForegroundColor Cyan

# -------------------------------
# Step 1 — Setup
# -------------------------------
Write-Host "Step 1: Running setup..." -ForegroundColor Yellow
go run "$projectPath\cmd\setup\main.go"

Start-Sleep -Seconds 1

# -------------------------------
# Step 2 — Start Processor
# -------------------------------
Write-Host "Step 2: Starting Processor..." -ForegroundColor Green

Start-Process powershell -ArgumentList @(
    "-NoExit",
    "-Command",
    "cd `"$projectPath`"; go run cmd/processor/main.go"
)

Start-Sleep -Seconds 2

# -------------------------------
# Step 3 — Start Publisher
# -------------------------------
Write-Host "Step 3: Starting Publisher..." -ForegroundColor Magenta

Start-Process powershell -ArgumentList @(
    "-NoExit",
    "-Command",
    "cd `"$projectPath`"; go run cmd/publisher/main.go"
)

Start-Sleep -Seconds 2

# -------------------------------
# Step 4 — Status Check
# -------------------------------
Write-Host "Step 4: Checking Status..." -ForegroundColor Blue
go run "$projectPath\cmd\status\main.go"

# -------------------------------
# Done
# -------------------------------
Write-Host ""
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "Execution Completed 🚀" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Cyan
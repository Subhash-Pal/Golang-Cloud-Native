param(
    [string]$ModuleRoot = $PSScriptRoot,
    [switch]$SkipDockerDesktop
)

$ErrorActionPreference = "Stop"

function Write-Step {
    param([string]$Message)
    Write-Host ""
    Write-Host "==> $Message" -ForegroundColor Cyan
}

function Test-IsAdmin {
    $identity = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($identity)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Ensure-Chocolatey {
    if (Get-Command choco -ErrorAction SilentlyContinue) {
        Write-Host "Chocolatey already installed." -ForegroundColor Green
        return
    }

    Write-Step "Installing Chocolatey"
    Set-ExecutionPolicy Bypass -Scope Process -Force
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
    Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
}

function Install-ChocoPackage {
    param([string]$PackageName)

    Write-Step "Installing $PackageName"
    choco install $PackageName -y --no-progress
}

function Clear-StaleCgoEnvironment {
    Write-Step "Clearing stale CGO compiler variables"

    foreach ($name in @('CC', 'CXX', 'CGO_CFLAGS', 'CGO_CPPFLAGS', 'CGO_CXXFLAGS', 'CGO_LDFLAGS')) {
        Remove-Item "Env:$name" -ErrorAction SilentlyContinue
        [Environment]::SetEnvironmentVariable($name, $null, 'User')
        [Environment]::SetEnvironmentVariable($name, $null, 'Machine')
    }

    if (Get-Command go -ErrorAction SilentlyContinue) {
        go env -u CC
        go env -u CXX
        go env -u CGO_CFLAGS
        go env -u CGO_CPPFLAGS
        go env -u CGO_CXXFLAGS
        go env -u CGO_LDFLAGS
    }
}

function Refresh-ProcessPathFromMachineAndUser {
    $machinePath = [Environment]::GetEnvironmentVariable('Path', 'Machine')
    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    if ([string]::IsNullOrWhiteSpace($userPath)) {
        $env:Path = $machinePath
    } else {
        $env:Path = "$machinePath;$userPath"
    }
}

if (-not (Test-IsAdmin)) {
    throw "Run this script in PowerShell as Administrator."
}

Write-Step "Preparing new Windows machine for Module 6"
Write-Host "Module root: $ModuleRoot"

Ensure-Chocolatey

Install-ChocoPackage -PackageName "golang"
Install-ChocoPackage -PackageName "mingw"

if (-not $SkipDockerDesktop) {
    Install-ChocoPackage -PackageName "docker-desktop"
}

Refresh-ProcessPathFromMachineAndUser
Clear-StaleCgoEnvironment

Write-Step "Verifying tools"
Write-Host "Go version:" -ForegroundColor Yellow
go version

Write-Host ""
Write-Host "GCC path:" -ForegroundColor Yellow
where.exe gcc

Write-Host ""
Write-Host "GCC version:" -ForegroundColor Yellow
gcc --version

Write-Host ""
Write-Host "CGO enabled:" -ForegroundColor Yellow
go env CGO_ENABLED

Write-Step "Next steps"
Write-Host "1. Close this window and open a fresh normal PowerShell." -ForegroundColor Green
if (-not $SkipDockerDesktop) {
    Write-Host "2. Start Docker Desktop and wait until it is fully running." -ForegroundColor Green
    Write-Host "3. Start PostgreSQL and Redis from the module root:" -ForegroundColor Green
    Write-Host "   docker compose up -d" -ForegroundColor White
    Write-Host "4. Run an hour folder with its helper script, for example:" -ForegroundColor Green
} else {
    Write-Host "2. Start PostgreSQL and Redis using your preferred local setup." -ForegroundColor Green
    Write-Host "3. Run an hour folder with its helper script, for example:" -ForegroundColor Green
}
Write-Host "   Set-Location '$ModuleRoot\Hour 47'" -ForegroundColor White
Write-Host "   .\run.ps1" -ForegroundColor White

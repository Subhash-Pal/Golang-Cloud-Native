Set-Location (Join-Path $PSScriptRoot "..")

function Run-Step {
    param(
        [string]$Title,
        [string]$Command
    )

    Write-Host ""
    Write-Host ("=" * 64) -ForegroundColor DarkGray
    Write-Host $Title -ForegroundColor Yellow
    Write-Host "Command: $Command" -ForegroundColor Green
    Write-Host ("=" * 64) -ForegroundColor DarkGray
    Invoke-Expression $Command
}

Run-Step -Title "Create or reuse bucket" -Command "go run . bucket create"
Run-Step -Title "Store initial version" -Command "go run . put app.config.version v1"
Run-Step -Title "Read stored value" -Command "go run . get app.config.version"
Run-Step -Title "List keys" -Command "go run . list"
Run-Step -Title "Health check" -Command "go run . health"
Run-Step -Title "Update version" -Command "go run . put app.config.version v2"
Run-Step -Title "Read latest version" -Command "go run . get app.config.version"

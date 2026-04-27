# Get the current folder path and wrap it in single quotes to handle spaces
$currentDir = $PSScriptRoot

# 1. Start NATS server (uses single quotes around $currentDir for safety)
Write-Host "Starting NATS Server..." -ForegroundColor Yellow

Start-Sleep -Seconds 2

# 2. Create KV bucket (run with local path prefix)
Write-Host "Creating KV bucket..." -ForegroundColor Cyan
go run .\kv_create.go

# 3. Start watcher (Terminal 1)
Write-Host "Starting Watcher..." -ForegroundColor Yellow
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$currentDir'; go run .\kv_watch.go"

Start-Sleep -Seconds 2

# 4. Put value (Terminal 2)
Write-Host "Putting value..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$currentDir'; go run .\kv_put.go"

Start-Sleep -Seconds 1

# 5. Read value
Write-Host "Reading value from KV:" -ForegroundColor Green
go run .\kv_get.go

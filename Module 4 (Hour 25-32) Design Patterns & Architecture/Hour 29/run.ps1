$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot
ls -Recurse
go run .

$graphvizBin = "D:\training_golang\tools\graphviz\graphviz-2.38\release\bin"

if (-not (Test-Path $graphvizBin)) {
	Write-Error "Graphviz bin folder not found at $graphvizBin"
	return
}

$pathEntries = $env:PATH -split ";"
if ($pathEntries -notcontains $graphvizBin) {
	$env:PATH = "$graphvizBin;$env:PATH"
}

Write-Host "Graphviz enabled for this PowerShell session."
Write-Host "You can now run: go tool pprof -http=:8081 cpu.prof"

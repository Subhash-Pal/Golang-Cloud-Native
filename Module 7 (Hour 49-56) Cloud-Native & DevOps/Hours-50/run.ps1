param(
    [ValidateSet("local", "docker")]
    [string]$Mode = "docker"
)

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

function Remove-DockerContainerIfExists {
    param([string]$Name)
    $containerNames = docker ps -a --format "{{.Names}}"
    if ($containerNames -contains $Name) { docker rm -f $Name | Out-Null }
}

function Remove-DockerImageIfExists {
    param([string]$Name)
    $imageNames = docker images --format "{{.Repository}}:{{.Tag}}"
    if (($imageNames -contains $Name) -or ($imageNames -contains "${Name}:latest")) { docker image rm -f $Name | Out-Null }
}

switch ($Mode) {
    "local" {
        $env:PORT = "18450"
        Write-Host "Running Hour 50 locally on http://localhost:18450"
        Write-Host "Endpoints: /  /build  /text"
        go run .
    }
    "docker" {
        $imageName = "hour50-multi-stage"
        $containerName = "hour50-multi-stage-test"
        $url = "http://localhost:18050/build"

        Remove-DockerContainerIfExists -Name $containerName
        Remove-DockerImageIfExists -Name $imageName
        try {
            Write-Host "Building multi-stage image $imageName"
            docker build --build-arg VERSION=1.0.0 --build-arg COMMIT=training -t $imageName .
            if ($LASTEXITCODE -ne 0) { throw "Docker build failed." }
            Write-Host "Starting verification container on http://localhost:18050"
            docker run -d --name $containerName -p 18050:8080 $imageName | Out-Null
            if ($LASTEXITCODE -ne 0) { throw "Docker run failed." }
            Start-Sleep -Seconds 3
            $response = Invoke-RestMethod -Uri $url
            $response | ConvertTo-Json -Depth 5
        }
        finally {
            Write-Host "Cleaning up container and image"
            Remove-DockerContainerIfExists -Name $containerName
            Remove-DockerImageIfExists -Name $imageName
        }
    }
}

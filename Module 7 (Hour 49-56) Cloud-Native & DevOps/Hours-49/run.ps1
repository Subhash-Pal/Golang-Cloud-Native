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
        $env:PORT = "18449"
        Write-Host "Running Hour 49 locally on http://localhost:18449"
        Write-Host "Endpoints: /  /healthz  /time"
        go run .
    }
    "docker" {
        $imageName = "hour49-basic-api"
        $containerName = "hour49-basic-api-test"
        $url = "http://localhost:18049/"

        Remove-DockerContainerIfExists -Name $containerName
        Remove-DockerImageIfExists -Name $imageName
        try {
            Write-Host "Building Docker image $imageName"
            docker build -t $imageName .
            if ($LASTEXITCODE -ne 0) { throw "Docker build failed." }
            Write-Host "Starting verification container on $url"
            docker run -d --name $containerName -p 18049:8080 $imageName | Out-Null
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

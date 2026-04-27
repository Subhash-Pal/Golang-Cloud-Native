param(
    [ValidateSet("local", "checks", "docker")]
    [string]$Mode = "checks"
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
        $env:PORT = "18455"
        Write-Host "Running Hour 55 locally on http://localhost:18455"
        Write-Host "Endpoints: /  /healthz"
        go run .
    }
    "checks" {
        Write-Host "Running tests and build checks"
        go test ./...
        go build .
    }
    "docker" {
        $imageName = "hour55-ci-api"
        $containerName = "hour55-ci-api-test"
        $url = "http://localhost:18055/healthz"

        Remove-DockerContainerIfExists -Name $containerName
        Remove-DockerImageIfExists -Name $imageName
        try {
            Write-Host "Building Docker image $imageName"
            docker build -t $imageName .
            if ($LASTEXITCODE -ne 0) { throw "Docker build failed." }
            Write-Host "Starting verification container on http://localhost:18055"
            docker run -d --name $containerName -p 18055:8080 $imageName | Out-Null
            if ($LASTEXITCODE -ne 0) { throw "Docker run failed." }
            Start-Sleep -Seconds 3
            $response = Invoke-RestMethod -Uri $url
            $response
        }
        finally {
            Write-Host "Cleaning up container and image"
            Remove-DockerContainerIfExists -Name $containerName
            Remove-DockerImageIfExists -Name $imageName
        }
    }
}

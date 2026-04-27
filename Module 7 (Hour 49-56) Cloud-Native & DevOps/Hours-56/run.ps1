param(
    [ValidateSet("local", "docker", "compose", "k8s")]
    [string]$Mode = "docker"
)

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

function Get-AvailablePort {
    param(
        [int[]]$Candidates
    )

    foreach ($port in $Candidates) {
        $listener = $null
        try {
            $listener = [System.Net.Sockets.TcpListener]::new([System.Net.IPAddress]::Loopback, $port)
            $listener.Start()
            return $port
        }
        catch {
        }
        finally {
            if ($listener) {
                $listener.Stop()
            }
        }
    }

    throw "No available host port found."
}

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

function Invoke-DockerComposeDownQuietly {
    $previousPreference = $ErrorActionPreference
    try {
        # Docker Compose often writes stop/remove progress to stderr even when the command succeeds.
        $ErrorActionPreference = "Continue"
        docker compose down --rmi local --remove-orphans *> $null
    }
    finally {
        $ErrorActionPreference = $previousPreference
    }
}

function Run-KubernetesApply {
    $imageRef = "docker.io/library/hour56-api:latest"

    Write-Host "Checking kubectl client"
    kubectl version --client
    if ($LASTEXITCODE -ne 0) { throw "kubectl client check failed." }

    Write-Host "Checking current Kubernetes context"
    kubectl config current-context
    if ($LASTEXITCODE -ne 0) { throw "kubectl context check failed." }

    Write-Host "Checking cluster connectivity"
    kubectl get nodes
    if ($LASTEXITCODE -ne 0) { throw "kubectl cannot reach the cluster." }

    Write-Host "Building Docker image hour56-api for Kubernetes"
    docker build -t hour56-api .
    if ($LASTEXITCODE -ne 0) { throw "Docker build failed for Kubernetes mode." }

    Write-Host "Confirming Kubernetes node can see $imageRef"
    docker exec desktop-control-plane ctr -n k8s.io images ls | findstr "docker.io/library/hour56-api:latest"
    if ($LASTEXITCODE -ne 0) { throw "Kubernetes node runtime cannot see $imageRef." }

    Write-Host "Applying Kubernetes manifests"
    kubectl apply -f .\k8s\
    if ($LASTEXITCODE -ne 0) { throw "kubectl apply failed." }

    Write-Host "Updating deployment image to $imageRef"
    kubectl set image deployment/hour56-api api=$imageRef
    if ($LASTEXITCODE -ne 0) { throw "kubectl set image failed." }

    Write-Host "Checking deployed resources"
    kubectl get deployment hour56-api
    kubectl get pods
    kubectl get svc

    kubectl wait --for=condition=available deployment/hour56-api --timeout=30s
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Deployment did not become available within 30 seconds."
        Write-Host "If you see ErrImageNeverPull or ImagePullBackOff, the cluster cannot see the local image yet."
        Write-Host "Check details with: kubectl describe pod -l app=hour56-api"
    }

    Write-Host "To test locally from Kubernetes, run: kubectl port-forward svc/hour56-api 8080:80"
}

switch ($Mode) {
    "local" {
        $env:PORT = "18456"
        Write-Host "Running Hour 56 locally on http://localhost:18456"
        Write-Host "Endpoints: /  /items  /orders  /healthz  /readyz"
        go run .
    }
    "docker" {
        $imageName = "hour56-api"
        $containerName = "hour56-api-test"
        $hostPort = Get-AvailablePort -Candidates @(18056, 18156, 18256, 18356, 19056)
        $url = "http://localhost:$hostPort/items"

        Remove-DockerContainerIfExists -Name $containerName
        Remove-DockerImageIfExists -Name $imageName
        try {
            Write-Host "Building Docker image $imageName"
            docker build -t $imageName .
            if ($LASTEXITCODE -ne 0) { throw "Docker build failed." }
            Write-Host "Starting verification container on http://localhost:$hostPort"
            docker run -d --name $containerName -p "${hostPort}:8080" $imageName | Out-Null
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
    "compose" {
        $hostPort = Get-AvailablePort -Candidates @(18056, 18156, 18256, 18356, 19056)
        $env:HOST_PORT = "$hostPort"
        try {
            Invoke-DockerComposeDownQuietly
            docker compose up -d --build
            if ($LASTEXITCODE -ne 0) { throw "Docker Compose startup failed." }

            Start-Sleep -Seconds 5
            Write-Host "Compose verification URL: http://localhost:$hostPort/items"
            $response = Invoke-RestMethod -Uri "http://localhost:$hostPort/items"
            $response | ConvertTo-Json -Depth 5
        }
        finally {
            Write-Host "Stopping stack and removing local images"
            Invoke-DockerComposeDownQuietly
            Remove-Item Env:HOST_PORT -ErrorAction SilentlyContinue
        }
    }
    "k8s" {
        Run-KubernetesApply
    }
}

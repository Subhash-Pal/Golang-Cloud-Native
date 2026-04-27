param(
    [ValidateSet("full", "local", "docker", "k8s")]
    [string]$Mode = "full"
)

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

function Remove-DockerContainerIfExists {
    param([string]$Name)

    $containerNames = docker ps -a --format "{{.Names}}"
    if ($containerNames -contains $Name) {
        docker rm -f $Name | Out-Null
    }
}

function Remove-DockerImageIfExists {
    param([string]$Name)

    $imageNames = docker images --format "{{.Repository}}:{{.Tag}}"
    if (($imageNames -contains $Name) -or ($imageNames -contains "${Name}:latest")) {
        docker image rm -f $Name | Out-Null
    }
}

function Run-DockerVerification {
    $imageName = "hour53-api"
    $containerName = "hour53-api-test"
    $url = "http://localhost:18053/config"

    Remove-DockerContainerIfExists -Name $containerName
    Remove-DockerImageIfExists -Name $imageName

    try {
        Write-Host "Building Docker image $imageName"
        docker build -t $imageName .
        if ($LASTEXITCODE -ne 0) { throw "Docker build failed." }

        Write-Host "Starting verification container on http://localhost:18053"
        docker run -d --name $containerName -p 18053:8080 -e APP_MODE=production -e LOG_LEVEL=info $imageName | Out-Null
        if ($LASTEXITCODE -ne 0) { throw "Docker run failed." }

        Start-Sleep -Seconds 3
        $response = Invoke-RestMethod -Uri $url
        Write-Host "Docker verification response:"
        $response | ConvertTo-Json -Depth 5
    }
    finally {
        Write-Host "Cleaning up container and image"
        Remove-DockerContainerIfExists -Name $containerName
        Remove-DockerImageIfExists -Name $imageName
    }
}

function Test-KubernetesAvailable {
    kubectl version --client *> $null
    if ($LASTEXITCODE -ne 0) { return $false }
    kubectl config current-context *> $null
    if ($LASTEXITCODE -ne 0) { return $false }
    kubectl get nodes *> $null
    return ($LASTEXITCODE -eq 0)
}

function Run-KubernetesApply {
    $imageRef = "docker.io/library/hour53-api:latest"

    Write-Host "Checking kubectl client"
    kubectl version --client
    if ($LASTEXITCODE -ne 0) { throw "kubectl client check failed." }

    Write-Host "Checking current Kubernetes context"
    kubectl config current-context
    if ($LASTEXITCODE -ne 0) { throw "kubectl context check failed." }

    Write-Host "Checking cluster connectivity"
    kubectl get nodes
    if ($LASTEXITCODE -ne 0) { throw "kubectl cannot reach the cluster." }

    Write-Host "Building Docker image hour53-api for Kubernetes"
    docker build -t hour53-api .
    if ($LASTEXITCODE -ne 0) { throw "Docker build failed for Kubernetes mode." }

    Write-Host "Confirming Kubernetes node can see $imageRef"
    docker exec desktop-control-plane ctr -n k8s.io images ls | findstr "docker.io/library/hour53-api:latest"
    if ($LASTEXITCODE -ne 0) { throw "Kubernetes node runtime cannot see $imageRef." }

    Write-Host "Applying Kubernetes manifests"
    kubectl apply -f .\k8s\
    if ($LASTEXITCODE -ne 0) { throw "kubectl apply failed." }

    Write-Host "Updating deployment image to $imageRef"
    kubectl set image deployment/hour53-api api=$imageRef
    if ($LASTEXITCODE -ne 0) { throw "kubectl set image failed." }

    Write-Host "Checking deployed resources"
    kubectl get deployment hour53-api
    kubectl get pods
    kubectl get svc
    kubectl get configmap hour53-config

    kubectl wait --for=condition=available deployment/hour53-api --timeout=30s
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Deployment did not become available within 30 seconds."
        Write-Host "If you see ErrImageNeverPull or ImagePullBackOff, the cluster cannot see the local image yet."
        Write-Host "Check details with: kubectl describe pod -l app=hour53-api"
    }

    Write-Host "To test locally from Kubernetes, run: kubectl port-forward svc/hour53-api 8080:80"
}

switch ($Mode) {
    "full" {
        Write-Host "Running Hour 53 full setup"
        Write-Host "Step 1: Docker verification"
        Run-DockerVerification

        Write-Host "Step 2: Kubernetes verification"
        if (Test-KubernetesAvailable) {
            Run-KubernetesApply
        } else {
            Write-Host "Kubernetes check skipped because kubectl is not configured for a reachable cluster."
            Write-Host "When your cluster is ready, run: powershell -ExecutionPolicy Bypass -File .\run.ps1 -Mode k8s"
        }
    }
    "local" {
        $env:APP_MODE = "development"
        $env:LOG_LEVEL = "debug"
        $env:PORT = "18453"
        Write-Host "Running Hour 53 locally on http://localhost:18453"
        Write-Host "Endpoints: /  /config  /healthz"
        go run .
    }
    "docker" {
        Run-DockerVerification
    }
    "k8s" {
        Run-KubernetesApply
    }
}

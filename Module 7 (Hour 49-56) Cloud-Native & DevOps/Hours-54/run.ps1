param(
    [ValidateSet("local", "docker", "k8s")]
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

function Run-KubernetesApply {
    $imageRef = "docker.io/library/hour54-api:latest"

    Write-Host "Checking kubectl client"
    kubectl version --client
    if ($LASTEXITCODE -ne 0) { throw "kubectl client check failed." }

    Write-Host "Checking current Kubernetes context"
    kubectl config current-context
    if ($LASTEXITCODE -ne 0) { throw "kubectl context check failed." }

    Write-Host "Checking cluster connectivity"
    kubectl get nodes
    if ($LASTEXITCODE -ne 0) { throw "kubectl cannot reach the cluster." }

    Write-Host "Building Docker image hour54-api for Kubernetes"
    docker build -t hour54-api .
    if ($LASTEXITCODE -ne 0) { throw "Docker build failed for Kubernetes mode." }

    Write-Host "Confirming Kubernetes node can see $imageRef"
    docker exec desktop-control-plane ctr -n k8s.io images ls | findstr "docker.io/library/hour54-api:latest"
    if ($LASTEXITCODE -ne 0) { throw "Kubernetes node runtime cannot see $imageRef." }

    Write-Host "Applying Kubernetes manifests"
    kubectl apply -f .\k8s\
    if ($LASTEXITCODE -ne 0) { throw "kubectl apply failed." }

    Write-Host "Updating deployment image to $imageRef"
    kubectl set image deployment/hour54-api api=$imageRef
    if ($LASTEXITCODE -ne 0) { throw "kubectl set image failed." }

    Write-Host "Checking deployed resources"
    kubectl get deployment hour54-api
    kubectl get pods
    kubectl get svc
    kubectl get hpa hour54-api

    kubectl wait --for=condition=available deployment/hour54-api --timeout=30s
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Deployment did not become available within 30 seconds."
        Write-Host "If you see ErrImageNeverPull or ImagePullBackOff, the cluster cannot see the local image yet."
        Write-Host "Check details with: kubectl describe pod -l app=hour54-api"
    }

    Write-Host "To test locally from Kubernetes, run: kubectl port-forward svc/hour54-api 8080:80"
}

switch ($Mode) {
    "local" {
        $env:STARTUP_DELAY_SECONDS = "5"
        $env:PORT = "18454"
        Write-Host "Running Hour 54 locally on http://localhost:18454"
        Write-Host "Endpoints: /  /healthz  /readyz  /metrics"
        go run .
    }
    "docker" {
        $imageName = "hour54-api"
        $containerName = "hour54-api-test"
        $url = "http://localhost:18054/readyz"

        Remove-DockerContainerIfExists -Name $containerName
        Remove-DockerImageIfExists -Name $imageName
        try {
            Write-Host "Building Docker image $imageName"
            docker build -t $imageName .
            if ($LASTEXITCODE -ne 0) { throw "Docker build failed." }
            Write-Host "Starting verification container on http://localhost:18054"
            docker run -d --name $containerName -p 18054:8080 -e STARTUP_DELAY_SECONDS=2 $imageName | Out-Null
            if ($LASTEXITCODE -ne 0) { throw "Docker run failed." }
            Start-Sleep -Seconds 4
            $response = Invoke-RestMethod -Uri $url
            $response | ConvertTo-Json -Depth 5
        }
        finally {
            Write-Host "Cleaning up container and image"
            Remove-DockerContainerIfExists -Name $containerName
            Remove-DockerImageIfExists -Name $imageName
        }
    }
    "k8s" {
        Run-KubernetesApply
    }
}

param(
    [ValidateSet("local", "docker", "k8s")]
    [string]$Mode = "docker"
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

function Build-LocalBinary {
    param([string]$BinaryName)

    $binDir = Join-Path $PSScriptRoot ".bin"
    if (-not (Test-Path $binDir)) {
        New-Item -ItemType Directory -Path $binDir | Out-Null
    }

    $binaryPath = Join-Path $binDir "${BinaryName}.exe"
    Write-Host "Building local binary at $binaryPath"
    go build -o $binaryPath .
    if ($LASTEXITCODE -ne 0) { throw "go build failed." }

    return $binaryPath
}

function Get-CurrentKubernetesContext {
    $context = kubectl config current-context 2>$null
    if ($LASTEXITCODE -ne 0) { throw "kubectl context check failed." }
    return $context.Trim()
}

function Import-ImageToCluster {
    param(
        [string]$ImageRef,
        [string]$ContextName
    )

    if ($ContextName -eq "minikube") {
        Write-Host "Loading $ImageRef into minikube"
        minikube image load $ImageRef
        if ($LASTEXITCODE -ne 0) { throw "minikube image load failed." }
        return
    }

    if ($ContextName -like "kind-*") {
        $clusterName = $ContextName.Substring(5)
        Write-Host "Loading $ImageRef into kind cluster $clusterName"
        kind load docker-image $ImageRef --name $clusterName
        if ($LASTEXITCODE -ne 0) { throw "kind load docker-image failed." }
        return
    }

    if ($ContextName -eq "docker-desktop") {
        Write-Host "Checking Docker Desktop cluster image visibility"
        docker exec docker-desktop-control-plane crictl images | findstr "hour54-api"
        if ($LASTEXITCODE -ne 0) {
            throw "Docker Desktop cluster cannot see $ImageRef."
        }
        return
    }

    Write-Host "No automatic image import is configured for context '$ContextName'."
    Write-Host "The deployment will use $ImageRef and relies on the cluster being able to access that image."
}

function Run-KubernetesApply {
    $imageName = "hour54-api"
    $imageRef = "docker.io/library/hour54-api:latest"

    Write-Host "Checking kubectl client"
    kubectl version --client
    if ($LASTEXITCODE -ne 0) { throw "kubectl client check failed." }

    Write-Host "Checking current Kubernetes context"
    $contextName = Get-CurrentKubernetesContext
    Write-Host $contextName

    Write-Host "Checking cluster connectivity"
    kubectl get nodes
    if ($LASTEXITCODE -ne 0) { throw "kubectl cannot reach the cluster." }

    Write-Host "Building Docker image $imageName for Kubernetes"
    docker build -t $imageName -t $imageRef .
    if ($LASTEXITCODE -ne 0) { throw "Docker build failed for Kubernetes mode." }

    Import-ImageToCluster -ImageRef $imageRef -ContextName $contextName

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

    kubectl wait --for=condition=available deployment/hour54-api --timeout=90s
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Deployment did not become available within 90 seconds."
        Write-Host "Check pod details with: kubectl describe pod -l app=hour54-api"
        throw "Deployment verification failed."
    }

    Write-Host "To test locally from Kubernetes, run: kubectl port-forward svc/hour54-api 8080:80"
}

switch ($Mode) {
    "local" {
        $env:STARTUP_DELAY_SECONDS = "5"
        $env:PORT = "18454"
        $binaryPath = Build-LocalBinary -BinaryName "hour54-api"
        Write-Host "Running Hour 54 locally on http://localhost:18454"
        Write-Host "Endpoints: /  /healthz  /readyz  /metrics"
        & $binaryPath
    }
    "docker" {
        $imageName = "hour54-api"
        $imageRef = "docker.io/library/hour54-api:latest"
        $containerName = "hour54-api-test"
        $url = "http://localhost:18054/readyz"

        Remove-DockerContainerIfExists -Name $containerName
        Remove-DockerImageIfExists -Name $imageName
        try {
            Write-Host "Building Docker image $imageName"
            docker build -t $imageName -t $imageRef .
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
            Remove-DockerImageIfExists -Name $imageRef
        }
    }
    "k8s" {
        Run-KubernetesApply
    }
}

# -------------------------------
# Ensure script runs from its own directory
# -------------------------------
Set-Location -Path $PSScriptRoot

Write-Host "===== STEP 1: Start Minikube ====="
minikube start

Write-Host "===== STEP 2: Configure Docker to use Minikube ====="
minikube docker-env --shell powershell | Invoke-Expression

Write-Host "===== STEP 3: Build Docker Images ====="
docker build -t quote-service ./quote-service
docker build -t api ./api

Write-Host "===== STEP 4: Deploy to Kubernetes ====="
kubectl apply -f k8s.yaml

Write-Host "===== STEP 5: Wait for Pods ====="
kubectl wait --for=condition=ready pod -l app=quote-service --timeout=90s
kubectl wait --for=condition=ready pod -l app=api --timeout=90s

Write-Host "===== STEP 6: Verify Pods ====="
kubectl get pods

Write-Host "===== STEP 7: Verify Services ====="
kubectl get svc

Write-Host "===== STEP 8: Get Service URL ====="
Write-Host "===== STEP 8: Start Minikube Tunnel ====="

# Run tunnel in background (non-blocking)
Start-Process powershell -ArgumentList "minikube service api" 

Start-Sleep -Seconds 5

Write-Host "===== STEP 9: Get URL ====="

# Extract only URL line cleanly
$url = (minikube service api --url | Select-Object -First 1).Trim()

if (-not $url -or $url -notmatch "^http") {
    Write-Host "❌ Could not detect URL automatically"
    Write-Host "Run manually: minikube service api --url"
    exit 1
}

Write-Host "API URL: $url"

Write-Host "===== STEP 10: Open Browser ====="
Start-Process $url
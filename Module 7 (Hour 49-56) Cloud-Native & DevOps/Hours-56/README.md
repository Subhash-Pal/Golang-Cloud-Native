# Hour 56: Mock Test - Deploy Containerized API

## Files

- `main.go`
- `main_test.go`
- `Dockerfile`
- `docker-compose.yml`
- `k8s\deployment.yaml`
- `k8s\service.yaml`
- `run-local.ps1`
- `run-docker.ps1`
- `run-compose.ps1`
- `run-k8s.ps1`
- `run.ps1`

## What This Demo Does

Hour 56 is a small inventory and order API that can run in four ways:

- local Go server
- Docker verification run
- Docker Compose verification run
- Kubernetes deployment on Docker Desktop

Main routes:

- `GET /`
- `GET /items`
- `POST /items`
- `GET /orders`
- `POST /orders`
- `GET /healthz`
- `GET /readyz`

## PowerShell Files

- `run-local.ps1`: starts the API locally and keeps it running until you stop it
- `run-docker.ps1`: builds a temporary image, verifies the API, then removes the container and image
- `run-compose.ps1`: starts a temporary Compose stack, verifies the API, then removes the stack and local image
- `run-k8s.ps1`: builds the image, applies the Kubernetes manifests, and waits for the deployment to become available
- `run.ps1`: shared runner used by the scripts above

You can also call the shared runner directly:

```powershell
powershell -ExecutionPolicy Bypass -File .\run.ps1 -Mode local
powershell -ExecutionPolicy Bypass -File .\run.ps1 -Mode docker
powershell -ExecutionPolicy Bypass -File .\run.ps1 -Mode compose
powershell -ExecutionPolicy Bypass -File .\run.ps1 -Mode k8s
```

## Run Local

```powershell
powershell -ExecutionPolicy Bypass -File .\run-local.ps1
```

The local server keeps running. Stop it with `Ctrl+C` after testing.

Then open:

```text
http://localhost:18456/
http://localhost:18456/items
http://localhost:18456/orders
http://localhost:18456/healthz
http://localhost:18456/readyz
```

## Run Docker Verification

```powershell
powershell -ExecutionPolicy Bypass -File .\run-docker.ps1
```

The script picks an available host port automatically, starts a temporary container, verifies `GET /items`, prints the response, and then removes the container and image.

Use the localhost URL printed by the script output.

## Run Compose Verification

```powershell
powershell -ExecutionPolicy Bypass -File .\run-compose.ps1
```

The script picks an available host port automatically, starts the Compose stack, verifies `GET /items`, prints the response, and then removes the stack and local image.

Use the localhost URL printed by the script output.

## Run Kubernetes

```powershell
powershell -ExecutionPolicy Bypass -File .\run-k8s.ps1
```

The Kubernetes runner:

- builds `hour56-api`
- confirms the Docker Desktop Kubernetes node can see `docker.io/library/hour56-api:latest`
- applies the manifests
- updates the deployment image
- waits for the deployment to become available

After deployment, check:

```powershell
kubectl get deployment hour56-api
kubectl get pods -l app=hour56-api
kubectl get svc hour56-api
```

Start local testing:

```powershell
kubectl port-forward svc/hour56-api 8080:80
```

Then open:

```text
http://localhost:8080/
http://localhost:8080/items
http://localhost:8080/orders
http://localhost:8080/healthz
http://localhost:8080/readyz
```

### Example API Calls

Create an item:

```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:18456/items -ContentType "application/json" -Body '{"name":"Mouse","stock":15}'
```

Create an order:

```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:18456/orders -ContentType "application/json" -Body '{"item_id":1,"quantity":2}'
```

## Recommended Test Order

1. Run `run-local.ps1`
2. Run `run-docker.ps1`
3. Run `run-compose.ps1`
4. Run `run-k8s.ps1`

This order moves from the simplest setup to the full Kubernetes deployment.

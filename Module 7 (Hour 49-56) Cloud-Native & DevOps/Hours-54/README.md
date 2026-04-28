# Hour 54: Health Probes and Scaling

## Files

- `main.go`
- `Dockerfile`
- `k8s\deployment.yaml`
- `k8s\service.yaml`
- `k8s\hpa.yaml`
- `run-local.ps1`
- `run-docker.ps1`
- `run-k8s.ps1`
- `run.ps1`

## Run Local

```powershell
powershell -ExecutionPolicy Bypass -File .\run-local.ps1
```

Then test:

```text
http://localhost:18454/
http://localhost:18454/healthz
http://localhost:18454/readyz
http://localhost:18454/metrics
```

The local runner now builds a stable `.bin\hour54-api.exe` and starts that binary instead of using `go run`, which helps on machines where Windows App Control blocks temporary Go build-cache executables.

`/readyz` is expected to return `503` briefly during the startup warm-up window and then switch to `200`.

## Run Docker Verification

```powershell
powershell -ExecutionPolicy Bypass -File .\run-docker.ps1
```

The script verifies `http://localhost:18054/readyz` and then cleans up.

## Run Kubernetes

```powershell
powershell -ExecutionPolicy Bypass -File .\run-k8s.ps1
```

The Kubernetes runner now:

- builds `hour54-api` and tags `docker.io/library/hour54-api:latest`
- auto-loads the image into `minikube` or `kind` when those contexts are active
- checks Docker Desktop image visibility when the active context is `docker-desktop`
- applies the manifests
- waits for the deployment to become available

After deployment, use:

```powershell
kubectl get deployment hour54-api
kubectl get pods
kubectl get svc
kubectl get hpa hour54-api
```

Start local testing:

```powershell
kubectl port-forward svc/hour54-api 8080:80
```

Then open:

```text
http://localhost:8080/
http://localhost:8080/healthz
http://localhost:8080/readyz
http://localhost:8080/metrics
```

If the deployment does not become available, check:

```powershell
kubectl describe pod -l app=hour54-api
```

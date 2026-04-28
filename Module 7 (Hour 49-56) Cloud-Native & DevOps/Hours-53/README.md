# Hour 53: Writing Deployment and Service YAML

## Files

- `main.go`
- `Dockerfile`
- `k8s\configmap.yaml`
- `k8s\deployment.yaml`
- `k8s\service.yaml`
- `run-full.ps1`
- `run-local.ps1`
- `run-docker.ps1`
- `run-k8s.ps1`
- `run.ps1`

## Run Everything With One Command

```powershell
powershell -ExecutionPolicy Bypass -File .\run-full.ps1
```

This runs Docker verification first and then runs the Kubernetes apply flow if your cluster is reachable.

The local runner now builds a stable `.bin\hour53-api.exe` and starts that binary instead of using `go run`, which helps on machines where Windows App Control blocks temporary Go build-cache executables.

## Run Local

```powershell
powershell -ExecutionPolicy Bypass -File .\run-local.ps1
```

Open:

```text
http://localhost:18453/
http://localhost:18453/config
http://localhost:18453/healthz
```

## Run Docker Verification

```powershell
powershell -ExecutionPolicy Bypass -File .\run-docker.ps1
```

## Run Kubernetes

```powershell
powershell -ExecutionPolicy Bypass -File .\run-k8s.ps1
```

The Kubernetes runner now:

- builds `hour53-api` and tags `docker.io/library/hour53-api:latest`
- auto-loads the image into `minikube` or `kind` when those contexts are active
- checks Docker Desktop image visibility when the active context is `docker-desktop`
- applies the manifests
- waits for the deployment to become available

After deployment, use:

```powershell
kubectl get deployment hour53-api
kubectl get pods
kubectl get svc
kubectl get configmap hour53-config
```

Start local testing:

```powershell
kubectl port-forward svc/hour53-api 8080:80
```

Then open:

```text
http://localhost:8080/
http://localhost:8080/config
http://localhost:8080/healthz
```

If the deployment does not become available, check:

```powershell
kubectl describe pod -l app=hour53-api
```

# Module 7: Cloud-Native & DevOps

This module contains one self-contained training project per hour from 49 to 56.

## Hours Covered

| Hour | Topic | Folder |
| --- | --- | --- |
| 49 | Docker fundamentals | `Hours-49` |
| 50 | Multi-stage Docker builds | `Hours-50` |
| 51 | Docker Compose setup | `Hours-51` |
| 52 | Kubernetes architecture | `Hours-52` |
| 53 | Deployment and Service YAML | `Hours-53` |
| 54 | Health probes and scaling | `Hours-54` |
| 55 | CI/CD pipeline setup | `Hours-55` |
| 56 | Mock test: deploy containerized API | `Hours-56` |

## Runner Files

Each hour now uses separate PowerShell files for each supported configuration.

Common script names:

- `run-local.ps1`
- `run-docker.ps1`
- `run-compose.ps1`
- `run-k8s.ps1`
- `run-checks.ps1`
- `run-full.ps1`
- `run.ps1`

Not every hour uses every script name. Each folder only includes the files needed for that topic.

## Basic Steps

1. Open PowerShell.
2. Move into the hour folder you want to run.
3. Run the matching script file for that configuration.
4. If you use a local script, keep the terminal open and browse the localhost URL shown by the script.
5. If you use a Docker or Compose script, the script verifies the app automatically and then cleans up.

## Example Commands

```powershell
powershell -ExecutionPolicy Bypass -File .\run-local.ps1
powershell -ExecutionPolicy Bypass -File .\run-docker.ps1
powershell -ExecutionPolicy Bypass -File .\run-compose.ps1
powershell -ExecutionPolicy Bypass -File .\run-k8s.ps1
powershell -ExecutionPolicy Bypass -File .\run-checks.ps1
powershell -ExecutionPolicy Bypass -File .\run-full.ps1
```

For Kubernetes-enabled hours, use the split Kubernetes runner:

```powershell
powershell -ExecutionPolicy Bypass -File .\run-k8s.ps1
```

The Kubernetes runner now builds the image, confirms the Docker Desktop node can see it, applies the manifests, and waits for the deployment to become available.

After the Kubernetes runner finishes, test the service with:

```powershell
kubectl get deployment
kubectl get pods
kubectl get svc
kubectl port-forward svc/<service-name> 8080:80
```

Then open the relevant localhost URLs for that hour in your browser.

# Kubernetes + Golang Microservices (Minikube) — Complete Workflow

## 📌 Overview

This project demonstrates a **production-style microservice architecture** using:

* Golang services
* Docker containers
* Kubernetes (Minikube)
* Internal service discovery
* External access via NodePort

---

## 🏗 Architecture

```
Client
   ↓
[ API Service (NodePort:30007) ]
   ↓
[ quote-service (ClusterIP) ]
```

### Components

| Component     | Description                      |
| ------------- | -------------------------------- |
| API Service   | Entry point, calls quote-service |
| quote-service | Internal service providing data  |
| Kubernetes    | Orchestrates containers          |
| Minikube      | Local Kubernetes cluster         |
| NodePort      | Exposes API externally           |

---

## ⚙️ Project Structure

```
go-k8s-demo/
│
├── api/
│   ├── main.go
│   ├── Dockerfile
│   └── go.mod
│
├── quote-service/
│   ├── main.go
│   ├── Dockerfile
│   └── go.mod
│
├── k8s.yaml
└── run.ps1
```

---

## 🚀 Execution Workflow

### Step 1 — Start Cluster

```powershell
minikube start
```

---

### Step 2 — Configure Docker Environment

```powershell
minikube docker-env --shell powershell | Invoke-Expression
```

This ensures images are built **inside Minikube**.

---

### Step 3 — Build Services

```powershell
docker build -t quote-service ./quote-service
docker build -t api ./api
```

---

### Step 4 — Deploy to Kubernetes

```powershell
kubectl apply -f k8s.yaml
```

---

### Step 5 — Verify Pods

```powershell
kubectl get pods
```

Expected:

```
STATUS = Running
RESTARTS = 0
```

---

### Step 6 — Verify Services

```powershell
kubectl get svc
```

Expected:

* `quote-service` → ClusterIP
* `api` → NodePort

---

### Step 7 — Get Minikube IP

```powershell
minikube ip
```

Example:

```
192.168.49.2
```

---

### Step 8 — Access API

```powershell
http://<MINIKUBE_IP>:30007
```

Example:

```powershell
Start-Process "http://192.168.49.2:30007"
```

---

### Step 9 — Test API

```powershell
Invoke-RestMethod -Uri "http://192.168.49.2:30007"
```

Expected response:

```json
{
  "source": "api",
  "quote": {
    "quote": "Kubernetes brings orchestration to your containers.",
    "timestamp": "..."
  },
  "timestamp": "..."
}
```

---

## 🔁 Automation

Run everything using:

```powershell
.\run.ps1
```

This script handles:

* Cluster startup
* Docker build
* Kubernetes deployment
* Pod readiness

---

## ⚠️ Important Notes

### 1. Avoid `minikube service` for automation

Reason:

* Requires terminal tunnel
* Blocks execution
* Not script-friendly

### 2. Use NodePort Instead

```
http://<MINIKUBE_IP>:30007
```

✔ Stable
✔ Script-friendly
✔ No tunnel required

---

### 3. Pod Restart Rule

Check:

```powershell
kubectl get pods
```

| Condition              | Meaning           |
| ---------------------- | ----------------- |
| Running + Restarts = 0 | Healthy           |
| Restarts increasing    | Application issue |

---

### 4. Minikube IP May Change

After restart:

```powershell
minikube ip
```

---

## 🧠 Key Concepts Learned

### Kubernetes Fundamentals

* Deployment
* Service (ClusterIP, NodePort)
* Pod lifecycle
* DNS-based service discovery

---

### Microservice Communication

```
api → http://quote-service:8081
```

✔ Internal DNS resolution
✔ Decoupled services

---

### Container Behavior

* Containers must **not exit unexpectedly**
* Services must handle failures gracefully
* Restart loops indicate application issues

---

## 🚀 Future Improvements

### 1. Add Health Probes

```yaml
livenessProbe:
  httpGet:
    path: /
    port: 8080
```

---

### 2. Use ConfigMap

Externalize:

```
QUOTE_SERVICE_URL
```

---

### 3. Scale Application

```powershell
kubectl scale deployment api --replicas=3
```

---

### 4. Replace NodePort with Ingress

Benefits:

* No IP + port
* Domain-based routing
* Production-ready

---

### 5. Add Observability

```powershell
kubectl logs -f deployment/api
kubectl top pods
```

---

## 🧠 Final Insight

This setup represents a **complete cloud-native pipeline**:

```
Build → Deploy → Discover → Expose → Validate
```

You are no longer just running containers—you are operating a **Kubernetes-based distributed system**.

---

## ✅ Status

| Layer             | Status |
| ----------------- | ------ |
| Docker Build      | ✔      |
| Kubernetes Deploy | ✔      |
| Service Discovery | ✔      |
| External Access   | ✔      |
| Stability         | ✔      |

---

## 📌 Summary

You have successfully implemented:

* Multi-service Golang system
* Kubernetes deployment model
* Internal service communication
* External API exposure
* Debugging and stabilization workflow

This is the **foundation of real-world backend systems**.

---

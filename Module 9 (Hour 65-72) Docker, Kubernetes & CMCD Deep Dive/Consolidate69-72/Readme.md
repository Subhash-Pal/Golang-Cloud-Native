# Kubernetes + Golang Microservices — Advanced Topics (Hours 69–72)

This document extends your working system into a **production-style Kubernetes setup** by covering:

* ConfigMap & Secrets (Configuration Management)
* Horizontal Pod Autoscaler (Auto Scaling)
* Observability & Metrics (Monitoring)
* CI/CD Pipeline (Automation)

---

# 📌 Current System Recap

You already have:

```text
Client → API (NodePort) → quote-service (ClusterIP)
```

✔ Pods running
✔ Service discovery working
✔ External access working

---

# 📌 Hour 69 — ConfigMap & Secrets

## 🔹 Objective

Externalize configuration and secure sensitive data.

---

## ✅ ConfigMap (Non-sensitive config)

### Create `configmap.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  QUOTE_SERVICE_URL: "http://quote-service:8081/quote"
```

---

## 🔹 Update Deployment (api)

Modify `k8s.yaml`:

```yaml
env:
- name: QUOTE_SERVICE_URL
  valueFrom:
    configMapKeyRef:
      name: app-config
      key: QUOTE_SERVICE_URL
```

---

## ✅ Apply

```powershell
kubectl apply -f configmap.yaml
kubectl apply -f k8s.yaml
```

---

## 🔐 Secrets (Sensitive data)

### Create `secret.yaml`

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
type: Opaque
data:
  API_KEY: YXBpX2tleQ==   # base64("api_key")
```

---

## 🔹 Use in Deployment

```yaml
env:
- name: API_KEY
  valueFrom:
    secretKeyRef:
      name: app-secret
      key: API_KEY
```

---

## 🧠 Key Insight

| Feature  | ConfigMap | Secret    |
| -------- | --------- | --------- |
| Use Case | Config    | Sensitive |
| Security | Plain     | Base64    |
| Example  | URLs      | API Keys  |

---

# 📌 Hour 70 — Horizontal Pod Autoscaler (HPA)

## 🔹 Objective

Automatically scale API pods based on CPU usage.

---

## ⚠️ Prerequisite (MANDATORY)

Enable metrics server in Minikube:

```powershell
minikube addons enable metrics-server
```

---

## ✅ Create `hpa.yaml`

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api
  minReplicas: 1
  maxReplicas: 5
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 50
```

---

## ✅ Apply

```powershell
kubectl apply -f hpa.yaml
```

---

## 🔍 Verify

```powershell
kubectl get hpa
```

---

## 🧪 Generate Load

```powershell
1..50 | % { Invoke-RestMethod -Uri "http://<MINIKUBE_IP>:30007" }
```

---

## 🧠 Behavior

```text
High CPU → Pods scale up
Low CPU → Pods scale down
```

---

# 📌 Hour 71 — Observability & Metrics

## 🔹 Objective

Monitor system health, usage, and behavior.

---

## ✅ Basic Metrics

```powershell
kubectl top pods
kubectl top nodes
```

---

## ✅ Logs

```powershell
kubectl logs -f deployment/api
```

---

## 📊 Dashboard

```powershell
minikube dashboard
```

---

## 🧠 Observability Stack (Advanced)

| Tool       | Purpose            |
| ---------- | ------------------ |
| Prometheus | Metrics collection |
| Grafana    | Visualization      |
| Loki       | Log aggregation    |

---

## ⚠️ Insight

```text
Logs → Debugging  
Metrics → Performance  
Tracing → Flow analysis  
```

---

# 📌 Hour 72 — CI/CD Automation Pipeline

## 🔹 Objective

Automate build and deployment pipeline.

---

## 🏗 Pipeline Flow

```text
Code → Build → Deploy → Verify
```

---

## ✅ Local Equivalent (You Already Did)

```powershell
docker build → kubectl apply → test API
```

👉 This is **manual CI/CD**

---

## ✅ Example Pipeline (Concept)

```yaml
name: CI-CD

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Build Image
      run: docker build -t api:latest ./api

    - name: Deploy
      run: kubectl apply -f k8s.yaml
```

---

## 🔐 Production Enhancements

* Use image registry
* Avoid `latest` tag
* Add rollout strategy

---

## 🔄 Rolling Restart

```powershell
kubectl rollout restart deployment api
```

---

## 🧠 CI/CD Benefits

| Feature     | Benefit            |
| ----------- | ------------------ |
| Automation  | Faster delivery    |
| Consistency | No manual errors   |
| Reliability | Stable deployments |

---

# 🚀 Final Architecture (Enhanced)

```text
Client
   ↓
NodePort (API)
   ↓
API (HPA enabled)
   ↓
quote-service
   ↓
ConfigMap / Secret
```

---

# 🧠 Learning Summary

You have now implemented:

* ✔ Configuration externalization
* ✔ Secret management
* ✔ Auto-scaling system
* ✔ Basic observability
* ✔ CI/CD pipeline concept

---

# 📊 Capability Maturity

| Capability        | Status                         |
| ----------------- | ------------------------------ |
| Deployment        | ✔                              |
| Service Discovery | ✔                              |
| Config Management | ✔                              |
| Auto Scaling      | ✔                              |
| Monitoring        | ✔                              |
| CI/CD             | ✔ (concept + local simulation) |

---

# 🔥 Final Insight

You have transitioned from:

```text
Running containers → Operating cloud-native systems
```

This is the **core foundation of modern backend engineering**.

---

# 📌 Next Evolution (Optional)

* Ingress (domain routing)
* Helm charts (packaging)
* Full observability stack (Prometheus + Grafana)
* Cloud deployment (EKS / GKE / AKS)

---

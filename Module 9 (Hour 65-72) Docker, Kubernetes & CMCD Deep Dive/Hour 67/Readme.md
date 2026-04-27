# Advanced Kubernetes Topics — ConfigMap, HPA, Observability & CI/CD

This document extends the base setup by introducing **production-grade capabilities**:

* Configuration management
* Auto-scaling
* Observability
* CI/CD automation

---

# 📌 Hour 69 — ConfigMap & Secrets

## 🔹 Why needed

Hardcoding values like:

```text
QUOTE_SERVICE_URL
```

is not scalable or secure.

---

## ✅ ConfigMap (Non-sensitive data)

### Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  QUOTE_SERVICE_URL: "http://quote-service:8081/quote"
```

---

## 🔹 Use in Deployment

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
```

---

## 🔐 Secrets (Sensitive data)

Used for:

* passwords
* tokens
* API keys

---

### Example

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
type: Opaque
data:
  API_KEY: YXBpX2tleQ==   # base64 encoded
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

## ⚠️ Key Difference

| Feature  | ConfigMap | Secret    |
| -------- | --------- | --------- |
| Purpose  | Config    | Sensitive |
| Encoding | Plain     | Base64    |
| Security | Low       | Higher    |

---

# 📌 Hour 70 — Horizontal Pod Autoscaler (HPA)

## 🔹 Objective

Automatically scale pods based on CPU usage.

---

## ⚠️ Prerequisite

Enable metrics server in Minikube:

```powershell
minikube addons enable metrics-server
```

---

## ✅ Define HPA

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
High load → CPU ↑ → Pods scale up
Low load → Pods scale down
```

---

# 📌 Hour 71 — Observability & Metrics

## 🔹 Why Observability?

To monitor:

* performance
* resource usage
* errors

---

## ✅ Enable Metrics

```powershell
kubectl top pods
kubectl top nodes
```

---

## 🔍 Logs

```powershell
kubectl logs -f deployment/api
```

---

## 📊 Kubernetes Dashboard

```powershell
minikube dashboard
```

---

## 🧠 Observability Stack (Real-world)

| Tool       | Purpose            |
| ---------- | ------------------ |
| Prometheus | Metrics collection |
| Grafana    | Visualization      |
| Loki       | Log aggregation    |

---

## ⚠️ Key Insight

```text
Logs → Debugging
Metrics → Performance
Tracing → Flow analysis
```

---

# 📌 Hour 72 — CI/CD Automation Pipeline

## 🔹 Objective

Automate:

```text
Code → Build → Test → Deploy
```

---

## 🏗 Pipeline Flow

```text
Git Push
   ↓
Build Docker Image
   ↓
Push Image
   ↓
Deploy to Kubernetes
```

---

## ✅ Example (GitHub Actions)

```yaml
name: CI-CD Pipeline

on:
  push:
    branches:
      - main

jobs:
  build-deploy:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout Code
      uses: actions/checkout@v3

    - name: Build Docker Image
      run: docker build -t api:latest ./api

    - name: Deploy to Kubernetes
      run: kubectl apply -f k8s.yaml
```

---

## 🔐 Production Enhancements

* Use container registry (Docker Hub / ECR)
* Use image tags (not `latest`)
* Add rollout strategy

---

## 🔄 Rolling Update

```powershell
kubectl rollout restart deployment api
```

---

## 🧠 CI/CD Benefits

| Feature     | Benefit                 |
| ----------- | ----------------------- |
| Automation  | No manual deploy        |
| Consistency | Same process every time |
| Speed       | Faster delivery         |
| Reliability | Reduced human error     |

---

# 🚀 Final Architecture (After Enhancements)

```text
Client
   ↓
Ingress (future)
   ↓
API (HPA enabled)
   ↓
quote-service
   ↓
ConfigMap / Secret
```

---

# 🧠 Final Learning Summary

You now understand:

* Configuration externalization
* Secure secret handling
* Auto-scaling systems
* Observability practices
* CI/CD automation

---

# ✅ Maturity Level Achieved

| Capability        | Status |
| ----------------- | ------ |
| Basic Deployment  | ✔      |
| Service Discovery | ✔      |
| Config Management | ✔      |
| Scaling           | ✔      |
| Monitoring        | ✔      |
| CI/CD             | ✔      |

---

# 📌 Final Insight

You have moved from:

```text
Running containers → Operating cloud-native systems
```

This is the foundation of modern backend infrastructure.

---

# Event-Driven Login System (Golang + Redis Pub/Sub)

## 📌 Overview

This project demonstrates **event-based communication** using:

* Golang (Gin API)
* Redis (Pub/Sub)
* Subscriber (event consumer)

Flow:

```
Client → Login API → Redis (Publish) → Subscriber (Consume)
```

---

## ⚙️ Prerequisites

* Go installed (`go version`)
* Docker installed
* PowerShell (Windows)

---

## 🚀 Step 1 — Start Redis

```powershell
docker run -d -p 6379:6379 --name redis-local redis:7-alpine
```

Verify:

```powershell
docker ps
```

---

## 📦 Step 2 — Initialize Project

```powershell
go mod init event-demo
go get github.com/gin-gonic/gin
go get github.com/redis/go-redis/v9
```

---

## 📄 Step 3 — Files

### `main.go`

* Login API
* Publishes event on successful login

### `subscriber.go`

* Listens to Redis channel
* Processes events

---

## ▶️ Step 4 — Run Application

### Terminal 1 — Run Subscriber

```powershell
go run subscriber.go
```

Output:

```
📡 Listening for events...
```

---

### Terminal 2 — Run API

```powershell
go run main.go
```

Output:

```
🚀 API running on http://localhost:8081
```

---

### Terminal 3 — Trigger Event

```powershell
Invoke-RestMethod -Uri "http://localhost:8081/login" `
-Method POST `
-Body (@{username="shubh"; password="1234"} | ConvertTo-Json) `
-ContentType "application/json"
```

---

## 📊 Expected Output

### API Response

```json
{
  "message": "login successful"
}
```

### Subscriber Output

```
🔥 Event received: user_logged_in:shubh
📧 Send login notification email
```

---

## ⚠️ Important Notes

* Subscriber must be running **before** API call
* Redis Pub/Sub does **not persist messages**
* If subscriber is down → event is lost

---

## 🧠 Architecture

```
Client
  ↓
Gin API (Producer)
  ↓
Redis Pub/Sub
  ↓
Subscriber (Consumer)
```

---

## 🧪 Test Multiple Events

```powershell
1..3 | % {
Invoke-RestMethod -Uri "http://localhost:8081/login" `
-Method POST `
-Body (@{username="shubh"; password="1234"} | ConvertTo-Json) `
-ContentType "application/json"
}
```

---

## 🔥 Next Improvements

* Use JSON event payloads
* Add multiple subscribers
* Implement retry mechanism
* Replace Redis Pub/Sub with Kafka/RabbitMQ for durability

---

## ✅ Summary

This project demonstrates:

* Event-driven architecture
* Decoupled services
* Asynchronous communication pattern

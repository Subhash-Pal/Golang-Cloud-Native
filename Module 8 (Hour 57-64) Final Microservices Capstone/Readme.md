```markdown
# 📘 Module 8 (Hour 58–61) Lab Manual  
## Microservices Capstone: API Gateway + Auth + Product + Order Services (Golang + Gin)

---

# 🎯 Objective

This lab demonstrates how to build a **basic microservices architecture using Golang and Gin**, consisting of:

- API Gateway (Reverse Proxy)
- Auth Service (Login mock)
- Product Service (List API)
- Order Service (Create Order API)
- Inter-service communication via HTTP

---

# 🧱 Final Architecture

```

Client
↓
API Gateway (:8080)
├── /auth      → Auth Service (:8001)
├── /products  → Product Service (:8002)
└── /orders    → Order Service (:8003)

```

---

# 📁 Project Structure

```

training_golang/
│
├── api-gateway/
│   ├── main.go
│   ├── config/
│   ├── gateway/
│   ├── go.mod
│
├── auth-service/
│   ├── main.go
│   ├── go.mod
│
├── product-service/
│   ├── main.go
│   ├── go.mod
│
└── order-service/
├── main.go
├── go.mod

````

---

# ⚙️ STEP 1 — API GATEWAY SETUP

## Install dependencies

```bash
go mod init api-gateway
go get github.com/gin-gonic/gin
````

---

## Run Gateway

```bash
go run main.go
```

Gateway runs on:

```
http://localhost:8080
```

---

## Gateway Responsibilities

* Route requests to services
* Reverse proxy HTTP calls
* Maintain service registry (config-based)

---

# 🔐 STEP 2 — AUTH SERVICE (:8001)

## Setup

```bash
go mod init auth-service
go get github.com/gin-gonic/gin
```

---

## main.go

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/login", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Auth Service is working",
		})
	})

	r.Run(":8001")
}
```

---

## Run Service

```bash
go run main.go
```

---

## Test Direct

```bash
curl.exe http://localhost:8001/login
```

---

# 📦 STEP 3 — PRODUCT SERVICE (:8002)

## Setup

```bash
go mod init product-service
go get github.com/gin-gonic/gin
```

---

## main.go

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/list", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"products": []string{"Laptop", "Phone", "Tablet"},
		})
	})

	r.Run(":8002")
}
```

---

## Run Service

```bash
go run main.go
```

---

## Test Direct

```bash
curl.exe http://localhost:8002/list
```

---

# 📦 STEP 4 — ORDER SERVICE (:8003)

## Setup

```bash
go mod init order-service
go get github.com/gin-gonic/gin
```

---

## main.go

```go
package main

import "github.com/gin-gonic/gin"

type OrderRequest struct {
	UserID   string `json:"user_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

func main() {
	r := gin.Default()

	r.POST("/create", func(c *gin.Context) {

		var req OrderRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		c.JSON(200, gin.H{
			"order_id": 1001,
			"user_id":  req.UserID,
			"product":  req.Product,
			"quantity": req.Quantity,
			"status":   "created",
		})
	})

	r.Run(":8003")
}
```

---

## Run Service

```bash
go run main.go
```

---

## Test Direct (PowerShell safe)

```powershell
Invoke-RestMethod -Method POST http://localhost:8003/create `
-ContentType "application/json" `
-Body '{"user_id":"u1","product":"Laptop","quantity":2}'
```

---

# 🌐 STEP 5 — API GATEWAY ROUTING

## config.go

```go
package config

type ServiceConfig struct {
	AuthService    string
	ProductService string
	OrderService   string
}

func LoadConfig() ServiceConfig {
	return ServiceConfig{
		AuthService:    "http://localhost:8001",
		ProductService: "http://localhost:8002",
		OrderService:   "http://localhost:8003",
	}
}
```

---

## router.go

```go
auth := r.Group("/auth")
auth.Any("/*path", ReverseProxy(cfg.AuthService, "/auth"))

products := r.Group("/products")
products.Any("/*path", ReverseProxy(cfg.ProductService, "/products"))

orders := r.Group("/orders")
orders.Any("/*path", ReverseProxy(cfg.OrderService, "/orders"))
```

---

## Run Gateway

```bash
go run main.go
```

---

# 🧪 STEP 6 — FULL SYSTEM TESTING

---

## 🔹 Auth via Gateway

```bash
curl.exe http://localhost:8080/auth/login
```

---

## 🔹 Product via Gateway

```bash
curl.exe http://localhost:8080/products/list
```

---

## 🔹 Order via Gateway

```powershell
Invoke-RestMethod -Method POST http://localhost:8080/orders/create `
-ContentType "application/json" `
-Body '{"user_id":"u1","product":"Laptop","quantity":2}'
```

---

# ⚠️ COMMON ISSUES

## ❌ 404 Not Found

* Route missing in gateway

---

## ❌ 502 Bad Gateway

* Service not running
* Wrong port
* Reverse proxy failure

---

## ❌ Invalid Request Body

* JSON malformed
* PowerShell escaping issue

---

## ❌ curl issues in Windows

Use:

```
Invoke-RestMethod (recommended)
```

---

# 📊 FINAL RESULT

You now have:

* API Gateway pattern ✔
* 3 independent microservices ✔
* Reverse proxy routing ✔
* GET + POST APIs ✔
* JSON-based communication ✔
* Distributed architecture foundation ✔

---

# 🚀 NEXT EXTENSIONS (Hour 62+)



---

```
```

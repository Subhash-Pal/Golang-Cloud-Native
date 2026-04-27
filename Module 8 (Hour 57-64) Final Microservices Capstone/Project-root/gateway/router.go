package gateway

import (
	"api-gateway/config"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg config.ServiceConfig) *gin.Engine {

	r := gin.Default()

	r.SetTrustedProxies(nil)

	r.Use(Logger())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Gateway running"})
	})

	// =========================
	// AUTH SERVICE
	// =========================
	auth := r.Group("/auth")
	{
		auth.Any("/*path", ReverseProxy(cfg.AuthService, "/auth"))
	}

	// =========================
	// PRODUCT SERVICE
	// =========================
	products := r.Group("/products")
	{
		products.Any("/*path", ReverseProxy(cfg.ProductService, "/products"))
	}

	// =========================
	// ORDER SERVICE (future ready)
	// =========================
	orders := r.Group("/orders")
	{
		orders.Any("/*path", ReverseProxy(cfg.OrderService, "/orders"))
	}

	return r
}
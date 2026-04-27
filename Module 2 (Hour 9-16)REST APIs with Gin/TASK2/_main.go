package main

import (
	"time"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		println("Request:", c.Request.Method,
			c.Request.URL.Path,
			"Duration:", end.Sub(start).String())

		// Replace your println line with this for higher precision:
        println("Request:", c.Request.Method,
                c.Request.URL.Path,
                "Duration:", end.Sub(start).Microseconds(), "µs")
	    }
}

func main() {
	r := gin.New()

	r.Use(LoggerMiddleware())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "OK"})
	})

	r.Run(":8080")
}
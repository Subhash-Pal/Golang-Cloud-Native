package routes

import (
	"github.com/gin-gonic/gin"
	"lab3-gin/handlers"
)

func RegisterUserRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/users", handlers.CreateUser)
		api.GET("/users", handlers.GetUsers)
		api.GET("/users/:id", handlers.GetUser)
		api.PUT("/users/:id", handlers.UpdateUser)
		api.DELETE("/users/:id", handlers.DeleteUser)
	}
}


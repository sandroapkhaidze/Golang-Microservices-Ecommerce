package http

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes for the user service
func SetupRouter(userHandler *UserHandler) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", userHandler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
		}
	}

	return router
}

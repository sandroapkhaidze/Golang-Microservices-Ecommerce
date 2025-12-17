package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(orderHandler *OrderHandler) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", orderHandler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder) // POST /api/v1/orders
		}
	}

	return router
}

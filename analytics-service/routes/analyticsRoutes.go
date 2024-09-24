package routes

import (
	"analytics-service/middleware" // Assuming you have existing auth middleware
	"github.com/gin-gonic/gin"
)

func RegisterAnalyticsRoutes(router *gin.Engine) {
	analytics := router.Group("/analytics")
	analytics.Use(middleware.AuthMiddleware()) // Use auth middleware
	{

	}
}

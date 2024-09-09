package routes

import (
	"github.com/gin-gonic/gin"
	"user-service/controllers"
	"user-service/middleware"
)

func AuthRoutes(r *gin.Engine) {
	// Unprotected routes
	//unprotected := r.Group("/api")
	//unprotected.GET("/users", controllers.GetUserss) // List all users (unprotected)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/users/:id", controllers.GetUser) // Get user by ID (protected)

}

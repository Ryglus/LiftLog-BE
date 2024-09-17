package routes

import (
	"auth-service/controllers"
	"auth-service/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	// Public routes
	unprotected := r.Group("/api/auth")
	unprotected.POST("/register", controllers.Register)
	unprotected.POST("/login", controllers.Login)
	unprotected.POST("/logout", controllers.Logout)

	// Protected route
	protected := r.Group("/api/auth")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/refresh-token", controllers.RefreshAccessToken)
}

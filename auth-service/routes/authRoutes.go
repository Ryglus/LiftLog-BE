package routes

import (
	"auth-service/controllers"
	"auth-service/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/api/auth/register", controllers.Register)
	r.POST("/api/auth/login", controllers.Login)
	r.POST("/api/auth/logout", controllers.Logout)

	// Protected route
	protected := r.Group("/api/auth")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/whoami", controllers.WhoAmI)
}

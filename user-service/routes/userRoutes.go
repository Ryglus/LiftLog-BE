package routes

import (
	"github.com/gin-gonic/gin"
	"user-service/controllers"
	"user-service/middleware"
)

func UserRoutes(r *gin.Engine) {
	unprotected := r.Group("/api")
	unprotected.GET("/search", controllers.SearchProfiles)
	unprotected.GET("/profile/:id", controllers.GetUserProfile)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.PUT("/profile", controllers.UpdateProfile)
	protected.GET("/profile/me", controllers.GetUserProfile)
}

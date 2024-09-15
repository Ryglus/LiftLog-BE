package routes

import (
	"github.com/gin-gonic/gin"
	"user-service/controllers"
	"user-service/middleware"
)

func UserRoutes(r *gin.Engine) {
	unprotected := r.Group("/api")

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	unprotected.GET("/profile/:id", controllers.GetUserProfile)
	// Profile management
	protected.PUT("/profile", controllers.UpdateProfile)
	protected.GET("/profile/me", controllers.GetUserProfile)
	protected.POST("/upload-profile-image", controllers.UploadProfileImage)

	// Search
	protected.GET("/search", controllers.SearchProfiles)

}

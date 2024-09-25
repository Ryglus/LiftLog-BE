package routes

import (
	"analytics-service/controllers"
	"analytics-service/middleware" // Assuming you have existing auth middleware
	"github.com/gin-gonic/gin"
)

func RegisterAnalyticsRoutes(router *gin.Engine) {
	protected := router.Group("/api/tracking")
	protected.Use(middleware.AuthMiddleware()) // All routes use the auth middleware
	{
		protected.PUT("/schedule", controllers.PutSchedule) // Create/Update Schedule
		protected.GET("/schedule", controllers.GetSchedule) // Get Schedule

		//protected.PUT("/workout", controllers.PutWorkout(pgRepo))     // Create/Update Workout
		//protected.GET("/workout/:id", controllers.GetWorkout(pgRepo)) // Get Workout by ID

		//protected.PUT("/log-workout", controllers.PutLogWorkout(influxRepo))              // Log workout
		//protected.GET("/log-workout/:exercise_id", controllers.GetLogWorkout(influxRepo)) // Get workout log for an exercise
	}
}

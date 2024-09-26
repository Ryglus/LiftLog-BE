package routes

import (
	"analytics-service/controllers"
	"analytics-service/middleware" // Assuming you have existing auth middleware
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	protected := router.Group("/api/tracking")
	protected.Use(middleware.AuthMiddleware()) // All routes use the auth middleware
	{
		protected.GET("/overview", controllers.GetTracking)
		protected.GET("/schedule", controllers.GetSchedules)

		protected.PUT("/schedule", controllers.PutSchedule)
		protected.PUT("/schedule/workouts", controllers.AssignWorkoutToSchedule)

		protected.GET("/workout", controllers.GetWorkout)

		protected.PUT("/workout", controllers.PutWorkout)
		protected.PUT("/workout/exercise", controllers.AssignExerciseToWorkout)

	}
}

package controllers

import (
	"analytics-service/models"
	"analytics-service/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSchedules retrieves the user's schedule
func GetSchedules(c *gin.Context) {
	JWTuser, _ := c.Get("user")
	userID := uint(JWTuser.(map[string]interface{})["user_id"].(float64))

	// Use the service to get the user's schedule
	schedule, err := services.GetUserSchedules(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

func AssignWorkoutToSchedule(c *gin.Context) {
	JWTuser, _ := c.Get("user")
	userID := uint(JWTuser.(map[string]interface{})["user_id"].(float64))

	var workoutToSchedule models.ScheduleWorkout
	if err := c.ShouldBindJSON(&workoutToSchedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Use the service to associate workout and schedule for each day
	if err := services.AssignWorkoutToSchedule(&workoutToSchedule, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workout assigned to schedule successfully"})
}

func PutSchedule(c *gin.Context) {
	JWTuser, _ := c.Get("user")
	userID := uint(JWTuser.(map[string]interface{})["user_id"].(float64))

	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use the service to save the schedule and check authorization
	if err := services.SaveUserSchedule(&schedule, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule saved successfully"})
}

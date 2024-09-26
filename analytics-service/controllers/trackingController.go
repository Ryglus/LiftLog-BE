package controllers

import (
	"analytics-service/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetTracking retrieves all schedules, workouts, and exercises for a user
func GetTracking(c *gin.Context) {
	// Extract the user ID from JWT
	JWTuser, _ := c.Get("user")
	userID := uint(JWTuser.(map[string]interface{})["user_id"].(float64))

	// Fetch all schedules, workouts, and exercises
	schedules, workouts, exercises, err := services.GetUserTracking(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the three arrays as a single response
	c.JSON(http.StatusOK, gin.H{
		"schedules": schedules,
		"workouts":  workouts,
		"exercises": exercises,
	})
}

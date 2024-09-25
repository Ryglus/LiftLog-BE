package controllers

import (
	"analytics-service/models"
	"analytics-service/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSchedule retrieves the user's schedule
func GetSchedule(c *gin.Context) {
	JWTuser, _ := c.Get("user")
	userID := uint(JWTuser.(map[string]interface{})["user_id"].(float64))

	// Use the service to get the user's schedule
	schedule, err := services.GetUserSchedule(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

// PutSchedule creates or updates the user's schedule
func PutSchedule(c *gin.Context) {
	JWTuser, _ := c.Get("user")
	userID := uint(JWTuser.(map[string]interface{})["user_id"].(float64))

	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the user ID from the JWT claims
	schedule.UserID = userID

	// Use the service to save the schedule
	if err := services.SaveUserSchedule(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule saved successfully"})
}

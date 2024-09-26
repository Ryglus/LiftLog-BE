package controllers

import (
	"analytics-service/models"
	"analytics-service/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// PutWorkout creates or updates the user's workout
func PutWorkout(c *gin.Context) {
	JWTuser, _ := c.Get("user")
	userID := uint(JWTuser.(map[string]interface{})["user_id"].(float64))

	var workout models.Workout
	if err := c.ShouldBindJSON(&workout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workout.UserID = userID
	// Use the service to save the workout and check authorization
	if err := services.SaveUserWorkout(&workout, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workout saved successfully"})
}

// GetWorkout fetches a workout by its ID, ensuring the user has access to it
func GetWorkout(c *gin.Context) {
	JWTuser, _ := c.Get("user")
	userID := uint(JWTuser.(map[string]interface{})["user_id"].(float64))

	// Extract workout ID
	workoutID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout ID"})
		return
	}

	workout, err := services.GetUserWorkout(uint(workoutID), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workout)
}

// AssignExerciseToWorkout assigns an exercise to a workout
func AssignExerciseToWorkout(c *gin.Context) {
	JWTuser, _ := c.Get("user")
	userID := uint(JWTuser.(map[string]interface{})["user_id"].(float64))

	var request struct {
		ExerciseID uint `json:"exercise_id"`
		WorkoutID  uint `json:"workout_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.AssignExerciseToWorkout(request.ExerciseID, request.WorkoutID, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercise assigned to workout successfully"})
}

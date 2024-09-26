package services

import (
	"analytics-service/models"
	"analytics-service/repositories"
	"fmt"
)

// GetUserTracking retrieves all schedules, workouts, and exercises for the user
func GetUserTracking(userID uint) ([]models.Schedule, []models.Workout, []models.Exercise, error) {
	schedules, workouts, exercises, err := repositories.GetTrackingData(userID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to fetch tracking data: %v", err)
	}

	return schedules, workouts, exercises, nil
}

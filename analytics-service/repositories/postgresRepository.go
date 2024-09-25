package repositories

import (
	"analytics-service/database"
	"analytics-service/models"
)

func GetSchedule(userID uint) (*models.Schedule, error) {
	var schedule models.Schedule
	err := database.PostgresDB.Preload("Workouts.Exercises").Where("user_id = ?", userID).First(&schedule).Error
	return &schedule, err
}

// SaveSchedule creates or updates a workout schedule
func SaveSchedule(schedule *models.Schedule) error {
	return database.PostgresDB.Save(schedule).Error
}

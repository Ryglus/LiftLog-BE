package repositories

import (
	"analytics-service/database"
	"analytics-service/models"
)

func GetSchedulesByUserId(userID uint) ([]models.Schedule, error) {
	var schedules []models.Schedule

	// Fetch schedules with active ones first, and preload related data
	err := database.PostgresDB.Preload("Workouts.Exercises").
		Where("user_id = ?", userID).
		Order("active DESC"). // Active schedules come first
		Find(&schedules).     // Use Find instead of First to return all records
		Error

	return schedules, err
}
func GetScheduleByScheduleId(scheduleId uint) (*models.Schedule, error) {
	var schedule models.Schedule
	err := database.PostgresDB.Preload("Workouts.Exercises").Where("id = ?", scheduleId).First(&schedule).Error
	return &schedule, err
}

// SaveSchedule creates or updates a workout schedule
func SaveSchedule(schedule *models.Schedule) error {
	return database.PostgresDB.Save(schedule).Error
}

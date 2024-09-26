package repositories

import (
	"analytics-service/database"
	"analytics-service/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// GetSchedulesByUserId retrieves all schedules for a user from the database
func GetSchedulesByUserId(userID uint) ([]models.Schedule, error) {
	var schedules []models.Schedule
	// Query to fetch schedules for the user, ordering by active schedules first
	err := database.PostgresDB.Preload("Workouts").Where("user_id = ?", userID).
		Order("active DESC").Find(&schedules).Error
	return schedules, err
}

// GetScheduleByID retrieves a specific schedule by its ID
func GetScheduleByID(scheduleID uint) (*models.Schedule, error) {
	var schedule models.Schedule
	err := database.PostgresDB.Preload("Workouts").Where("id = ?", scheduleID).First(&schedule).Error
	return &schedule, err
}

// GetScheduleByScheduleId retrieves a specific schedule by ID (used in service for authorization checks)
func GetScheduleByScheduleId(scheduleID uint) (*models.Schedule, error) {
	var schedule models.Schedule
	err := database.PostgresDB.Where("id = ?", scheduleID).First(&schedule).Error
	return &schedule, err
}

// SaveSchedule saves or updates a schedule in the database
func SaveSchedule(schedule *models.Schedule) error {
	return database.PostgresDB.Save(schedule).Error
}

func AddOrUpdateScheduleWorkout(updateScheduleWorkoutData *models.ScheduleWorkout) error {

	var scheduleWorkout models.ScheduleWorkout
	err := database.PostgresDB.Where("workout_id = ? AND schedule_id = ?", updateScheduleWorkoutData.WorkoutID, updateScheduleWorkoutData.ScheduleID).First(&scheduleWorkout).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to query schedule workout association: %v", err)
	}

	// If it exists, update the days of split; otherwise, create a new association
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create a new association

		if err := database.PostgresDB.Create(&updateScheduleWorkoutData).Error; err != nil {
			return fmt.Errorf("failed to associate workout with schedule: %v", err)
		}

	} else {
		scheduleWorkout.DaysOfSplit = updateScheduleWorkoutData.DaysOfSplit
		// Update the existing association with the new days of the split
		if err := database.PostgresDB.Save(&scheduleWorkout).Error; err != nil {
			return fmt.Errorf("failed to update schedule workout association: %v", err)
		}

	}

	return nil
}

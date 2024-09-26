package repositories

import (
	"analytics-service/database"
	"analytics-service/models"
)

// GetTrackingData fetches all schedules, workouts, and exercises for the user
func GetTrackingData(userID uint) ([]models.Schedule, []models.Workout, []models.Exercise, error) {
	var schedules []models.Schedule
	var workouts []models.Workout
	var exercises []models.Exercise
	var scheduleWorkouts []models.ScheduleWorkout

	// Fetch all schedules for the user
	err := database.PostgresDB.Where("user_id = ?", userID).Find(&schedules).Error
	if err != nil {
		return nil, nil, nil, err
	}

	// Fetch all workouts for the user
	err = database.PostgresDB.Where("user_id = ?", userID).Find(&workouts).Error
	if err != nil {
		return nil, nil, nil, err
	}

	// Fetch all exercises for the user
	err = database.PostgresDB.Where("user_id = ?", userID).Find(&exercises).Error
	if err != nil {
		return nil, nil, nil, err
	}

	// Fetch all schedule-workout assignments for the fetched schedules
	err = database.PostgresDB.Where("schedule_id IN ?", getScheduleIDs(schedules)).Find(&scheduleWorkouts).Error
	if err != nil {
		return nil, nil, nil, err
	}

	// Add the associated ScheduleWorkout data to each schedule
	for i := range schedules {
		schedules[i].ScheduleWorkouts = filterScheduleWorkouts(scheduleWorkouts, schedules[i].ID)
	}

	// Return all three arrays (schedules, workouts, exercises)
	return schedules, workouts, exercises, nil
}

// Helper function to get all schedule IDs from the schedules array
func getScheduleIDs(schedules []models.Schedule) []uint {
	ids := make([]uint, len(schedules))
	for i, schedule := range schedules {
		ids[i] = schedule.ID
	}
	return ids
}

// Helper function to filter schedule workouts for a specific schedule ID
func filterScheduleWorkouts(scheduleWorkouts []models.ScheduleWorkout, scheduleID uint) []models.ScheduleWorkout {
	var filtered []models.ScheduleWorkout
	for _, sw := range scheduleWorkouts {
		if sw.ScheduleID == scheduleID {
			filtered = append(filtered, sw)
		}
	}
	return filtered
}

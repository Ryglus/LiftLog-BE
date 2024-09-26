package services

import (
	"analytics-service/models"
	"analytics-service/repositories"
	"fmt"
	"log"
)

// GetUserSchedules retrieves the user's schedule, handling any business logic
func GetUserSchedules(userID uint) ([]models.Schedule, error) {
	// Example of additional business logic (could be validation, formatting, etc.)
	schedule, err := repositories.GetSchedulesByUserId(userID)
	if err != nil {
		return nil, err
	}

	// You could apply additional business logic here if necessary
	return schedule, nil
}
func AssignWorkoutToSchedule(workoutToSchedule *models.ScheduleWorkout, userID uint) error {
	// Check if the user owns the schedule
	if workoutToSchedule.ScheduleID != 0 {
		existingSchedule, err := repositories.GetScheduleByScheduleId(workoutToSchedule.ScheduleID)
		if err != nil {
			return fmt.Errorf("schedule not found")
		}

		// Check if the user owns the schedule
		if existingSchedule.UserID != userID {
			return fmt.Errorf("unauthorized: you do not own this schedule")
		}
	}
	if workoutToSchedule.ScheduleID != 0 {
		existingSchedule, err := repositories.GetWorkoutByID(workoutToSchedule.WorkoutID)
		if err != nil {
			return fmt.Errorf("workout not found")
		}

		// Check if the user owns the schedule
		if existingSchedule.UserID != userID {
			return fmt.Errorf("unauthorized: you do not own this workout")
		}
	}
	log.Print(workoutToSchedule.DaysOfSplit)
	// Associate the workout with the schedule and days of the split
	return repositories.AddOrUpdateScheduleWorkout(workoutToSchedule)
}

// SaveUserSchedule saves or updates the user's schedule, ensuring the user owns the schedule
func SaveUserSchedule(schedule *models.Schedule, userID uint) error {
	// If the schedule ID exists, check authorization
	if schedule.ID != 0 {
		existingSchedule, err := repositories.GetScheduleByScheduleId(schedule.ID)
		if err != nil {
			return fmt.Errorf("schedule not found")
		}

		// Check if the user owns the schedule
		if existingSchedule.UserID != userID {
			return fmt.Errorf("unauthorized: you do not own this schedule")
		}
	}

	// Set the user ID for new schedules
	if schedule.ID == 0 {
		schedule.UserID = userID
	}

	// Save the schedule
	return repositories.SaveSchedule(schedule)
}

package services

import (
	"analytics-service/models"
	"analytics-service/repositories"
	"fmt"
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

// SaveUserSchedule saves or updates the user's schedule, ensuring the user owns the schedule
func SaveUserSchedule(schedule *models.Schedule, userID uint) error {
	// Fetch the existing schedule to ensure the user owns it (for update scenario)
	existingSchedule, err := repositories.GetScheduleByScheduleId(schedule.ID)
	if err != nil {
		return fmt.Errorf("schedule not found")
	}

	// Ensure the schedule belongs to the current user
	if existingSchedule != nil && existingSchedule.UserID != userID {
		return fmt.Errorf("unauthorized: you do not have permission to edit this schedule")
	}

	// Example of business logic before saving a schedule (validation)
	if schedule.SplitInterval <= 0 {
		return fmt.Errorf("invalid split interval")
	}

	// Call the repository to save the schedule
	return repositories.SaveSchedule(schedule)
}

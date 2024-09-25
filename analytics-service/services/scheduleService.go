package services

import (
	"analytics-service/models"
	"analytics-service/repositories"
	"fmt"
)

// GetUserSchedule retrieves the user's schedule, handling any business logic
func GetUserSchedule(userID uint) (*models.Schedule, error) {
	// Example of additional business logic (could be validation, formatting, etc.)
	schedule, err := repositories.GetSchedule(userID)
	if err != nil {
		return nil, err
	}

	// You could apply additional business logic here if necessary
	return schedule, nil
}

// SaveUserSchedule saves or updates the user's schedule
func SaveUserSchedule(schedule *models.Schedule) error {
	// Example of business logic before saving a schedule
	// For instance, validation like ensuring the split interval is valid
	if schedule.SplitInterval <= 0 {
		return fmt.Errorf("invalid split interval")
	}

	// Call the repository to save the schedule
	return repositories.SaveSchedule(schedule)
}

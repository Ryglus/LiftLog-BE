package services

import (
	"analytics-service/models"
	"analytics-service/repositories"
	"fmt"
)

// GetUserWorkout fetches a specific workout based on the workout ID and user ID
func GetUserWorkout(workoutID uint, userID uint) (*models.Workout, error) {
	workout, err := repositories.GetWorkoutByID(workoutID)
	if err != nil {
		return nil, fmt.Errorf("workout not found")
	}

	// Ensure the workout belongs to a schedule owned by the user
	var schedules []models.Schedule
	if err := repositories.GetSchedulesForWorkout(workoutID, &schedules); err != nil {
		return nil, fmt.Errorf("unable to fetch schedules for this workout")
	}

	// Check if any of the schedules belong to the user
	for _, schedule := range schedules {
		if schedule.UserID == userID {
			return workout, nil
		}
	}

	return nil, fmt.Errorf("unauthorized: you do not have permission to view this workout")
}

// SaveUserWorkout creates or updates the user's workout
func SaveUserWorkout(workout *models.Workout, userID uint) error {
	// If the workout ID exists, check authorization
	if workout.ID != 0 {
		// Ensure the workout belongs to a schedule owned by the user
		var schedules []models.Schedule
		if err := repositories.GetSchedulesForWorkout(workout.ID, &schedules); err != nil {
			return fmt.Errorf("unable to fetch schedules for this workout")
		}

		// Check if any of the schedules belong to the user
		isAuthorized := false
		for _, schedule := range schedules {
			if schedule.UserID == userID {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			return fmt.Errorf("unauthorized: you do not own this workout")
		}
	}

	// Save the workout
	return repositories.SaveWorkout(workout)
}

// AssignExerciseToWorkout associates an exercise with a workout
func AssignExerciseToWorkout(exerciseID uint, workoutID uint, userID uint) error {
	// Ensure the workout belongs to a schedule owned by the user
	workout, err := GetUserWorkout(workoutID, userID)
	if err != nil {
		return err
	}

	// Associate the exercise with the workout
	return repositories.AddExerciseToWorkout(exerciseID, workout.ID)
}

package repositories

import (
	"analytics-service/database"
	"analytics-service/models"
)

// GetSchedulesForWorkout retrieves the schedules associated with a workout
func GetSchedulesForWorkout(workoutID uint, schedules *[]models.Schedule) error {
	return database.PostgresDB.Model(&models.Workout{ID: workoutID}).Association("Schedules").Find(schedules)
}

// GetWorkoutByID retrieves a workout by its ID
func GetWorkoutByID(workoutID uint) (*models.Workout, error) {
	var workout models.Workout
	err := database.PostgresDB.Preload("Exercises").Where("id = ?", workoutID).First(&workout).Error
	return &workout, err
}

// SaveWorkout creates or updates a workout
func SaveWorkout(workout *models.Workout) error {
	return database.PostgresDB.Save(workout).Error
}

// AddExerciseToWorkout associates an exercise with a workout
func AddExerciseToWorkout(exerciseID uint, workoutID uint) error {
	// Use GORM to associate the exercise with the workout
	workout := models.Workout{ID: workoutID}
	exercise := models.Exercise{ID: exerciseID}

	return database.PostgresDB.Model(&workout).Association("Exercises").Append(&exercise)
}

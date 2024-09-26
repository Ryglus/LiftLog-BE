package models

import "time"
import "github.com/lib/pq"

// Schedule Model
type Schedule struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	UserID           uint              `json:"-"`
	Title            string            `json:"title"`
	Active           bool              `json:"active"`
	StartDate        time.Time         `json:"start_date"`
	SplitInterval    int               `json:"split_interval"`
	ScheduleWorkouts []ScheduleWorkout `json:"schedule_workouts"` // Added this field
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

// ScheduleWorkout Join Table for Schedule and Workout
type ScheduleWorkout struct {
	ID          uint          `gorm:"primaryKey"`                      // Auto-increment ID
	ScheduleID  uint          `json:"schedule_id"`                     // Foreign key for schedule
	WorkoutID   uint          `json:"workout_id"`                      // Foreign key for workout
	DaysOfSplit pq.Int64Array `gorm:"type:int[]" json:"days_of_split"` // Array of days this workout is assigned to within the schedule
}
type ScheduleWorkoutResponse struct {
	WorkoutID   uint          `json:"schedule_id"`                     // Foreign key for workout
	DaysOfSplit pq.Int64Array `gorm:"type:int[]" json:"days_of_split"` // Array of days this workout is assigned to within the schedule
}

// Workout Model
type Workout struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `json:"-"`
	WorkoutName string     `json:"workout_name"`
	Color       string     `json:"color"`                                         // Hex code or color name for the workout tile
	Image       string     `json:"image"`                                         // URL or file path to the image
	Exercises   []Exercise `gorm:"many2many:workout_exercises;" json:"exercises"` // Many-to-many for exercises
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Exercise Model
type Exercise struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `json:"-"`
	WorkoutID       uint      `json:"workout_id"`
	ExerciseName    string    `json:"exercise_name"`
	Description     string    `json:"description"`
	DefaultRepRange RepRange  `gorm:"type:json" json:"default_rep_range"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RepRange Model for default rep range
type RepRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

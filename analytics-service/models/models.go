package models

import "time"

//TODO: Add image and colour saving

type Schedule struct {
	ID            uint      `gorm:"primaryKey" json:"id"` // This maps from the "id" field in the request
	UserID        uint      `json:"-"`                    // Do not map this from the request
	Title         string    `json:"title"`
	Active        bool      `json:"active"`
	StartDate     time.Time `json:"start_date"`
	SplitInterval int       `json:"split_interval"` // e.g., 7 for weekly, 8 for 8-day split
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Workouts      []Workout `json:"workouts"`
}

type Workout struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ScheduleID  uint       `json:"schedule_id"`
	WorkoutName string     `json:"workout_name"`
	Exercises   []Exercise `json:"exercises"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Exercise struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	WorkoutID    uint      `json:"workout_id"`
	ExerciseName string    `json:"exercise_name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

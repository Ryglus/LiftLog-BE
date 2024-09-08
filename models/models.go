// models.go
package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username     string `gorm:"unique"`
	PasswordHash string
	Email        string
	Foods        []Food
	Supplements  []Supplement
	Lifts        []Lift
	Workouts     []WorkoutSplit
}

type Food struct {
	gorm.Model
	UserID   uint
	Date     time.Time
	MealType string
	FoodItem string
	Quantity float64
	Calories float64
	Protein  float64
	Carbs    float64
	Fats     float64
}

type Supplement struct {
	gorm.Model
	UserID         uint
	Date           time.Time
	SupplementName string
	Taken          bool
}

type Lift struct {
	gorm.Model
	UserID             uint
	Date               time.Time
	ExerciseID         uint
	Reps               int
	Weight             float64
	Completed          bool
	AdditionalExercise bool
}

type WorkoutSplit struct {
	gorm.Model
	UserID      uint
	Date        time.Time
	DayOfWeek   string
	WorkoutType string
	Exercises   string // JSON or text representation
}

type Exercise struct {
	gorm.Model
	Name     string
	Category string
}

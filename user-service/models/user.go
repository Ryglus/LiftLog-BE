package models

import (
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"unique"`
	FullName     string    `json:"full_name"`
	ProfileImage string    `json:"profile_image"`
	Bio          string    `json:"bio"`
	Location     string    `json:"location"`
	Visibility   string    `json:"visibility" gorm:"default:'public'"` // public, private
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type SearchedUser struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	ProfileImage string `json:"profile_image"`
}

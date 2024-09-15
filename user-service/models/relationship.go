package models

import "time"

type Relationship struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	UserID           uint      `json:"user_id"`           // The user who is following/has sent a request
	TargetUserID     uint      `json:"target_user_id"`    // The user being followed/receiving the request
	RelationshipType string    `json:"relationship_type"` // e.g. "friend", "follow"
	Status           string    `json:"status"`            // e.g. "pending", "accepted", "rejected"
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

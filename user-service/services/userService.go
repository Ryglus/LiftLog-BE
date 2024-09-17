package services

import (
	"user-service/models"
	"user-service/repositories"
)

// SearchProfiles searches for profiles by query (username or bio)
func SearchProfiles(query string) ([]models.User, error) {
	return repositories.SearchProfiles(query)
}

// UpdateUserProfile updates the user's profile based on the provided fields
func UpdateUserProfile(userID uint, username, bio, location, profileImage string) error {
	user, err := repositories.FindUserByID(userID)
	if err != nil {
		return err
	}

	// Update fields only if provided
	if username != "" {
		user.Username = username
	}
	if bio != "" {
		user.Bio = bio
	}
	if location != "" {
		user.Location = location
	}
	if profileImage != "" {
		user.ProfileImage = profileImage
	}

	// Save the updated profile
	return repositories.UpdateUserProfile(user)
}

// GetUserProfile returns a user by ID
func GetUserProfile(id float64) (*models.User, error) {
	return repositories.FindUserByIDStr(id)
}

package repositories

import (
	"user-service/database"
	"user-service/models"
)

// SearchProfiles searches for users based on a query
func SearchProfiles(query string) ([]models.User, error) {
	var users []models.User
	if err := database.DB.Where("username LIKE ?", "%"+query+"%").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// FindUserByID returns a user by ID
func FindUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByIDStr returns a user by ID in string form
func FindUserByIDStr(id float64) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserProfile saves a user's updated profile to the database
func UpdateUserProfile(user *models.User) error {
	return database.DB.Save(user).Error
}

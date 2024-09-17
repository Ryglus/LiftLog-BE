package services

import (
	"auth-service/auth"
	"auth-service/database"
	"auth-service/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrMissingFields        = errors.New("all fields are required")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrFailedToHashPassword = errors.New("failed to hash password")
	ErrFailedToCreateToken  = errors.New("failed to create token")
)

// RegisterUser registers a new user
func RegisterUser(username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return ErrMissingFields
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrFailedToHashPassword
	}

	user := models.User{Username: username, Email: email, PasswordHash: string(hash)}
	if result := database.DB.Create(&user); result.Error != nil {
		return result.Error
	}

	return nil
}

// LoginUser handles user login and returns access and refresh tokens
func LoginUser(username, password string) (string, string, error) {
	if username == "" || password == "" {
		return "", "", ErrMissingFields
	}

	// Find user by username
	var user models.User
	if err := database.DB.Where("email = ?", username).First(&user).Error; err != nil {
		return "", "", ErrUserNotFound
	}

	// Compare hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Generate access token
	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", ErrFailedToCreateToken
	}

	// Generate refresh token
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", ErrFailedToCreateToken
	}

	return accessToken, refreshToken, nil
}
func RefreshTokens(userID uint) (string, error) {

	// Generate new access token
	accessToken, err := auth.GenerateAccessToken(userID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

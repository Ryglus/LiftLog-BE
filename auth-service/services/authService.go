package services

import (
	"LiftLog-BE/auth"
	"LiftLog-BE/database"
	"LiftLog-BE/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrMissingFields        = errors.New("all fields are required")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrFailedToHashPassword = errors.New("failed to hash password")
	ErrFailedToCreateToken  = errors.New("failed to create token")
	ErrExpiredRefreshToken  = errors.New("refresh token expired")
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
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
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
func RefreshTokens(refreshTokenStr string) (string, string, error) {
	// Parse and validate the refresh token
	refreshToken, err := auth.ParseToken(refreshTokenStr)
	if err != nil {
		return "", "", ErrExpiredRefreshToken
	}

	claims, ok := refreshToken.Claims.(*auth.Claims)
	if !ok || !refreshToken.Valid {
		return "", "", ErrExpiredRefreshToken
	}

	// Generate new access token
	accessToken, err := auth.GenerateAccessToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token
	newRefreshToken, err := auth.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

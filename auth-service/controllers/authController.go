package controllers

import (
	"auth-service/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := services.RegisterUser(input.Username, input.Email, input.Password)
	if errors.Is(err, services.ErrMissingFields) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	} else if errors.Is(err, services.ErrFailedToHashPassword) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	accessToken, refreshToken, err := services.LoginUser(input.Email, input.Password)
	if errors.Is(err, services.ErrMissingFields) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	} else if errors.Is(err, services.ErrUserNotFound) || errors.Is(err, services.ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else if errors.Is(err, services.ErrFailedToCreateToken) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
func RefreshAccessToken(c *gin.Context) {
	// Retrieve claims from the context (set by AuthMiddleware)
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	// Assert claims type
	claims, ok := userClaims.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims format"})
		return
	}

	// Ensure it's a refresh token by checking the "version" field
	if claims["version"] != "REFRESH" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
		return
	}

	// Extract the userID from the claims
	userID := uint(claims["user_id"].(float64)) // JWT numbers are float64

	// Generate a new access token
	accessToken, err := services.RefreshTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the new access token
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func Logout(c *gin.Context) {
	// Clear refresh token by setting it to expire immediately
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

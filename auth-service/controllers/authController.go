package controllers

import (
	"auth-service/services"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func WhoAmI(c *gin.Context) {
	// Get userID from context (set by the AuthMiddleware)
	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	log.Print(user)
	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"email":    user.Email,
	})
}

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
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	accessToken, refreshToken, err := services.LoginUser(input.Username, input.Password)
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

func Logout(c *gin.Context) {
	// Clear refresh token by setting it to expire immediately
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

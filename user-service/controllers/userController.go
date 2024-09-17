package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"user-service/database"
	"user-service/models"
	"user-service/services"

	"github.com/gin-gonic/gin"
)

// SearchProfiles allows users to search for other profiles
func SearchProfiles(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	users, err := services.SearchProfiles(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search profiles"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func UpdateProfile(c *gin.Context) {
	JWTuser, _ := c.Get("user")
	userID := JWTuser.(map[string]interface{})["user_id"].(float64)

	profileImage, err := services.HandleFileUpload(c, "profile_image")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload profile image"})
		return
	}

	err = services.UpdateUserProfile(uint(userID), c.PostForm("username"), c.PostForm("bio"), c.PostForm("location"), profileImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetUserProfile(c *gin.Context) {
	JWTuser, _ := c.Get("user")
	userID := JWTuser.(map[string]interface{})["user_id"].(float64)
	user, err := services.GetUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UploadProfileImage(c *gin.Context) {
	// Extract the user ID from the context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the uploaded file from the request
	file, err := c.FormFile("profile_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Set the path where the file will be saved
	uploadDir := "./uploads"
	fileName := fmt.Sprintf("user_%v_%s", userID, filepath.Base(file.Filename)) // e.g., user_1_profile.jpg
	filePath := filepath.Join(uploadDir, fileName)

	// Ensure the uploads directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Save the file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Update the user's profile image in the database
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.ProfileImage = filePath

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	// Respond with the updated profile image path
	c.JSON(http.StatusOK, gin.H{
		"message":     "Profile image uploaded successfully",
		"profile_url": filePath,
	})
}

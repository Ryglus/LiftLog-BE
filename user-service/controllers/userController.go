package controllers

import (
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"user-service/database"
	"user-service/models"

	"github.com/gin-gonic/gin"
)

// SearchProfiles allows users to search for other profiles
func SearchProfiles(c *gin.Context) {
	query := c.Query("query")

	// Ensure the query is not empty
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var users []models.User
	if err := database.DB.Where("username LIKE ? OR bio LIKE ?", "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search profiles"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func UpdateProfile(c *gin.Context) {
	// Extract the user ID from the JWT token
	JWTuser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claimsMap, ok := JWTuser.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	userID, _ := claimsMap["user_id"].(float64)

	// Check if the user already exists in the database
	var user models.User
	if err := database.DB.Where("id = ?", uint(userID)).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User doesn't exist, create a new user profile
			createNewProfile(c, uint(userID))
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	// User exists, update the profile fields
	updateProfile(c, &user)
}

// createNewProfile handles profile creation when the user does not exist
func createNewProfile(c *gin.Context, userID uint) {
	user := models.User{
		ID:       userID,
		Username: c.PostForm("username"),
		Bio:      c.PostForm("bio"),
		Location: c.PostForm("location"),
	}

	// Upload profile image if provided
	filePath, err := handleFileUpload(c, "profile_image")
	if err != nil {
		if err != http.ErrMissingFile { // Ignore missing files, it's optional
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload profile image"})
			return
		}
	} else {
		user.ProfileImage = filePath
	}

	// Create the new user profile in the database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Profile created successfully"})
}

// updateProfile updates the existing profile with only the provided fields
func updateProfile(c *gin.Context, user *models.User) {
	// Update fields only if they are provided
	if username := c.PostForm("username"); username != "" {
		user.Username = username
	}
	if bio := c.PostForm("bio"); bio != "" {
		user.Bio = bio
	}
	if location := c.PostForm("location"); location != "" {
		user.Location = location
	}

	// Upload profile image if provided
	filePath, err := handleFileUpload(c, "profile_image")
	if err != nil {
		if err != http.ErrMissingFile { // Ignore missing files, it's optional
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload profile image"})
			return
		}
	} else {
		user.ProfileImage = filePath
	}

	// Save the updated profile in the database
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// handleFileUpload processes file upload and saves the file, returning the file path
func handleFileUpload(c *gin.Context, formKey string) (string, error) {
	file, header, err := c.Request.FormFile(formKey)
	if err != nil {
		return "", err // Will be handled by the caller
	}

	filePath := fmt.Sprintf("./uploads/%s", header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// GetUserProfile returns the user's profile by ID (visible profiles only)
func GetUserProfile(c *gin.Context) {
	// Extract the user ID from the request parameters
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Only return public profiles or the user's own profile
	requestingUserID, _ := c.Get("user_id") // AuthMiddleware sets this
	if user.Visibility == "private" && requestingUserID != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "This profile is private"})
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

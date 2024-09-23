package controllers

import (
	"log"
	"net/http"
	"user-service/repositories"
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

	// Find the user by ID to get the current profile image path
	user, err := repositories.FindUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	profileImage, err := services.HandleFileUpload(c, "profile_image", "pfp")
	if err != nil && err != http.ErrMissingFile {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload profile image"})
		return
	}

	// If a new image is uploaded, delete the old one
	if profileImage != "" && user.ProfileImage != "" {
		if err := services.DeleteOldFile(user.ProfileImage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old profile image"})
			return
		}
		user.ProfileImage = profileImage // Update with new image path
	}
	//TODO: non posted updates not updated
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

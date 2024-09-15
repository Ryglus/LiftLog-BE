package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/database"
	"user-service/models"
)

// SendFriendRequest allows users to send a friend request or follow
func SendFriendRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var relationship models.Relationship
	if err := c.BindJSON(&relationship); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Ensure the target user ID is set
	if relationship.TargetUserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target user is required"})
		return
	}

	// Create a new friend/follow request
	relationship.UserID = userID.(uint)
	relationship.Status = "pending" // By default, the request is pending

	if err := database.DB.Create(&relationship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request sent successfully"})
}

// ManageFriendRequest allows users to accept or reject friend requests
func ManageFriendRequest(c *gin.Context) {
	relationshipID := c.Param("id")
	var relationship models.Relationship
	if err := database.DB.First(&relationship, relationshipID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	// Accept or reject the request
	status := c.Query("status") // "accepted" or "rejected"
	if status == "accepted" || status == "rejected" {
		relationship.Status = status
		if err := database.DB.Save(&relationship).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Request " + status})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
	}
}

// GetFriends retrieves all accepted friends for a user
func GetFriends(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var friends []models.User
	database.DB.Raw(`
		SELECT users.* 
		FROM users 
		JOIN relationships 
		ON users.id = relationships.target_user_id 
		WHERE relationships.user_id = ? 
		AND relationships.status = 'accepted'
	`, userID).Scan(&friends)

	c.JSON(http.StatusOK, friends)
}

package controllers

import (
	"analytics-service/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PutExercise(c *gin.Context) {
	//JWTuser, _ := c.Get("user")
	//userID := uint(JWTuser.(map[string]interface{})["user_id"].(float64))

	var exercise models.Exercise
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercise saved successfully"})
}

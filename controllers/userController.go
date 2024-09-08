package controllers

import (
	"LiftLog-BE/database"
	"LiftLog-BE/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWorkoutSplit(c *gin.Context) {
	date := c.Param("date")
	var split models.WorkoutSplit
	if err := database.DB.Where("date = ?", date).First(&split).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout split not found"})
		return
	}
	c.JSON(http.StatusOK, split)
}

func AddLift(c *gin.Context) {
	var lift models.Lift
	if err := c.ShouldBindJSON(&lift); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&lift)
	c.JSON(http.StatusOK, lift)
}

// More handlers for Food, Supplements, and Exercises can be added similarly

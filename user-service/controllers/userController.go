package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetUser atm gets userJwt info that's just for testing tho
func GetUser(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, user)
}

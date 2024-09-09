package middleware

import (
	"LiftLog-BE/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT and adds user info to the context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Remove Bearer if present
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// Validate token and retrieve claims
		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Attach user ID to context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

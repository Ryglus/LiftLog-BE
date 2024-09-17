package middleware

import (
	"auth-service/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
		claims, err := auth.ValidateToken(strings.TrimPrefix(tokenStr, "Bearer "))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

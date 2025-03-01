package middleware

import (
	"net/http"
	"user_crud/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware - Protect routes with JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || !utils.ValidateJWT(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

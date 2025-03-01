package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 5) // 1 request per second, burst up to 5

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}

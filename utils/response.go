package utils

import "github.com/gin-gonic/gin"

func SendSuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(200, gin.H{"success": true, "message": message, "data": data})
}

func SendErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"success": false, "error": message})
}

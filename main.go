package main

import (
	"user_crud/config"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize database
	config.ConnectDB()

	// Register Routes
	UserRoutes(r)

	// Start Server
	r.Run(":8080")
}

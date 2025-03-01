package cmd

import (
	"user_crud/config"
	"user_crud/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize database
	config.ConnectDB()

	// Register Routes
	routes.UserRoutes(r)

	// Start Server
	r.Run(":8080")
}

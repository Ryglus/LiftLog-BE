package main

import (
	"analytics-service/database"
	"analytics-service/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize InfluxDB
	database.InitInfluxDB()
	defer database.CloseInfluxDB()

	// Setup Gin router
	router := gin.Default()

	// Register analytics routes
	routes.RegisterAnalyticsRoutes(router)

	// Start the server
	err := router.Run(":8082")
	if err != nil {
		return
	} // Start the server on port 8082
}

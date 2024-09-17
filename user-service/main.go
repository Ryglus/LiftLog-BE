package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"user-service/database"
	"user-service/routes"
)

func main() {
	// Connect to the database
	database.ConnectDatabase()

	// Initialize Gin router
	r := gin.Default()

	// Enable CORS for all origins or specific origins
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Add your frontend origin here
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Static("/uploads", "./uploads")
	// Register routes
	routes.UserRoutes(r)

	// Run server
	err := r.Run(":8081")
	if err != nil {
		return
	}
}

package main

import (
	"LiftLog-BE/database"
	"LiftLog-BE/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// Register routes
	routes.AuthRoutes(r)

	// Run server
	err := r.Run(":8080")
	if err != nil {
		return
	}
}

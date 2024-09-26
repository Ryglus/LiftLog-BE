package main

import (
	"analytics-service/database"
	"analytics-service/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize PostgreSQL
	database.InitPostgresDB()

	// Initialize InfluxDB
	database.InitInfluxDB()
	defer database.CloseInfluxDB()

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
	routes.RegisterRoutes(r)

	// Run server
	err := r.Run(":8082")
	if err != nil {
		return
	}
}

/*
schedule
should contain day of a week/month or whatever split ratio is set (like 8 day split ext) so it should just continue over that from the start, and wintin those days it will have a workout

workout will be collection of exersizes

and then the user will be able to add how much they lifted and sets ext to each workout, as it will be presented for them what they gonna do that specific day.

so thats the whole workout plan scheme

then i wanna have simular things but with suplmenets and and food (suplements can be set to whatever schedule, but food has not fixed timings or days, it can be added whenever.

from that after i would like to generate some analytics, weight lifted, how did they improove ext, but thats for after we have all made up.
*/

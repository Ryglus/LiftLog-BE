package database

import (
	"analytics-service/models"
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	PostgresDB   *gorm.DB
	InfluxClient influxdb2.Client
	InfluxOrg    string
	InfluxBucket string
)

func InitPostgresDB() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve PostgreSQL configuration from environment variables
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		host, port, user, password, dbname)

	// Connect to PostgreSQL using GORM
	PostgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Auto-migrate your PostgreSQL models
	err = PostgresDB.AutoMigrate(&models.Schedule{}, &models.Workout{}, &models.Exercise{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	fmt.Println("Connected to PostgreSQL and migrated models successfully!")
}

func InitInfluxDB() {
	// Load environment variables if not already loaded
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, relying on environment variables")
	}

	// Retrieve InfluxDB configuration from environment variables
	url := os.Getenv("INFLUXDB_URL")
	token := os.Getenv("INFLUXDB_TOKEN")
	InfluxOrg = os.Getenv("INFLUXDB_ORG")
	InfluxBucket = os.Getenv("INFLUXDB_BUCKET")

	if url == "" || token == "" || InfluxOrg == "" || InfluxBucket == "" {
		log.Fatalf("InfluxDB configuration missing in .env file")
	}

	// Initialize InfluxDB client
	InfluxClient = influxdb2.NewClient(url, token)

	// Test the connection
	health, err := InfluxClient.Health(context.Background())
	if err != nil || health.Status != "pass" {
		log.Fatalf("InfluxDB connection failed: %v", err)
	}

	fmt.Println("Connected to InfluxDB successfully!")
}

func CloseInfluxDB() {
	InfluxClient.Close()
}

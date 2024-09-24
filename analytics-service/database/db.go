package database

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var InfluxClient influxdb2.Client
var Org = "liftlog"
var Bucket = "health_analytics"

func InitInfluxDB() {
	// InfluxDB configuration
	url := "http://localhost:8086"
	token := "your-token-here" // Replace with your InfluxDB token

	// Create a new InfluxDB client
	InfluxClient = influxdb2.NewClient(url, token)
	fmt.Println("Connected to InfluxDB!")
}

// Close the client
func CloseInfluxDB() {
	InfluxClient.Close()
}

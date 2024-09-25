package repositories

import (
	"context"
	"fmt"
	"time"

	"analytics-service/database"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// LogWorkout logs or updates a workout in InfluxDB
func LogWorkout(userID, exerciseID string, weight float64, sets, reps int) error {
	writeAPI := database.InfluxClient.WriteAPIBlocking(database.InfluxOrg, database.InfluxBucket)

	point := influxdb2.NewPointWithMeasurement("workout_logs").
		AddTag("user_id", userID).
		AddTag("exercise_id", exerciseID).
		AddField("weight", weight).
		AddField("sets", sets).
		AddField("reps", reps).
		SetTime(time.Now())

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		fmt.Printf("Error logging workout to InfluxDB: %v\n", err)
		return err
	}

	fmt.Println("Workout logged successfully.")
	return nil
}

// GetWorkoutLogs retrieves logs for a specific exercise and user
func GetWorkoutLogs(userID, exerciseID string, start, stop time.Time) (string, error) {
	queryAPI := database.InfluxClient.QueryAPI(database.InfluxOrg)

	query := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "workout_logs")
			|> filter(fn: (r) => r["user_id"] == "%s")
			|> filter(fn: (r) => r["exercise_id"] == "%s")
	`, database.InfluxBucket, start.Format(time.RFC3339), stop.Format(time.RFC3339), userID, exerciseID)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Error querying workout logs from InfluxDB: %v\n", err)
		return "", err
	}

	var logs string
	for result.Next() {
		logs += fmt.Sprintf("Time: %s, Weight: %v, Sets: %v, Reps: %v\n",
			result.Record().Time().Format(time.RFC3339),
			result.Record().ValueByKey("weight"),
			result.Record().ValueByKey("sets"),
			result.Record().ValueByKey("reps"),
		)
	}

	if result.Err() != nil {
		fmt.Printf("Error in query result: %v\n", result.Err())
		return "", result.Err()
	}

	return logs, nil
}

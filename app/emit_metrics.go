package main

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	token := "6_VwzCyty_-XoSk9JaNwsyZySo4094rkjsZqI3kSw5nYVvj8El3L4brb2BFhC3LID1C9id2IXcc7QvsZ-6YoVg=="
	// Set up a connection to InfluxDB
	client := influxdb2.NewClient("http://localhost:8086", token)

	fmt.Println("Emitting metrics to InfluxDB")
	// get non-blocking write client
	writeAPI := client.WriteAPI("private", "default")

	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 18.5, "max": 30},
		time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// create point using fluent style
	p = influxdb2.NewPointWithMeasurement("stat").
		AddTag("unit", "temperature").
		AddField("avg", 18.2).
		AddField("max", 30).
		SetTime(time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()

	// Close the connection to InfluxDB
	defer client.Close()

	fmt.Println("Metrics emitted successfully")
}

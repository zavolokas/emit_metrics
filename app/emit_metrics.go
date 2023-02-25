package main

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func writeMetrics(client influxdb2.Client) {
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
}

func writeAnnot(client influxdb2.Client) {
	fmt.Println("Emitting annotations to InfluxDB")
	// get non-blocking write client
	writeAPI := client.WriteAPI("private", "default")

	p := influxdb2.NewPointWithMeasurement("events").
		AddTag("type", "accident").
		AddTag("severity", "mid").
		AddField("title", "accident_1").
		AddField("text", "something happened and we don't know what").
		SetTime(time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()
}

func main() {
	token := "6_VwzCyty_-XoSk9JaNwsyZySo4094rkjsZqI3kSw5nYVvj8El3L4brb2BFhC3LID1C9id2IXcc7QvsZ-6YoVg=="
	// Set up a connection to InfluxDB
	client := influxdb2.NewClient("http://localhost:8086", token)

	// writeMetrics(client)
	writeAnnot(client)

	// Close the connection to InfluxDB
	defer client.Close()

	fmt.Println("Metrics emitted successfully")
}

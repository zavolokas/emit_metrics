package testmetrics

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	log "github.com/sirupsen/logrus"
)

func WriteMetrics(client influxdb2.Client) {
	log.Info("Emitting metrics to InfluxDB")
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
	log.Info("Metrics emitted successfully")
}

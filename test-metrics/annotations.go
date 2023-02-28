package testmetrics

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	log "github.com/sirupsen/logrus"
)

func WriteAnnot(client influxdb2.Client) {
	log.Info("Emitting annotations to InfluxDB")
	// get non-blocking write client
	writeAPI := client.WriteAPI("private", "default")

	timeEnd := time.Now()
	startTime := timeEnd.Add(-(1 * time.Minute))

	p := influxdb2.NewPointWithMeasurement("events").
		AddTag("type", "testing").
		AddTag("severity", "mid").
		AddField("title", "range").
		// AddField("timeEnd", time.Now().Add(-(5*time.Minute))).
		AddField("text", "working on grafana integration").
		SetTime(startTime)
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()
	log.Info("Annotations emitted successfully")
}

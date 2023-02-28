package main

import (
	"os"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	nestmetrics "github.com/zavolokas/emit_metrics/nest-metrics"
)

func main() {
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	err := godotenv.Load()
	if err != nil {
		log.WithError(err).Fatal("Error loading .env file")
		return
	}

	token := os.Getenv("INFLUXDB_TOKEN")
	// Set up a connection to InfluxDB
	client := influxdb2.NewClient("http://localhost:8086", token)
	defer client.Close()

	nestConfig := nestmetrics.Config{
		ClientID:     os.Getenv("NEST_CLIENT_ID"),
		ClientSecret: os.Getenv("NEST_CLIENT_SECRET"),
		ProjectID:    os.Getenv("NEST_PROJECT_ID"),
	}

	// testmetrics.WriteMetrics(client)
	// testmetrics.WriteAnnot(client)

	nestFreqSec, err := strconv.Atoi(os.Getenv("NEST_FREQUENCY_SEC"))
	if err != nil {
		log.WithError(err).Fatal("failed to get nest frequency", err)
		return
	}
	for {
		err := nestmetrics.EmitNestMetrics(nestConfig, client)
		if err != nil {
			log.WithError(err).Warn("failed to emit nest metrics", err)
		}
		time.Sleep(time.Duration(nestFreqSec) * time.Second)
	}
}

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func emitNestMetrics(cfg NestClientConfig, client influxdb2.Client) error {
	log.Info("Getting Nest Device...")
	nestDevice, err := getNestDevice(cfg)
	if err != nil {
		log.WithError(err).Warn("failed to get nest device")
		return err
	}

	log.Info("Emitting Thermostat metrics to InfluxDB...")

	writeAPI := client.WriteAPI("private", "default")

	// create point using fluent style
	p := influxdb2.NewPointWithMeasurement("nest_thermostat").
		// AddTag("unit", "temperature").
		AddField("actual", nestDevice.TemperatureActual).
		AddField("set", nestDevice.TemperatureSet).
		AddField("mode", nestDevice.Mode).
		AddField("status", nestDevice.Status).
		AddField("humid", nestDevice.HumidityPercent).
		SetTime(time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()

	log.Info("Metrics emitted successfully")
	return nil
}

type NestClientConfig struct {
	ClientID     string
	ClientSecret string
	ProjectID    string
}

type NestDevice struct {
	TemperatureActual float64
	TemperatureSet    float64
	Mode              int
	Status            int
	HumidityPercent   float64
}

func getNestDevice(cfg NestClientConfig) (NestDevice, error) {
	return NestDevice{}, fmt.Errorf("not implemented")
}

func writeMetrics(client influxdb2.Client) {
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

func writeAnnot(client influxdb2.Client) {
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

	nestConfig := NestClientConfig{
		ClientID:     os.Getenv("NEST_CLIENT_ID"),
		ClientSecret: os.Getenv("NEST_CLIENT_SECRET"),
		ProjectID:    os.Getenv("NEST_PROJECT_ID"),
	}

	// Will not emit metrics if SW is not set to asda
	if os.Getenv("SW") == "asda" {
		writeMetrics(client)
		writeAnnot(client)
	}
	nestFreqSec, err := strconv.Atoi(os.Getenv("NEST_FREQUENCY_SEC"))
	if err != nil {
		log.WithError(err).Fatal("failed to get nest frequency", err)
		return
	}
	for {
		err := emitNestMetrics(nestConfig, client)
		if err != nil {
			log.WithError(err).Warn("failed to emit nest metrics", err)
		}
		time.Sleep(time.Duration(nestFreqSec) * time.Second)
	}
}

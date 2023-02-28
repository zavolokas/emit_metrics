package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	nestmetrics "github.com/zavolokas/emit_metrics/nest-metrics"
)

const (
	exitCodeErr       = 1
	exitCodeInterrupt = 2
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

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		os.Exit(exitCodeInterrupt)
	}()

	token := os.Getenv("INFLUXDB_TOKEN")
	// Set up a connection to InfluxDB
	client := influxdb2.NewClient("http://localhost:8086", token)
	defer client.Close()

	nestFreqSec, err := strconv.Atoi(os.Getenv("NEST_FREQUENCY_SEC"))
	if err != nil {
		log.WithError(err).Fatal("failed to get nest frequency", err)
		return
	}

	nestConfig := nestmetrics.Config{
		ClientID:     os.Getenv("NEST_CLIENT_ID"),
		ClientSecret: os.Getenv("NEST_CLIENT_SECRET"),
		ProjectID:    os.Getenv("NEST_PROJECT_ID"),
		EmitFreqSec:  nestFreqSec,
	}

	if err := run(ctx, nestConfig, client); err != nil {
		log.WithError(err).Fatal("error running program")
		os.Exit(exitCodeErr)
	}
}

func run(ctx context.Context, nestConfig nestmetrics.Config, client influxdb2.Client) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:

			// testmetrics.WriteMetrics(client)
			// testmetrics.WriteAnnot(client)

			err := nestmetrics.EmitNestMetrics(nestConfig, client)
			if err != nil {
				log.WithError(err).Warn("failed to emit nest metrics", err)
			}
			time.Sleep(time.Duration(nestConfig.EmitFreqSec) * time.Second)
		}
	}
}

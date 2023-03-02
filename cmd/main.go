package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

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

	projectID := os.Getenv("NEST_PROJECT_ID")
	authURL := strings.Replace(os.Getenv("NEST_AUTH_URL_PATTERN"), "%PROJECT_ID%", projectID, 1)

	nestConfig := nestmetrics.Config{
		ClientID:     os.Getenv("NEST_CLIENT_ID"),
		ClientSecret: os.Getenv("NEST_CLIENT_SECRET"),
		ProjectID:    projectID,
		AuthURL:      authURL,
		TokenURL:     os.Getenv("NEST_TOKEN_URL"),
		RedirectURL:  os.Getenv("NEST_REDIRECT_URL"),
		Scopes:       strings.Split(os.Getenv("NEST_SCOPES"), ","),
		RefreshToken: os.Getenv("NEST_REFRESH_TOKEN"),
		EmitFreqSec:  nestFreqSec,
	}

	nestClient := nestmetrics.NewNestClient(nestConfig, client)

	if err := run(ctx, nestConfig, nestClient, client); err != nil {
		log.WithError(err).Fatal("error running program")
		os.Exit(exitCodeErr)
	}
}

func run(ctx context.Context, nestConfig nestmetrics.Config, nestClient nestmetrics.I, client influxdb2.Client) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:

			// testmetrics.WriteMetrics(client)
			// testmetrics.WriteAnnot(client)

			err := nestClient.EmitNestMetrics()
			if err != nil {
				log.WithError(err).Warn("failed to emit nest metrics", err)
			}
			time.Sleep(time.Duration(nestConfig.EmitFreqSec) * time.Second)
		}
	}
}

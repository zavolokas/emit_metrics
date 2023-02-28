package nestmetrics

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	log "github.com/sirupsen/logrus"
)

func EmitNestMetrics(cfg Config, client influxdb2.Client) error {
	log.Info("getting nest device data point...")
	dp, err := getDataPoint(cfg)
	if err != nil {
		log.WithError(err).Warn("failed to get nest device data")
		return err
	}

	log.Info("emitting thermostat metrics to influxdb...")

	writeAPI := client.WriteAPI("private", "default")

	// create point using fluent style
	p := influxdb2.NewPointWithMeasurement("nest_thermostat").
		// AddTag("unit", "temperature").
		AddField("actual", dp.TemperatureActual).
		AddField("set", dp.TemperatureSet).
		AddField("mode", dp.Mode).
		AddField("status", dp.Status).
		AddField("humid", dp.HumidityPercent).
		SetTime(time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()

	log.Info("metrics emitted successfully")
	return nil
}

func getDataPoint(cfg Config) (DataPoint, error) {
	return DataPoint{}, fmt.Errorf("not implemented")
}

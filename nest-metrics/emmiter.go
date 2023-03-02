package nestmetrics

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type nestClient struct {
	httpClient   *http.Client
	influxClient influxdb2.Client
	projectID    string
}

type I interface {
	EmitNestMetrics() error
}

func NewNestClient(cfg Config, influxClient influxdb2.Client) I {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       cfg.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.AuthURL,
			TokenURL: cfg.TokenURL,
		},
		RedirectURL: cfg.RedirectURL,
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Visit the URL for the auth dialog: %v\n\n Paste the AuthZ code here:\n", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)

	return &nestClient{
		httpClient:   client,
		influxClient: influxClient,
		projectID:    cfg.ProjectID,
	}
}

func (nc *nestClient) EmitNestMetrics() error {
	log.Info("getting nest device data point...")
	dp, err := getDataPoint(nc.httpClient, nc.projectID)
	if err != nil {
		log.WithError(err).Warn("failed to get nest device data")
		return err
	}

	log.Info("emitting thermostat metrics to influxdb...")

	writeAPI := nc.influxClient.WriteAPI("private", "default")

	// create point using fluent style
	p := influxdb2.NewPointWithMeasurement("nest_thermostat").SetTime(time.Now())
	// AddTag("unit", "temperature").
	if dp.TemperatureActual != nil {
		p = p.AddField("actual", *dp.TemperatureActual)
	}
	if dp.TemperatureSet != nil {
		p = p.AddField("set", *dp.TemperatureSet)
	}
	if dp.HumidityPercent != nil {
		p = p.AddField("humid", *dp.HumidityPercent)
	}
	if dp.Mode != nil {
		p = p.AddField("mode", *dp.Mode)
	}

	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()

	log.Info("metrics emitted successfully")
	return nil
}

func getDataPoint(httpClient *http.Client, projectID string) (DataPoint, error) {
	url := fmt.Sprintf("https://smartdevicemanagement.googleapis.com/v1/enterprises/%s/devices", projectID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithError(err).Warn("failed to create nest device request")
		return DataPoint{}, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.WithError(err).Warn("failed to get nest device data")
		return DataPoint{}, err
	}

	if resp.StatusCode != http.StatusOK {
		log.WithField("status_code", resp.StatusCode).Warn("failed to get nest device data")
		return DataPoint{}, fmt.Errorf("failed to get nest device data: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Warn("failed to read nest device response body")
		return DataPoint{}, err
	}
	//Convert the body to type string
	sb := string(body)
	log.Debug(sb)

	devicesResponse := NestDeviceResponse{}
	err = json.Unmarshal(body, &devicesResponse)
	if err != nil {
		log.WithError(err).Warn("failed to unmarshal nest device response")
		return DataPoint{}, err
	}

	if devicesResponse.Devices == nil || len(devicesResponse.Devices) == 0 {
		log.Warn("no devices found")
		return DataPoint{}, err
	}

	device := devicesResponse.Devices[0]

	dp := DataPoint{}

	if device.Traits["sdm.devices.traits.Temperature"]["ambientTemperatureCelsius"] != nil {
		tmp := device.Traits["sdm.devices.traits.Temperature"]["ambientTemperatureCelsius"].(float64)
		dp.TemperatureActual = &tmp
	} else {
		log.Warn("no ambient temperature found")
	}
	if device.Traits["sdm.devices.traits.ThermostatTemperatureSetpoint"]["heatCelsius"] != nil {
		tmp := device.Traits["sdm.devices.traits.ThermostatTemperatureSetpoint"]["heatCelsius"].(float64)
		dp.TemperatureSet = &tmp
	} else {
		log.Warn("no set temperature found")
	}
	if device.Traits["sdm.devices.traits.Humidity"]["ambientHumidityPercent"] != nil {
		tmp := device.Traits["sdm.devices.traits.Humidity"]["ambientHumidityPercent"].(float64)
		dp.HumidityPercent = &tmp
	} else {
		log.Warn("no humidity found")
	}
	if device.Traits["sdm.devices.traits.ThermostatHvac"]["status"] != nil {
		var heat int
		switch device.Traits["sdm.devices.traits.ThermostatHvac"]["status"] {
		case "HEATING":
			heat = 1
		case "COOLING":
			heat = -1
		default:
			heat = 0
		}
		dp.Mode = &heat
	} else {
		log.Warn("no HVAC Status found")
	}

	log.
		WithFields(
			log.Fields{
				"actual_temperature": *dp.TemperatureActual,
				"set_temperature":    *dp.TemperatureSet,
				"humidity":           *dp.HumidityPercent,
				"mode":               *dp.Mode,
			}).
		Info("got nest device data point")

	return dp, nil
}

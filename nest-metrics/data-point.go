package nestmetrics

type DataPoint struct {
	TemperatureActual *float64
	TemperatureSet    *float64
	Mode              *int
	HumidityPercent   *float64
}

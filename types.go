package smhi

import (
	"time"
)

// PrecipitationCategory constants.
const (
	NoPrecipitation PrecipitationCategory = iota
	Snow
	SnowAndRain
	Rain
	Drizzle
	FreezingRain
	FreezingDrizzle
)

// WeatherSymbol constants.
const (
	ClearSky WeatherSymbol = iota + 1
	NearlyClearSky
	VariableCloudiness
	HalfclearSky
	CloudySky
	Overcast
	Fog
	LightRainShowers
	ModerateRainShowers
	HeavyRainShowers
	Thunderstorm
	LightSleetShowers
	ModerateSleetShowers
	HeavySleetShowers
	LightSnowShowers
	ModerateSnowShowers
	HeavySnowShowers
	LightRain
	ModerateRain
	HeavyRain
	Thunder
	LightSleet
	ModerateSleet
	HeavySleet
	LightSnowfall
	ModerateSnowfall
	HeavySnowfall
)

type PrecipitationCategory uint8

type Coordinate []float64

type WeatherSymbol uint8

type Geometry struct {
	Type        string
	Coordinates Coordinate
}

// PointForecastAPI defines the data structure that is returned by the SMHI
// point forecast API
type PointForecastAPI struct {
	ApprovedTime  string                    `json:"approvedTime"`
	CreatedTime   string                    `json:"createdTime"`
	ReferenceTime string                    `json:"referenceTime"`
	Geometry      Geometry                  `json:"geometry"`
	TimeSeries    []PointForecastTimeSeries `json:"timeSeries"`
}

type PointForecastTimeSeries struct {
	Time                        string            `json:"time"`
	IntervalParametersStartTime string            `json:"intervalParametersStartTime"`
	Data                        PointForecastData `json:"data"`
}

type PointForecastData struct {
	AirTemperature                        float64       `json:"air_temperature"`
	WindFromDirection                     uint16        `json:"wind_from_direction"`
	WindSpeed                             float64       `json:"wind_speed"`
	WindSpeedOfGust                       float64       `json:"wind_speed_of_gust"`
	RelativeHumidity                      uint8         `json:"relative_humidity"`
	AirPressureAtMeanSeaLevel             float64       `json:"air_pressure_at_mean_sea_level"`
	VisibilityInAir                       float64       `json:"visibility_in_air"`
	ThunderstormProbability               uint8         `json:"thunderstorm_probability"`
	ProbabilityOfFrozenPrecipitation      float64       `json:"probability_of_frozen_precipitation"`
	CloudAreaFraction                     uint8         `json:"cloud_area_fraction"`
	LowTypeCloudAreaFraction              uint8         `json:"low_type_cloud_area_fraction"`
	MediumTypeCloudAreaFraction           uint8         `json:"medium_type_cloud_area_fraction"`
	HighTypeCloudAreaFraction             uint8         `json:"high_type_cloud_area_fraction"`
	CloudBaseAltitude                     float64       `json:"cloud_base_altitude"`
	CloudTopAltitude                      float64       `json:"cloud_top_altitude"`
	PrecipitationAmountMeanDeterministic  float64       `json:"precipitation_amount_mean_deterministic"`
	PrecipitationAmountMean               float64       `json:"precipitation_amount_mean"`
	PrecipitationAmountMin                float64       `json:"precipitation_amount_min"`
	PrecipitationAmountMax                float64       `json:"precipitation_amount_max"`
	PrecipitationAmountMedian             float64       `json:"precipitation_amount_median"`
	ProbabilityOfPrecipitation            uint8         `json:"probability_of_precipitation"`
	PrecipitationFrozenPart               int8          `json:"precipitation_frozen_part"`
	PredominantPrecipitationTypeAtSurface uint8         `json:"predominant_precipitation_type_at_surface"`
	SymbolCode                            WeatherSymbol `json:"symbol_code"`
}

// Forecast defines the structure that holds the converted TimeSeries data
// from the data returned by the SMHI point forecast API.
type Forecast struct {
	Hash                               string
	Timestamp                          time.Time
	IntervalParametersStartTime        time.Time
	AirPressure                        float64
	AirTemperature                     float64
	HorizontalVisibility               float64
	MaximumPrecipitationIntensity      float64
	MeanPrecipitationIntensity         float64
	MeanValueOfHighLevelCloudCover     uint8
	MeanValueOfLowLevelCloudCover      uint8
	MeanValueOfMediumLevelCloudCover   uint8
	MeanValueOfTotalCloudCover         uint8
	MedianPrecipitationIntensity       float64
	MinimumPrecipitationIntensity      float64
	ProbabilityOfFrozenPrecipitation   float64
	ProbabilityOfPrecipitation         uint8
	PercentOfPrecipitationInFrozenForm int8
	PrecipitationCategory              PrecipitationCategory
	PrecipitationCategoryDescription   map[string]string
	RelativeHumidity                   uint8
	ThunderProbability                 uint8
	WeatherSymbol                      WeatherSymbol
	WeatherSymbolDescription           map[string]string
	WindDirection                      uint16
	WindGustSpeed                      float64
	WindSpeed                          float64
	WindSpeedDescription               map[string]string
	CloudBaseAltitude                  float64
	CloudTopAltitude                   float64
	DeterministicPrecipitationAmount   float64
}

// PointForecast holds the data for a complete PointForecast request.
type PointForecast struct {
	ApprovedTime  time.Time
	CreatedTime   time.Time
	ReferenceTime time.Time
	Geometry      Geometry
	TimeSeries    []Forecast
}

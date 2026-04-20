package smhi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	forecastURL = "https://opendata-download-metfcst.smhi.se/api/category/snow1g/version/1/geotype/point/lon/%f/lat/%f/data.json"
)

// GetPointForecast fetches a forecast from the SMHI API for the given
// longitude and latitude.
func GetPointForecast(lon, lat float64) (*PointForecast, error) {
	var err error

	// Fetch the forecast for the given longitude and latitude.
	var res *http.Response
	if res, err = http.Get(fmt.Sprintf(forecastURL, lon, lat)); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read all of the data into a buffer.
	var data []byte
	if data, err = ioutil.ReadAll(res.Body); err != nil {
		return nil, err
	}

	// Decode the data into the data structure that's defined by SMHI.
	var decodedData PointForecastAPI
	if err = json.Unmarshal(data, &decodedData); err != nil {
		return nil, err
	}

	// Create a new copy of the data in a structure that is defined by us,
	// which makes it easier to find the given temperature etc.
	var ret *PointForecast
	if ret, err = toPointForecast(&decodedData); err != nil {
		return nil, err
	}

	return ret, nil
}

// toPointForecast convers the PointForecastAPI object to a PointForecase
// object.
func toPointForecast(d *PointForecastAPI) (*PointForecast, error) {
	var ret PointForecast
	var err error
	approvedTime := d.CreatedTime
	if approvedTime == "" {
		approvedTime = d.ApprovedTime
	}

	// Fill if with some basic data.
	if ret.ApprovedTime, err = time.Parse(time.RFC3339, approvedTime); err != nil {
		return nil, err
	}
	ret.CreatedTime = ret.ApprovedTime
	if d.CreatedTime != "" {
		if ret.CreatedTime, err = time.Parse(time.RFC3339, d.CreatedTime); err != nil {
			return nil, err
		}
	}
	if d.CreatedTime == "" {
		ret.CreatedTime = ret.ApprovedTime
	}
	if ret.ReferenceTime, err = time.Parse(time.RFC3339, d.ReferenceTime); err != nil {
		return nil, err
	}
	ret.Geometry = Geometry{
		Type:        d.Geometry.Type,
		Coordinates: d.Geometry.Coordinates,
	}

	// Iterate over the time series and construct a Forecast map for each
	// timestamp.
	for _, t := range d.TimeSeries {
		var f Forecast
		if f.Timestamp, err = time.Parse(time.RFC3339, t.Time); err != nil {
			return nil, err
		}
		if f.IntervalParametersStartTime, err = time.Parse(time.RFC3339, t.IntervalParametersStartTime); err != nil {
			return nil, err
		}

		f.AirPressure = t.Data.AirPressureAtMeanSeaLevel
		f.AirTemperature = t.Data.AirTemperature
		f.HorizontalVisibility = t.Data.VisibilityInAir
		f.MaximumPrecipitationIntensity = t.Data.PrecipitationAmountMax
		f.MeanPrecipitationIntensity = t.Data.PrecipitationAmountMean
		f.MeanValueOfHighLevelCloudCover = t.Data.HighTypeCloudAreaFraction
		f.MeanValueOfLowLevelCloudCover = t.Data.LowTypeCloudAreaFraction
		f.MeanValueOfMediumLevelCloudCover = t.Data.MediumTypeCloudAreaFraction
		f.MeanValueOfTotalCloudCover = t.Data.CloudAreaFraction
		f.MedianPrecipitationIntensity = t.Data.PrecipitationAmountMedian
		f.MinimumPrecipitationIntensity = t.Data.PrecipitationAmountMin
		f.ProbabilityOfFrozenPrecipitation = t.Data.ProbabilityOfFrozenPrecipitation
		f.ProbabilityOfPrecipitation = t.Data.ProbabilityOfPrecipitation
		f.PercentOfPrecipitationInFrozenForm = t.Data.PrecipitationFrozenPart
		f.PrecipitationCategory = PrecipitationCategory(t.Data.PredominantPrecipitationTypeAtSurface)
		f.PrecipitationCategoryDescription = getPrecipitationCategoryDescriptions(f.PrecipitationCategory)
		f.RelativeHumidity = t.Data.RelativeHumidity
		f.ThunderProbability = t.Data.ThunderstormProbability
		f.WeatherSymbol = t.Data.SymbolCode
		f.WeatherSymbolDescription = getWeatherSymbolDescription(f.WeatherSymbol)
		f.WindDirection = t.Data.WindFromDirection
		f.WindGustSpeed = t.Data.WindSpeedOfGust
		f.WindSpeed = t.Data.WindSpeed
		f.WindSpeedDescription = getWindSpeedDescription(f.WindSpeed)
		f.CloudBaseAltitude = t.Data.CloudBaseAltitude
		f.CloudTopAltitude = t.Data.CloudTopAltitude
		f.DeterministicPrecipitationAmount = t.Data.PrecipitationAmountMeanDeterministic

		f.Hash = getHash(&f)

		ret.TimeSeries = append(ret.TimeSeries, f)
	}

	return &ret, nil
}

// getPrecipitationCategoryDescriptions returns a friendly precipitation
// category description.
func getPrecipitationCategoryDescriptions(pc PrecipitationCategory) map[string]string {
	ret := make(map[string]string)

	switch pc {
	case NoPrecipitation:
		ret["sv-SE"] = "Ingen nederbörd"
		ret["en-US"] = "No precipitation"
		break
	case Snow:
		ret["sv-SE"] = "Snö"
		ret["en-US"] = "Snow"
		break
	case SnowAndRain:
		ret["sv-SE"] = "Snö och regn"
		ret["en-US"] = "Snow and rain"
		break
	case Rain:
		ret["sv-SE"] = "Regn"
		ret["en-US"] = "Rain"
		break
	case Drizzle:
		ret["sv-SE"] = "Duggregn"
		ret["en-US"] = "Drizzle"
		break
	case FreezingRain:
		ret["sv-SE"] = "Frysande regn"
		ret["en-US"] = "Freezing rain"
		break
	case FreezingDrizzle:
		ret["sv-SE"] = "Underkylt regn"
		ret["en-US"] = "Freezing drizzle"
		break
	}

	return ret
}

// getWeatherSymbolDescription returns a friendly weather symbol description.
func getWeatherSymbolDescription(ws WeatherSymbol) map[string]string {
	ret := make(map[string]string)

	switch ws {
	case ClearSky:
		ret["sv-SE"] = "Klar himmel"
		ret["en-US"] = "Clear sky"
		break
	case NearlyClearSky:
		ret["sv-SE"] = "Nästan klar himmel"
		ret["en-US"] = "Nearly clear sky"
		break
	case VariableCloudiness:
		ret["sv-SE"] = "Växlande molnighet"
		ret["en-US"] = "Variable cloudiness"
		break
	case HalfclearSky:
		ret["sv-SE"] = "Halvklar himmel"
		ret["en-US"] = "Halfclear sky"
		break
	case CloudySky:
		ret["sv-SE"] = "Molnig himmel"
		ret["en-US"] = "Cloudy sky"
		break
	case Overcast:
		ret["sv-SE"] = "Mulet"
		ret["en-US"] = "Overcast"
		break
	case Fog:
		ret["sv-SE"] = "Dimma"
		ret["en-US"] = "Fog"
		break
	case LightRainShowers:
		ret["sv-SE"] = "Lätta regnskurar"
		ret["en-US"] = "Light rain showers"
		break
	case ModerateRainShowers:
		ret["sv-SE"] = "Måttliga regnskurar"
		ret["en-US"] = "Moderate rain showers"
		break
	case HeavyRainShowers:
		ret["sv-SE"] = "Kraftiga regnskurar"
		ret["en-US"] = "Heavy rain showers"
		break
	case Thunderstorm:
		ret["sv-SE"] = "Åskoväder"
		ret["en-US"] = "Thunderstorm"
		break
	case LightSleetShowers:
		ret["sv-SE"] = "Lätta regnskurar"
		ret["en-US"] = "Light sleet showers"
		break
	case ModerateSleetShowers:
		ret["sv-SE"] = "Måttliga regnskurar"
		ret["en-US"] = "Moderate sleet showers"
		break
	case HeavySleetShowers:
		ret["sv-SE"] = "Kraftiga regnskurar"
		ret["en-US"] = "Heavy sleet showers"
		break
	case LightSnowShowers:
		ret["sv-SE"] = "Lätta snöbyar"
		ret["en-US"] = "Light snow showers"
		break
	case ModerateSnowShowers:
		ret["sv-SE"] = "Måttliga snöbyar"
		ret["en-US"] = "Moderate snow showers"
		break
	case HeavySnowShowers:
		ret["sv-SE"] = "Kraftiga snöbyar"
		ret["en-US"] = "Heavy snow showers"
		break
	case LightRain:
		ret["sv-SE"] = "Duggregn"
		ret["en-US"] = "Light rain"
		break
	case ModerateRain:
		ret["sv-SE"] = "Måttligt regn"
		ret["en-US"] = "Moderate rain"
		break
	case HeavyRain:
		ret["sv-SE"] = "Kraftigt regn"
		ret["en-US"] = "Heavy rain"
		break
	case Thunder:
		ret["sv-SE"] = "Åska"
		ret["en-US"] = "Thunder"
		break
	case LightSleet:
		ret["sv-SE"] = "Lätt snöblandat regn"
		ret["en-US"] = "Light sleet"
		break
	case ModerateSleet:
		ret["sv-SE"] = "Måttligt snöblandat regn"
		ret["en-US"] = "Moderate sleet"
		break
	case HeavySleet:
		ret["sv-SE"] = "Kraftigt snöblandat regn"
		ret["en-US"] = "Heavy sleet"
		break
	case LightSnowfall:
		ret["sv-SE"] = "Lätt snöfall"
		ret["en-US"] = "Light snowfall"
		break
	case ModerateSnowfall:
		ret["sv-SE"] = "Måttligt snöfall"
		ret["en-US"] = "Moderate snowfall"
		break
	case HeavySnowfall:
		ret["sv-SE"] = "Kraftigt snöfall"
		ret["en-US"] = "Heavy snowfall"
		break
	}

	return ret
}

// getWindSpeedDescription returns a friendly name for the wind speed.
func getWindSpeedDescription(windSpeed float64) map[string]string {
	ret := make(map[string]string)

	if windSpeed <= 0.2 {
		ret["sv-SE"] = "Stiltje"
		ret["en-US"] = "Calm"
	} else if windSpeed >= 0.3 && windSpeed <= 1.5 {
		ret["sv-SE"] = "Nästan stiltje"
		ret["en-US"] = "Light air"
	} else if windSpeed >= 1.6 && windSpeed <= 3.3 {
		ret["sv-SE"] = "Lätt bris"
		ret["en-US"] = "Light breeze"
	} else if windSpeed >= 3.4 && windSpeed <= 5.4 {
		ret["sv-SE"] = "God bris"
		ret["en-US"] = "Gentle breeze"
	} else if windSpeed >= 5.5 && windSpeed <= 7.9 {
		ret["sv-SE"] = "Frisk bris"
		ret["en-US"] = "Moderate breeze"
	} else if windSpeed >= 8.0 && windSpeed <= 10.7 {
		ret["sv-SE"] = "Styv bris"
		ret["en-US"] = "Fresh breeze"
	} else if windSpeed >= 10.8 && windSpeed <= 13.8 {
		ret["sv-SE"] = "Hård bris"
		ret["en-US"] = "Strong breeze"
	} else if windSpeed >= 13.9 && windSpeed <= 17.1 {
		ret["sv-SE"] = "Styv kuling"
		ret["en-US"] = "Moderate gale"
	} else if windSpeed >= 17.2 && windSpeed <= 20.7 {
		ret["sv-SE"] = "Hård kuling"
		ret["en-US"] = "Fresh gate"
	} else if windSpeed >= 20.8 && windSpeed <= 24.4 {
		ret["sv-SE"] = "Halv storm"
		ret["en-US"] = "Strong gale"
	} else if windSpeed >= 24.5 && windSpeed <= 28.4 {
		ret["sv-SE"] = "Storm"
		ret["en-US"] = "Storm"
	} else if windSpeed >= 28.5 && windSpeed <= 32.6 {
		ret["sv-SE"] = "Svår storm"
		ret["en-US"] = "Violent storm"
	} else {
		ret["sv-SE"] = "Orkan"
		ret["en-US"] = "Hurricane"
	}
	return ret
}

func getHash(f *Forecast) string {
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v",
		f.AirPressure,
		f.AirTemperature,
		f.HorizontalVisibility,
		f.MaximumPrecipitationIntensity,
		f.MeanPrecipitationIntensity,
		f.MeanValueOfHighLevelCloudCover,
		f.MeanValueOfLowLevelCloudCover,
		f.MeanValueOfMediumLevelCloudCover,
		f.MeanValueOfTotalCloudCover,
		f.MedianPrecipitationIntensity,
		f.MinimumPrecipitationIntensity,
		f.PercentOfPrecipitationInFrozenForm,
		f.PrecipitationCategory,
		f.RelativeHumidity,
		f.ThunderProbability,
		f.WeatherSymbol,
		f.WindDirection,
		f.WindGustSpeed,
		f.WindSpeed,
	)

}

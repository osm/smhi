package smhi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	forecastURL = "https://opendata-download-metfcst.smhi.se/api/category/pmp3g/version/2/geotype/point/lon/%f/lat/%f/data.json"
)

// GetPointForecast fetches a forecast from the SMHI API for the given
// longitude and latitude.
func GetPointForcecast(lon, lat float64) (*PointForecast, error) {
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

	// Fill if with some basic data.
	if ret.ApprovedTime, err = time.Parse(time.RFC3339, d.ApprovedTime); err != nil {
		return nil, err
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
		f.Timestamp, err = time.Parse(time.RFC3339, t.ValidTime)

		for _, p := range t.Parameters {
			switch p.Name {
			case "msl":
				f.AirPressure = p.Values[0]
				break
			case "t":
				f.AirTemperature = p.Values[0]
				break
			case "vis":
				f.HorizontalVisibility = p.Values[0]
				break
			case "wd":
				f.WindDirection = uint8(p.Values[0])
				break
			case "ws":
				f.WindSpeed = p.Values[0]
				break
			case "r":
				f.RelativeHumidity = uint8(p.Values[0])
				break
			case "tstm":
				f.ThunderProbability = uint8(p.Values[0])
				break
			case "tcc_mean":
				f.MeanValueOfTotalCloudCover = uint8(p.Values[0])
				break
			case "lcc_mean":
				f.MeanValueOfLowLevelCloudCover = uint8(p.Values[0])
				break
			case "mcc_mean":
				f.MeanValueOfMediumLevelCloudCover = uint8(p.Values[0])
				break
			case "hcc_mean":
				f.MeanValueOfHighLevelCloudCover = uint8(p.Values[0])
				break
			case "gust":
				f.WindGustSpeed = p.Values[0]
				break
			case "pmin":
				f.MinimumPrecipitationIntensity = p.Values[0]
				break
			case "pmax":
				f.MaximumPrecipitationIntensity = p.Values[0]
				break
			case "spp":
				f.PercentOfPrecipitationInFrozenForm = int8(p.Values[0])
				break
			case "pcat":
				f.PrecipitationCategory = PrecipitationCategory(p.Values[0])
				f.PrecipitationCategoryDescription = getPrecipitationCategoryDescriptions(f.PrecipitationCategory)
				break
			case "pmean":
				f.MeanPrecipitationIntensity = p.Values[0]
				break
			case "pmedian":
				f.MedianPrecipitationIntensity = p.Values[0]
				break
			case "Wsymb2":
				f.WeatherSymbol = WeatherSymbol(p.Values[0])
				f.WeatherSymbolDescription = getWeatherSymbolDescription(f.WeatherSymbol)
				break
			}
		}

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
		ret["SE"] = "Ingen nederbörd"
		ret["US"] = "No precipitation"
		break
	case Snow:
		ret["SE"] = "Snö"
		ret["US"] = "Snow"
		break
	case SnowAndRain:
		ret["SE"] = "Snö och regn"
		ret["US"] = "Snow and rain"
		break
	case Rain:
		ret["SE"] = "Regn"
		ret["US"] = "Rain"
		break
	case Drizzle:
		ret["SE"] = "Duggregn"
		ret["US"] = "Drizzle"
		break
	case FreezingRain:
		ret["SE"] = "Frysande regn"
		ret["US"] = "Freezing rain"
		break
	case FreezingDrizzle:
		ret["SE"] = "Underkylt regn"
		ret["US"] = "Freezing drizzle"
		break
	}

	return ret
}

// getWeatherSymbolDescription returns a friendly weather symbol description.
func getWeatherSymbolDescription(ws WeatherSymbol) map[string]string {
	ret := make(map[string]string)

	switch ws {
	case ClearSky:
		ret["SE"] = "Klar himmel"
		ret["US"] = "Clear sky"
		break
	case NearlyClearSky:
		ret["SE"] = "Nästan klar himmel"
		ret["US"] = "Nearly clear sky"
		break
	case VariableCloudiness:
		ret["SE"] = "Växlande molnighet"
		ret["US"] = "Variable cloudiness"
		break
	case HalfclearSky:
		ret["SE"] = "Halvklar himmel"
		ret["US"] = "Halfclear sky"
		break
	case CloudySky:
		ret["SE"] = "Molnig himmel"
		ret["US"] = "Cloudy sky"
		break
	case Overcast:
		ret["SE"] = "Mulet"
		ret["US"] = "Overcast"
		break
	case Fog:
		ret["SE"] = "Dimma"
		ret["US"] = "Fog"
		break
	case LightRainShowers:
		ret["SE"] = "Lätta regnskurar"
		ret["US"] = "Light rain showers"
		break
	case ModerateRainShowers:
		ret["SE"] = "Måttliga regnskurar"
		ret["US"] = "Moderate rain showers"
		break
	case HeavyRainShowers:
		ret["SE"] = "Kraftiga regnskurar"
		ret["US"] = "Heavy rain showers"
		break
	case Thunderstorm:
		ret["SE"] = "Åskoväder"
		ret["US"] = "Thunderstorm"
		break
	case LightSleetShowers:
		ret["SE"] = "Lätta regnskurar"
		ret["US"] = "Light sleet showers"
		break
	case ModerateSleetShowers:
		ret["SE"] = "Måttliga regnskurar"
		ret["US"] = "Moderate sleet showers"
		break
	case HeavySleetShowers:
		ret["SE"] = "Kraftiga regnskurar"
		ret["US"] = "Heavy sleet showers"
		break
	case LightSnowShowers:
		ret["SE"] = "Lätta snöbyar"
		ret["US"] = "Light snow showers"
		break
	case ModerateSnowShowers:
		ret["SE"] = "Måttliga snöbyar"
		ret["US"] = "Moderate snow showers"
		break
	case HeavySnowShowers:
		ret["SE"] = "Kraftiga snöbyar"
		ret["US"] = "Heavy snow showers"
		break
	case LightRain:
		ret["SE"] = "Duggregn"
		ret["US"] = "Light rain"
		break
	case ModerateRain:
		ret["SE"] = "Måttligt regn"
		ret["US"] = "Moderate rain"
		break
	case HeavyRain:
		ret["SE"] = "Kraftigt regn"
		ret["US"] = "Heavy rain"
		break
	case Thunder:
		ret["SE"] = "Åska"
		ret["US"] = "Thunder"
		break
	case LightSleet:
		ret["SE"] = "Lätt snöblandat regn"
		ret["US"] = "Light sleet"
		break
	case ModerateSleet:
		ret["SE"] = "Måttligt snöblandat regn"
		ret["US"] = "Moderate sleet"
		break
	case HeavySleet:
		ret["SE"] = "Kraftigt snöblandat regn"
		ret["US"] = "Heavy sleet"
		break
	case LightSnowfall:
		ret["SE"] = "Lätt snöfall"
		ret["US"] = "Light snowfall"
		break
	case ModerateSnowfall:
		ret["SE"] = "Måttligt snöfall"
		ret["US"] = "Moderate snowfall"
		break
	case HeavySnowfall:
		ret["SE"] = "Kraftigt snöfall"
		ret["US"] = "Heavy snowfall"
		break
	}

	return ret
}

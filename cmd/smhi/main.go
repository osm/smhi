package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/osm/smhi"
)

func main() {
	lon := flag.Float64("lon", 0, "longitude")
	lat := flag.Float64("lat", 0, "latitude")
	flag.Parse()

	var f *smhi.PointForecast
	var err error
	if f, err = smhi.GetPointForecast(*lon, *lat); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, t := range f.TimeSeries {
		fmt.Println(t.Timestamp, t.WeatherSymbolDescription["sv-SE"], t.AirTemperature, "C")
	}
}

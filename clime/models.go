package clime

import (
	"time"

	"github.com/barnardn/go-clime/compass"
	"github.com/barnardn/go-clime/openweathermap"
	"github.com/bobg/go-generics/slices"
	"github.com/moznion/go-optional"
)

type GeoLocation struct {
	Lat, Lon float32
}

type TempUnits int8

const (
	Celcius TempUnits = iota
	Fahrenheight
)

type Velocity struct {
	Magnitude float32
	Direction compass.CompassDirection
}

type Temperature struct {
	Magnitude float32
	Units     TempUnits
}

type CurrentTemps struct {
	Current   Temperature
	FeelsLike Temperature
	Low       Temperature
	High      Temperature
	Humidity  int
}

type LocationInfo struct {
	City           string
	Country        string
	Position       GeoLocation
	TimezoneOffset int
	Sunrise        time.Duration
	Sunset         time.Duration
}

type WindInfo struct {
	Current Velocity
	Gust    optional.Option[Velocity]
}

type ConditionsDescription struct {
	Summary string
	Details string
}

type CurrentConditions struct {
	Location           LocationInfo
	Visibility         int
	TemperatureDetails CurrentTemps
	Conditions         []ConditionsDescription
	Wind               optional.Option[WindInfo]
	Clouds             optional.Option[int]
	Rain               optional.Option[float32]
	Snow               optional.Option[float32]
}

func NewCurrentConditions(owmData openweathermap.Container) CurrentConditions {
	location := LocationInfo{
		City:           owmData.Name,
		Country:        owmData.Sys.Country,
		Position:       GeoLocation{Lat: owmData.Coord.Lat, Lon: owmData.Coord.Lon},
		TimezoneOffset: owmData.Timezone,
		Sunrise:        owmData.Sys.Sunrise,
		Sunset:         owmData.Sys.Sunset,
	}
	tempDeets := CurrentTemps{
		Current:   Temperature{Magnitude: owmData.Main.Temp, Units: Celcius},
		FeelsLike: Temperature{Magnitude: owmData.Main.FeelsLike, Units: Celcius},
		Low:       Temperature{Magnitude: owmData.Main.TempMin, Units: Celcius},
		High:      Temperature{Magnitude: owmData.Main.TempMax, Units: Celcius},
		Humidity:  owmData.Main.Humidity,
	}
	conditions, _ := slices.Map(
		owmData.Weather,
		func(idx int, w openweathermap.Conditions) (ConditionsDescription, error) {
			return ConditionsDescription{Summary: w.Main, Details: w.Description}, nil
		},
	)

	wind := WindInfo{
		Current: Velocity{
			Magnitude: owmData.Wind.Speed,
			Direction: *compass.FromDegrees(float32(owmData.Wind.Deg)),
		},
		Gust: optional.None[Velocity](),
	}
	return CurrentConditions{
		Location:           location,
		Visibility:         owmData.Visibility,
		TemperatureDetails: tempDeets,
		Conditions:         conditions,
		Wind:               optional.Some(wind),
		Rain:               optional.Some(owmData.Rain.Rate),
		Snow:               optional.Some(owmData.Snow.Rate),
	}
}

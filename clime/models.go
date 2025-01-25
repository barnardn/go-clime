package clime

import (
	"fmt"
	"strings"
	"time"

	"github.com/barnardn/go-clime/compass"
	"github.com/barnardn/go-clime/openweathermap"
	"github.com/bobg/go-generics/slices"
	"github.com/moznion/go-optional"
)

type GeoLocation struct {
	Lat, Lon float32
}

func (g *GeoLocation) String() string {
	return fmt.Sprintf("Lat: %0.8f, %0.8f", g.Lat, g.Lon)
}

type TempUnits uint8

const (
	Celcius TempUnits = iota
	Fahrenheit
)

type VelocityUnits uint8

const (
	MPerSec VelocityUnits = iota
	MPH
)

type Velocity struct {
	Magnitude float32
	Direction compass.CompassDirection
	Units     VelocityUnits
}

func (v *Velocity) String() string {
	return fmt.Sprintf("%s at %0.2f m/sec", v.Direction.String(), v.Magnitude)
}

type Temperature struct {
	Magnitude float32
	Units     TempUnits
}

func (t *Temperature) String() string {
	unitLabel := "℃"
	if t.Units == Fahrenheit {
		unitLabel = "℉"
	}
	return fmt.Sprintf("%0.2f%s", t.Magnitude, unitLabel)
}

type CurrentTemps struct {
	Current   Temperature
	FeelsLike Temperature
	Low       Temperature
	High      Temperature
	Humidity  int
}

func (ct *CurrentTemps) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Current: %s\n", ct.Current.String())
	fmt.Fprintf(&b, "Feels Like: %s\n", ct.FeelsLike.String())
	fmt.Fprintf(&b, "Low: %s\n", ct.Low.String())
	fmt.Fprintf(&b, "High: %s\n", ct.High.String())
	return b.String()
}

type LocationInfo struct {
	City           string
	Country        string
	Position       GeoLocation
	TimezoneOffset int
	Sunrise        time.Time
	Sunset         time.Time
}

func (li *LocationInfo) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s, %s\n", li.City, li.Country)
	fmt.Fprintf(&b, "%s\n", li.Position.String())
	fmt.Fprintf(&b, "UTC Offset %d\n", li.TimezoneOffset)
	sunrise := li.Sunrise.Local().Format("Mon Jan 2 15:04:05")
	sunset := li.Sunset.Local().Format("Mon Jan 2 15:04:05")
	fmt.Fprintf(&b, "Sunrise: %s\n", sunrise)
	fmt.Fprintf(&b, "Sunset: %s\n", sunset)
	return b.String()
}

type WindInfo struct {
	Current Velocity
	Gust    optional.Option[Velocity]
}

func (wi *WindInfo) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Wind Speed: %s\n", wi.Current.String())
	wi.Gust.IfSome(func(g Velocity) {
		fmt.Fprintf(&b, "Gusts: %s\n", g.String())
	})
	return b.String()
}

type ConditionsDescription struct {
	Summary string
	Details string
}

func (c *ConditionsDescription) String() string {
	return fmt.Sprintf("%s: %s\n", c.Summary, c.Details)
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

func (cc *CurrentConditions) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, "Current Weather Conditions")
	fmt.Fprintln(&b, "------- ------- ----------")
	fmt.Fprintln(&b, cc.Location.String())
	fmt.Fprintln(&b, cc.TemperatureDetails.String())
	fmt.Fprintf(&b, "Visibility: %d Meters\n", cc.Visibility)
	for _, cd := range cc.Conditions {
		fmt.Fprint(&b, cd.String())
	}
	cc.Wind.IfSome(func(wi WindInfo) {
		fmt.Fprintln(&b, wi.String())
	})
	cc.Clouds.IfSome(func(v int) {
		fmt.Fprintf(&b, "Cloud Cover: %d%%\n", v)
	})
	cc.Rain.IfSome(func(v float32) {
		fmt.Fprintf(&b, "Rain: %0.1f mm/hour", v)
	})
	cc.Snow.IfSome(func(v float32) {
		fmt.Fprintf(&b, "Snow: %0.1f mm/hour", v)
	})

	return b.String()
}

func NewCurrentConditions(owmData openweathermap.Container) CurrentConditions {
	location := LocationInfo{
		City:           owmData.Name,
		Country:        owmData.Sys.Country,
		Position:       GeoLocation{Lat: owmData.Coord.Lat, Lon: owmData.Coord.Lon},
		TimezoneOffset: owmData.Timezone,
		Sunrise:        time.Unix(owmData.Sys.Sunrise, 0),
		Sunset:         time.Unix(owmData.Sys.Sunset, 0),
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

	wind := optional.Map(owmData.Wind, func(w openweathermap.Wind) WindInfo {
		return WindInfo{
			Current: Velocity{
				Magnitude: w.Speed,
				Direction: *compass.FromDegrees(float32(w.Deg)),
				Units:     MPerSec,
			},
			Gust: optional.Map(w.Gust, func(v float32) Velocity {
				return Velocity{
					Magnitude: v,
					Direction: *compass.FromDegrees(float32(w.Deg)),
					Units:     MPerSec,
				}
			}),
		}
	})

	rain := optional.Map(owmData.Rain, func(r openweathermap.HourlyRate) float32 {
		return r.Rate
	})

	snow := optional.Map(owmData.Rain, func(s openweathermap.HourlyRate) float32 {
		return s.Rate
	})

	return CurrentConditions{
		Location:           location,
		Visibility:         owmData.Visibility,
		TemperatureDetails: tempDeets,
		Conditions:         conditions,
		Wind:               wind,
		Rain:               rain,
		Snow:               snow,
	}
}

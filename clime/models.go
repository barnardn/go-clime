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

type HourlyRate uint8

const (
	MMPerHour HourlyRate = iota
	InchPerHour
)

type RateOfFall struct {
	Magnitude float32
	Units     HourlyRate
}

func (r *RateOfFall) String() string {
	units := "mm/hour"
	if r.Units == InchPerHour {
		units = "inch/hour"
	}
	return fmt.Sprintf("%0.2f %s", r.Magnitude, units)
}

type DistanceUnits uint8

const (
	KM DistanceUnits = iota
	MILES
)

type Distance struct {
	Magnitude float32
	Units     DistanceUnits
}

func (d *Distance) String() string {
	units := "KM"
	if d.Units == MILES {
		units = "Miles"
	}
	return fmt.Sprintf("%0.2f %s", d.Magnitude, units)
}

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
	units := "m/sec"
	if v.Units == MPH {
		units = "mph"
	}
	return fmt.Sprintf("%s at %0.2f %s", v.Direction.String(), v.Magnitude, units)
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
	Visibility         Distance
	TemperatureDetails CurrentTemps
	Conditions         []ConditionsDescription
	Wind               optional.Option[WindInfo]
	Clouds             optional.Option[int]
	Rain               optional.Option[RateOfFall]
	Snow               optional.Option[RateOfFall]
}

func (cc *CurrentConditions) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, "Current Weather Conditions")
	fmt.Fprintln(&b, "------- ------- ----------")
	fmt.Fprintln(&b, cc.Location.String())
	fmt.Fprintln(&b, cc.TemperatureDetails.String())
	fmt.Fprintf(&b, "Max Visibility: %s\n", cc.Visibility.String())
	for _, cd := range cc.Conditions {
		fmt.Fprint(&b, cd.String())
	}
	cc.Wind.IfSome(func(wi WindInfo) {
		fmt.Fprintln(&b, wi.String())
	})
	cc.Clouds.IfSome(func(v int) {
		fmt.Fprintf(&b, "Cloud Cover: %d%%\n", v)
	})
	cc.Rain.IfSome(func(v RateOfFall) {
		fmt.Fprintln(&b, v.String())
	})
	cc.Snow.IfSome(func(v RateOfFall) {
		fmt.Fprintln(&b, v.String())
	})

	return b.String()
}

func NewCurrentConditions(owmData openweathermap.Container, units UnitsOfMeasure) CurrentConditions {
	location := LocationInfo{
		City:           owmData.Name,
		Country:        owmData.Sys.Country,
		Position:       GeoLocation{Lat: owmData.Coord.Lat, Lon: owmData.Coord.Lon},
		TimezoneOffset: owmData.Timezone,
		Sunrise:        time.Unix(owmData.Sys.Sunrise, 0),
		Sunset:         time.Unix(owmData.Sys.Sunset, 0),
	}
	tempUnits := Celcius
	if units == Imperial {
		tempUnits = Fahrenheit
	}
	tempDeets := CurrentTemps{
		Current:   Temperature{Magnitude: owmData.Main.Temp, Units: tempUnits},
		FeelsLike: Temperature{Magnitude: owmData.Main.FeelsLike, Units: tempUnits},
		Low:       Temperature{Magnitude: owmData.Main.TempMin, Units: tempUnits},
		High:      Temperature{Magnitude: owmData.Main.TempMax, Units: tempUnits},
		Humidity:  owmData.Main.Humidity,
	}
	conditions, _ := slices.Map(
		owmData.Weather,
		func(idx int, w openweathermap.Conditions) (ConditionsDescription, error) {
			return ConditionsDescription{Summary: w.Main, Details: w.Description}, nil
		},
	)
	windUnits := MPerSec
	if units == Imperial {
		windUnits = MPH
	}
	wind := optional.Map(owmData.Wind, func(w openweathermap.Wind) WindInfo {
		return WindInfo{
			Current: Velocity{
				Magnitude: w.Speed,
				Direction: *compass.FromDegrees(float32(w.Deg)),
				Units:     windUnits,
			},
			Gust: optional.Map(w.Gust, func(v float32) Velocity {
				return Velocity{
					Magnitude: v,
					Direction: *compass.FromDegrees(float32(w.Deg)),
					Units:     windUnits,
				}
			}),
		}
	})

	rain := optional.Map(owmData.Rain, func(r openweathermap.HourlyRate) RateOfFall {
		rateUnits := MMPerHour
		rate := r.Rate
		if units == Imperial {
			rateUnits = InchPerHour
			rate = 0.0393701 * r.Rate
		}
		return RateOfFall{Magnitude: rate, Units: rateUnits}
	})

	snow := optional.Map(owmData.Rain, func(s openweathermap.HourlyRate) RateOfFall {
		rateUnits := MMPerHour
		rate := s.Rate
		if units == Imperial {
			rateUnits = InchPerHour
			rate = 0.0393701 * s.Rate
		}
		return RateOfFall{Magnitude: rate, Units: rateUnits}
	})

	visUnits := KM
	maxVisInKM := float32(owmData.Visibility / 1000)
	if units == Imperial {
		maxVisInKM *= 0.621371
		visUnits = MILES
	}
	vis := Distance{Magnitude: maxVisInKM, Units: visUnits}

	return CurrentConditions{
		Location:           location,
		Visibility:         vis,
		TemperatureDetails: tempDeets,
		Conditions:         conditions,
		Wind:               wind,
		Rain:               rain,
		Snow:               snow,
	}
}

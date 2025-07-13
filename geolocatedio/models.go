package geolocatedio

import (
	"fmt"
	"strings"
)

type Coordinates struct {
	Lat, Lon float32
}

type LocationInfo struct {
	Ip          string
	CountryCode string
	CountryName string
	RegionName  string
	RegionCode  string
	CityName    string
	District    string
	Coord       Coordinates
	ZipCode     string
	TimeZone    string
}

type standardLookup struct {
	Ip            string
	Version       string
	AddressType   string
	ContinentCode string
	ContinentName string
	CountryCode   string
	CountryName   string
	RegionName    string
	RegionCode    string
	CityName      string
	District      string
	Latitude      float32
	Longitude     float32
	ZipCode       string
	TimeZone      string
	IddCode       string
}

func (c *Coordinates) String() string {
	return fmt.Sprintf("Lat: %0.8f, %0.8f", c.Lat, c.Lon)
}

func (l *LocationInfo) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "IP Address: %s\n", l.Ip)
	fmt.Fprintf(&b, "Country: %s\n", l.CountryName)
	fmt.Fprintf(&b, "State: %s\n", l.RegionName)
	fmt.Fprintf(&b, "City: %s\n", l.CityName)
	fmt.Fprintf(&b, "Timezone: %s\n", l.TimeZone)
	fmt.Fprintf(&b, "Zip code: %s\n", l.ZipCode)
	fmt.Fprintf(&b, "Coordinates %s", l.Coord.String())

	return b.String()
}

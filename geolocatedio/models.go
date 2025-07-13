package geolocatedio

type Coordinates struct {
	Lat, Lon float32
}

type LocationInfo struct {
	Ip            string
	CountryCode   string
	ContinentName string
	RegionName    string
	RegionCode    string
	CityName      string
	District      string
	Coord         Coordinates
	ZipCode       string
	TimeZone      string
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

package openweathermap

import (
	"github.com/moznion/go-optional"
)

type Coordinates struct {
	Lat, Lon float32
}

type Weather struct {
	Temp      float32
	FeelsLike float32 `json:"feels_like"`
	TempMin   float32 `json:"temp_min"`
	TempMax   float32 `json:"temp_max"`
	Pressure  int
	Humidity  int
}

type System struct {
	Id      int
	Country string
	Sunrise int64
	Sunset  int64
}

type Wind struct {
	Speed float32
	Deg   int
	Gust  optional.Option[float32]
}

type Cloud struct {
	All int
}

type HourlyRate struct {
	Rate float32 `json:",1h"`
}

type Conditions struct {
	Id          int
	Main        string
	Description string
	icon        string
}

type Container struct {
	Id         int
	Timezone   int
	Name       string
	Visibility int
	Coord      Coordinates
	Sys        System
	Main       Weather
	Weather    []Conditions
	Wind       optional.Option[Wind]
	Clouds     optional.Option[Cloud]
	Rain       optional.Option[HourlyRate]
	Snow       optional.Option[HourlyRate]
}

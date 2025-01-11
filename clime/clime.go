package clime

import (
	"github.com/barnardn/go-clime/openweathermap"
)

type WeatherClient[RawType any] interface {
	CurrentConditions(zip string) (&RawType, error)
}

type Clime[RawType any] struct {
	C ConditionsByZip[RawType]
}

// instantiates a clime object that is generic over it's network client
func thing() {
	owclient := openweathermap.NewClient("")
	owclient.CurrentConditions("")
	climer := Clime[openweathermap.Client]{C: owclient}
}

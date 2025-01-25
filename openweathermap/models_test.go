package openweathermap

import (
	"encoding/json"
	"testing"

	"github.com/moznion/go-optional"
	"github.com/stretchr/testify/assert"
)

func TestDeserialize(t *testing.T) {
	var jsonBlob = `
{
"coord": {
"lon": -85.5639,
"lat": 42.1938
},
"weather": [
    {
      "id": 600,
      "main": "Snow",
      "description": "light snow",
      "icon": "13n"
    },
    {
      "id": 701,
      "main": "Mist",
      "description": "mist",
      "icon": "50n"
    }
],
"base": "stations",
"main": {
"temp": 274.18,
"feels_like": 272.18,
"temp_min": 274.16,
"temp_max": 275.29,
"pressure": 1027,
"humidity": 74,
"sea_level": 1027,
"grnd_level": 994
},
"visibility": 10000,
"wind": {
"speed": 1.5,
"deg": 22
},
"clouds": {
"all": 100
},
"dt": 1735075673,
"sys": {
	"type": 1,
	"id": 3378,
	"country": "US",
	"sunrise": 1735045750,
	"sunset": 1735078494
},
"timezone": -18000,
"id": 101,
"name": "Portage",
"cod": 200
}`
	var container Container
	err := json.Unmarshal([]byte(jsonBlob), &container)
	assert.Nil(t, err)
	assert.Equal(t, "US", container.Sys.Country)
	assert.Equal(t, "Portage", container.Name)
	assert.Equal(t, float32(274.18), container.Main.Temp)

	assert.Equal(t, float32(272.18), container.Main.FeelsLike)
	assert.Equal(t, float32(274.16), container.Main.TempMin)
	assert.Equal(t, float32(275.29), container.Main.TempMax)
	assert.Equal(t, float32(1.5), container.Wind.Unwrap().Speed)
	assert.Equal(t, 22, container.Wind.Unwrap().Deg)
	assert.Equal(t, optional.None[float32](), container.Wind.Unwrap().Gust)
	assert.Equal(t, 100, container.Clouds.Unwrap().All)

}

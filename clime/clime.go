package clime

type WeatherClient[APIType any] interface {
	CurrentConditions(zip string) (*APIType, error)
}

type UnitsOfMeasure int

const (
	Metric UnitsOfMeasure = iota
	Imperial
)

type ClimeClient[APIType any] struct {
	client WeatherClient[APIType]
}

func NewClient[APIType any](weatherClient WeatherClient[APIType]) *ClimeClient[APIType] {
	return &ClimeClient[APIType]{client: weatherClient}
}

func (c *ClimeClient[APIType]) CurrentConditions(zip string) (*APIType, error) {
	return c.client.CurrentConditions(zip)
}

func (c *ClimeClient[APIType]) AsyncConditions(zip string) (chan APIType, chan error) {
	weatherChannel := make(chan APIType)
	errorChannel := make(chan error)
	go func() {
		result, error := c.client.CurrentConditions(zip)
		if error != nil {
			errorChannel <- error
			return
		}
		weatherChannel <- *result
	}()
	return weatherChannel, errorChannel
}

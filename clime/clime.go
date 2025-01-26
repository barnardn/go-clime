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

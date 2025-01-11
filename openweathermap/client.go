package openweathermap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	host          = "api.openweathermap.org"
	conditionsURL = "data/2.5/weather?appid=%s&zip=%s&units=metric"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (client *Client) CurrentConditions(zip string) (*Container, error) {
	path := fmt.Sprintf(conditionsURL, client.apiKey, zip)
	link := fmt.Sprintf("http://%s/%s", host, path)
	fetchUrl, err := url.Parse(link)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Bad url %v", err))
	}
	resp, err := http.Get(fetchUrl.String())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("bad response %v", err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cant read response %v", err))
	}
	var conditions Container
	err = json.Unmarshal([]byte(body), &conditions)
	return &conditions, err
}

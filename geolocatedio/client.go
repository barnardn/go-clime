package geolocatedio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	host               = "us-west-1.geolocated.io"
	standardLookupPath = "ip/%s?api-key=%s"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (client *Client) GeoLocation(ipAddress string) (chan LocationInfo, chan error) {
	geoChan := make(chan LocationInfo)
	errChan := make(chan error)
	go func() {
		result, err := client.geoLocation(ipAddress)
		if err != nil {
			errChan <- err
			return
		}
		geoChan <- *result
	}()
	return geoChan, errChan
}

func (client *Client) geoLocation(ipAddress string) (*LocationInfo, error) {
	path := fmt.Sprintf(standardLookupPath, ipAddress, client.apiKey)
	link := fmt.Sprintf("https://%s/%s", host, path)
	fetchUrl, err := url.Parse(link)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("bad url %v", err))
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
	var locationInfo standardLookup
	err = json.Unmarshal([]byte(body), &locationInfo)
	li := LocationInfo{
		Ip:          locationInfo.Ip,
		CountryCode: locationInfo.CountryCode,
		CountryName: locationInfo.CountryName,
		RegionName:  locationInfo.RegionName,
		RegionCode:  locationInfo.RegionCode,
		CityName:    locationInfo.CityName,
		District:    locationInfo.District,
		ZipCode:     locationInfo.ZipCode,
		TimeZone:    locationInfo.TimeZone,
		Coord: Coordinates{
			Lat: locationInfo.Latitude,
			Lon: locationInfo.Longitude,
		},
	}
	return &li, err
}

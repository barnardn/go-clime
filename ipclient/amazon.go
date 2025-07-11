package ipclient

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	host = "checkip.amazonaws.com"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (client *Client) FetchIP() (*string, error) {
	link := fmt.Sprintf("http://%s", host)
	resp, err := http.Get(link)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("bad response %v", err))
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't read response %v", err))
	}
	ip := string(body)
	return &ip, nil
}

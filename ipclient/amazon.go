package ipclient

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	host = "checkip.amazonaws.com"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (client *Client) PublicIP() (chan string, chan error) {
	ipChan := make(chan string)
	errChan := make(chan error)
	go func() {
		result, err := client.publicIP()
		if err != nil {
			errChan <- err
			return
		}
		ipChan <- *result
	}()
	return ipChan, errChan
}

func (client *Client) publicIP() (*string, error) {
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
	trimmed := strings.TrimSpace(ip)
	return &trimmed, nil
}

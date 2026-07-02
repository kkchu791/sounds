package broker

import (
	"github.com/kkchu791/sounds/internal/infra/http/restclient"
)

type Client struct {
	rc *restclient.RestClient
}

func NewClient(baseURL string) *Client {
	return &Client{
		rc: restclient.NewRestClient(baseURL),
	}
}

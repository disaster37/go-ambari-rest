package client

import (
	"crypto/tls"
	"github.com/go-resty/resty"
)

type AmbariClient struct {
	client *resty.Client
}

type Response struct {
	Href *string `json:"href,omitempty"`
}

func New(baseUrl string, login string, password string) *AmbariClient {
	return &AmbariClient{
		client: resty.New().SetHostURL(baseUrl).SetHeader("X-Requested-By", "ambari").SetBasicAuth(login, password),
	}
}

func (c *AmbariClient) SetClient(client *resty.Client) {

	if client == nil {
		panic("Client can't be empty")
	}

	c.client = client
}

func (c *AmbariClient) Client() *resty.Client {
	return c.client
}

func (c *AmbariClient) DisableVerifySSL() {
	c.client = c.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

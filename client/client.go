// This file permit to manage the Ambari client API

package client

import (
	"crypto/tls"
	"github.com/go-resty/resty"
)

// Ambari client object
type AmbariClient struct {
	client *resty.Client
}
type Response struct {
	Href *string `json:"href,omitempty"`
}

// New permit to create new Ambari client
// It return AmbariClient
func New(baseUrl string, login string, password string) *AmbariClient {
	return &AmbariClient{
		client: resty.New().SetHostURL(baseUrl).SetHeader("X-Requested-By", "ambari").SetBasicAuth(login, password),
	}
}

// Pertmit to set custom resty.Client for advance option
func (c *AmbariClient) SetClient(client *resty.Client) {

	if client == nil {
		panic("Client can't be empty")
	}

	c.client = client
}

// Client permit to return resty.Client Object
func (c *AmbariClient) Client() *resty.Client {
	return c.client
}

// DisableVerifySSL permit to disable the SSL certificat check when call Ambari webservice
func (c *AmbariClient) DisableVerifySSL() {
	c.client = c.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

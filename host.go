package client

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/host-resources.md

type Host struct {
	Response
	HostInfo *HostComponentInfo `json:"Hosts"`
}

type HostInfo struct {
	ClusterName      string `json:"cluster_name"`
	HostName         string `json:"host_name"`
	MaintenanceState string `json:"maintenance_state"`
	Rack             string `json:"rack_info"`
}

func (h *Host) String() string {
	json, _ := json.Marshal(h)
	return string(json)
}

func (c *AmbariClient) CreateHost(host *Host) (*Host, error) {

	if host == nil {
		panic("Host can't be nil")
	}

	log.Debug("Host: %s", host.String())

	// Create the Host component
	path := fmt.Sprintf("/clusters/%s/hosts/%s", host.HostInfo.ClusterName, host.HostInfo.HostName)
	jsonData, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, errors.New(resp.Status())
	}

	host, err = c.Host(host.HostInfo.ClusterName, host.HostInfo.HostName)
	if err != nil {
		return nil, err
	}
	if host == nil {
		return nil, errors.New("Can't get host that just created")
	}

	return host, nil

}

// Get cluster by ID is not supported by ambari api
func (c *AmbariClient) Host(clusterName string, hostName string) (*Host, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostName == "" {
		panic("HostName can't be empty")
	}

	path := fmt.Sprintf("/clusters/%s/hosts/%s", clusterName, hostName)

	// Get the host components
	resp, err := c.Client().R().Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to get: ", resp)
	host := &Host{}
	err = json.Unmarshal(resp.Body(), host)
	if err != nil {
		return nil, err
	}
	log.Debug("Host: ", host)

	return host, nil
}

func (c *AmbariClient) UpdateHost(host *Host) (*Host, error) {

	if host == nil {
		panic("Host can't be nil")
	}
	log.Debug("Host: ", host)

	path := fmt.Sprintf("/clusters/%s/hosts/%s", host.HostInfo.ClusterName, host.HostInfo.HostName)
	jsonData, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Put(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to update: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, errors.New(resp.Status())
	}

	// Get the Host
	host, err = c.Host(host.HostInfo.ClusterName, host.HostInfo.HostName)
	if err != nil {
		return nil, err
	}
	if host == nil {
		return nil, errors.New("Can't get host that just updated")
	}

	log.Debug("Host: ", host)

	return host, err

}

func (c *AmbariClient) DeleteHost(clusterName string, hostName string) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostName == "" {
		panic("HostName can't be empty")
	}
	path := fmt.Sprintf("/clusters/%s/hosts/%s", clusterName, hostName)

	resp, err := c.Client().R().Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response to delete host: ", resp)
	if resp.StatusCode() >= 300 {
		return errors.New(resp.Status())
	}

	return nil

}

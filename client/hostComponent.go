package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/host-component-resources.md

type HostComponent struct {
	HostComponentInfo *HostComponentInfo `json:"HostRoles"`
}

type HostComponentInfo struct {
	ClusterName   string `json:"cluster_name"`
	ComponentName string `json:"component_name,omitempty"`
	Hostname      string `json:"host_name,omitempty"`
	State         string `json:"state,omitempty"`
	DesiredState  string `json:"desired_state,omitempty"`
}

func (h *HostComponent) String() string {
	json, _ := json.Marshal(h)
	return string(json)
}

func (c *AmbariClient) CreateHostComponent(hostComponent *HostComponent) (*HostComponent, error) {

	if hostComponent == nil {
		panic("HostComponent can't be nil")
	}

	log.Debug("HostComponent: %s", hostComponent.String())

	// Create the Host component
	path := fmt.Sprintf("/clusters/%s/hosts/%s/host_components/%s", hostComponent.HostComponentInfo.ClusterName, hostComponent.HostComponentInfo.Hostname, hostComponent.HostComponentInfo.ComponentName)
	resp, err := c.Client().R().Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	// Force to install
	hostComponent.HostComponentInfo.State = "INSTALLED"
	hostComponent, err = c.UpdateHostComponent(hostComponent)
	if err != nil {
		return nil, err
	}

	return hostComponent, err

}

// Get cluster by ID is not supported by ambari api
func (c *AmbariClient) HostComponent(clusterName string, hostname string, componentName string) (*HostComponent, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("Hostname can't be empty")
	}
	if componentName == "" {
		panic("ComponentName can't be empty")
	}
	path := fmt.Sprintf("/clusters/%s/hosts/%s/host_components/%s", clusterName, hostname, componentName)

	// Get the host components
	resp, err := c.Client().R().Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to get: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		} else {
			return nil, NewAmbariError(resp.StatusCode(), resp.Status())
		}
	}
	hostComponent := &HostComponent{}
	err = json.Unmarshal(resp.Body(), hostComponent)
	if err != nil {
		return nil, err
	}
	log.Debug("HostComponent: ", hostComponent)

	return hostComponent, nil
}

func (c *AmbariClient) UpdateHostComponent(hostComponent *HostComponent) (*HostComponent, error) {

	if hostComponent == nil {
		panic("HostComponent can't be nil")
	}
	log.Debug("HostComponent: ", hostComponent)

	// Update the Cluster
	path := fmt.Sprintf("/clusters/%s/hosts/%s/host_components/%s", hostComponent.HostComponentInfo.ClusterName, hostComponent.HostComponentInfo.Hostname, hostComponent.HostComponentInfo.ComponentName)
	jsonData, err := json.Marshal(hostComponent)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Put(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to update: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	// Get the HostComponent
	hostComponent, err = c.HostComponent(hostComponent.HostComponentInfo.ClusterName, hostComponent.HostComponentInfo.Hostname, hostComponent.HostComponentInfo.ComponentName)
	if err != nil {
		return nil, err
	}
	if hostComponent == nil {
		return nil, NewAmbariError(500, "Can't get hostComponent that just updated")
	}

	log.Debug("HostComponent: ", hostComponent)

	return hostComponent, err

}

func (c *AmbariClient) DeleteHostComponent(clusterName string, hostname string, componentName string) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("HostName can't be empty")
	}
	if componentName == "" {
		panic("ComponentName can't be empty")
	}
	path := fmt.Sprintf("/clusters/%s/hosts/%s/host_components/%s", clusterName, hostname, componentName)

	resp, err := c.Client().R().Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response to delete hostComponent: ", resp)
	if resp.StatusCode() >= 300 {
		return NewAmbariError(resp.StatusCode(), resp.Status())
	}

	return nil

}

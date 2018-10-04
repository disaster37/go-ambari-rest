// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/host-component-resources.md
// This file permit to manage host component resource in Ambari API

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

// Object that reflect the Ambari API
type HostComponent struct {
	HostComponentInfo *HostComponentInfo `json:"HostRoles"`
}
type HostComponentInfo struct {
	ClusterName   string `json:"cluster_name"`
	ComponentName string `json:"component_name,omitempty"`
	Hostname      string `json:"host_name,omitempty"`
	State         string `json:"state,omitempty"`
	DesiredState  string `json:"desired_state,omitempty"`
	ServiceName   string `json:"service_name,omitempty"`
}

func (h *HostComponent) CleanBeforeSave() {
	h.HostComponentInfo.DesiredState = ""
}

// String permit to display the struct as JSON object
func (h *HostComponent) String() string {
	json, _ := json.Marshal(h)
	return string(json)
}

// CreateHostComponent permit to enable component in service to deploy it on host in second time
// It return HostComponent after the creation. The component is in INIT state. If component already exist, it do nothink.
// It return error if component not exist on service or if there are some error in API call.
func (c *AmbariClient) CreateHostComponent(hostComponent *HostComponent) (*HostComponent, error) {

	if hostComponent == nil {
		panic("HostComponent can't be nil")
	}
	log.Debugf("HostComponent: %s", hostComponent.String())

	// Check if hostcomponent is already installed
	hostComponentTemp, err := c.HostComponent(hostComponent.HostComponentInfo.ClusterName, hostComponent.HostComponentInfo.Hostname, hostComponent.HostComponentInfo.ComponentName)
	if err != nil {
		return nil, err
	}
	if hostComponentTemp != nil {
		return hostComponentTemp, nil
	}

	// Create the Host component
	hostComponent.CleanBeforeSave()
	hostComponent.HostComponentInfo.State = SERVICE_INSTALLED
	path := fmt.Sprintf("/clusters/%s/hosts/%s/host_components/%s", hostComponent.HostComponentInfo.ClusterName, hostComponent.HostComponentInfo.Hostname, hostComponent.HostComponentInfo.ComponentName)
	resp, err := c.Client().R().Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	hostComponent, err = c.HostComponent(hostComponent.HostComponentInfo.ClusterName, hostComponent.HostComponentInfo.Hostname, hostComponent.HostComponentInfo.ComponentName)
	if err != nil {
		return nil, err
	}
	for hostComponent.HostComponentInfo.State != SERVICE_INSTALLED {
		time.Sleep(5 * time.Second)
		hostComponent, err = c.HostComponent(hostComponent.HostComponentInfo.ClusterName, hostComponent.HostComponentInfo.Hostname, hostComponent.HostComponentInfo.ComponentName)
		if err != nil {
			return nil, err
		}

	}

	return hostComponent, nil

}

// HostComponent permit to load the host component
// It return the host component or nil if the component is not found
// It return error if there are some error in API call.
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

// UpdateHostComponent permit to update the host component
// It return  the host component
// It return error if the host component is not found or if there are some error in API call
func (c *AmbariClient) UpdateHostComponent(hostComponent *HostComponent) (*HostComponent, error) {

	if hostComponent == nil {
		panic("HostComponent can't be nil")
	}
	log.Debug("HostComponent: ", hostComponent)

	// Update the Cluster
	hostComponent.CleanBeforeSave()
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

// StopHostComponent permit to stop component on host
// It return the HostComponent
// Not wait that host component is stopped if service is in maintenance state or if host is in maintenance state, because it can't stop it
// It will return error if component is not found or if there are some error on API call
func (c *AmbariClient) StopHostComponent(clusterName string, hostname string, componentName string) (*HostComponent, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("Hostname can't be empty")
	}
	if componentName == "" {
		panic("ComponentName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("Hostname: ", hostname)
	log.Debug("ComponentName: ", componentName)

	// Load the hostComponent
	hostComponent, err := c.HostComponent(clusterName, hostname, componentName)
	if err != nil {
		return nil, err
	}
	if hostComponent == nil {
		return nil, NewAmbariError(404, "Component %s not found on host %s in cluster %s", componentName, hostname, clusterName)
	}

	// Check if components is already stopped
	if hostComponent.HostComponentInfo.State == SERVICE_STOPPED && hostComponent.HostComponentInfo.DesiredState == SERVICE_STOPPED {
		log.Debugf("Component %s on host %s is already stopped", componentName, hostname)
		return hostComponent, nil
	}

	// Update the host components
	hostComponent.HostComponentInfo.State = SERVICE_STOPPED
	hostComponent, err = c.UpdateHostComponent(hostComponent)
	if err != nil {
		return nil, err
	}
	// Check if it can stop the component
	service, err := c.Service(clusterName, hostComponent.HostComponentInfo.ServiceName)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, NewAmbariError(404, "Service %s not found in cluster %s", hostComponent.HostComponentInfo.ServiceName, clusterName)
	}

	host, err := c.HostOnCluster(clusterName, hostComponent.HostComponentInfo.Hostname)
	if err != nil {
		return nil, err
	}
	if host == nil {
		return nil, NewAmbariError(404, "Host %s not found in cluster %s", hostComponent.HostComponentInfo.Hostname, clusterName)
	}

	for service.ServiceInfo.MaintenanceState == MAINTENANCE_STATE_OFF && host.HostInfo.MaintenanceState == MAINTENANCE_STATE_OFF && hostComponent.HostComponentInfo.State != SERVICE_STOPPED {
		time.Sleep(5 * time.Second)
		hostComponent, err = c.HostComponent(hostComponent.HostComponentInfo.ClusterName, hostComponent.HostComponentInfo.Hostname, hostComponent.HostComponentInfo.ComponentName)
		if err != nil {
			return nil, err
		}
	}

	return hostComponent, nil

}

// StartHostComponent permit to start component on host
// It will check that the component is not on client category, because it can't start.
// It return the HostComponent
// It will return error if component is not found or if there are some error on API call
func (c *AmbariClient) StartHostComponent(clusterName string, hostname string, componentName string) (*HostComponent, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("Hostname can't be empty")
	}
	if componentName == "" {
		panic("ComponentName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("Hostname: ", hostname)
	log.Debug("ComponentName: ", componentName)

	// Load the hostComponent
	hostComponent, err := c.HostComponent(clusterName, hostname, componentName)
	if err != nil {
		return nil, err
	}
	if hostComponent == nil {
		return nil, NewAmbariError(404, "Component %s not found on host %s in cluster %s", componentName, hostname, clusterName)
	}

	// Check if components is already started
	if hostComponent.HostComponentInfo.State == SERVICE_STARTED && hostComponent.HostComponentInfo.DesiredState == SERVICE_STARTED {
		log.Debugf("Component %s on host %s is already started", componentName, hostname)
		return hostComponent, nil
	}

	// Get the category of the component
	component, err := c.Component(clusterName, hostComponent.HostComponentInfo.ServiceName, componentName)
	if err != nil {
		return nil, err
	}
	if component == nil {
		return nil, NewAmbariError(404, "Component %s not found in service %s on cluster %s", componentName, hostComponent.HostComponentInfo.ServiceName, clusterName)
	}
	if component.ServiceComponentInfo.Category == COMPONENT_CLIENT {
		log.Debugf("Component %s is client, it can't start", componentName)
		return hostComponent, nil
	}

	// Update the host components
	hostComponent.HostComponentInfo.State = SERVICE_STARTED
	hostComponent, err = c.UpdateHostComponent(hostComponent)
	if err != nil {
		return nil, err
	}

	// Check if it can start the component
	service, err := c.Service(clusterName, hostComponent.HostComponentInfo.ServiceName)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, NewAmbariError(404, "Service %s not found in cluster %s", hostComponent.HostComponentInfo.ServiceName, clusterName)
	}

	host, err := c.HostOnCluster(clusterName, hostComponent.HostComponentInfo.Hostname)
	if err != nil {
		return nil, err
	}
	if host == nil {
		return nil, NewAmbariError(404, "Host %s not found in cluster %s", hostComponent.HostComponentInfo.Hostname, clusterName)
	}
	for service.ServiceInfo.MaintenanceState == MAINTENANCE_STATE_OFF && host.HostInfo.MaintenanceState == MAINTENANCE_STATE_OFF && hostComponent.HostComponentInfo.State != SERVICE_STARTED {
		time.Sleep(5 * time.Second)
		hostComponent, err = c.HostComponent(hostComponent.HostComponentInfo.ClusterName, hostComponent.HostComponentInfo.Hostname, hostComponent.HostComponentInfo.ComponentName)
		if err != nil {
			return nil, err
		}
	}
	return hostComponent, nil
}

// DeleteHostComponent permit to delete component on host
// It will stop the component before to delete it
// It will return error if component is not found or if there are some error during API call
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

	_, err := c.StopHostComponent(clusterName, hostname, componentName)
	if err != nil {
		return err
	}

	// Then delete host components
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

// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/component-resources.md
// This file permit to manage Component in Ambari API

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const (
	COMPONENT_CLIENT = "CLIENT"
)

type Component struct {
	ComponentInfo  *ComponentInfo  `json:"ServiceComponentInfo"`
	HostComponents []HostComponent `json:"host_components"`
}
type ComponentInfo struct {
	ClusterName   string `json:"cluster_name,omitempty"`
	ServiceName   string `json:"service_name,omitempty"`
	ComponentName string `json:"component_name,omitempty"`
	State         string `json:"state,omitempty"`
	Category      string `json:"category,omitempty"`
}

// String permit to return Component as Json string
func (s *Component) String() string {
	json, _ := json.Marshal(s)
	return string(json)
}

// CreateComponent permit to create new component
// It return Component if all right fine
// It return error if something wrong when API call
func (c *AmbariClient) CreateComponent(component *Component) (*Component, error) {

	if component == nil {
		panic("Component can't be nil")
	}
	log.Debugf("Component: %s", component.String())

	path := fmt.Sprintf("/clusters/%s/services/%s/components/%s", component.ComponentInfo.ClusterName, component.ComponentInfo.ServiceName, component.ComponentInfo.ComponentName)
	resp, err := c.Client().R().Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	component, err = c.Component(component.ComponentInfo.ClusterName, component.ComponentInfo.ServiceName, component.ComponentInfo.ComponentName)
	if err != nil {
		return nil, err
	}
	if component == nil {
		return nil, NewAmbariError(500, "Can't get component that just created")
	}

	log.Debugf("Return component: %s", component)

	return component, nil

}

// Component permit to get Component item
// It return Component if found
// It return nil if service Component not found
// It return error if something wrong when API call
func (c *AmbariClient) Component(clusterName string, serviceName string, componentName string) (*Component, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if serviceName == "" {
		panic("ServiceName can't be empty")
	}
	if componentName == "" {
		panic("ComponentName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("ServiceName: ", serviceName)
	log.Debug("ComponentName: ", componentName)

	path := fmt.Sprintf("/clusters/%s/services/%s/components/%s", clusterName, serviceName, componentName)
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
	component := &Component{}
	err = json.Unmarshal(resp.Body(), component)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return component: %s", component)

	return component, nil
}

// Delete component permit to delete existing component
// before to delete component, we need to remove them from host
// It return error if component not exist
// It return error if something wrong when API call
func (c *AmbariClient) DeleteComponent(clusterName string, serviceName string, componentName string) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if serviceName == "" {
		panic("ServiceName can't be empty")
	}
	if componentName == "" {
		panic("ComponentName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("ServiceName: ", serviceName)
	log.Debug("ComponentName: ", componentName)

	// Check if component exist
	component, err := c.Component(clusterName, serviceName, componentName)
	if err != nil {
		return err
	}
	if component == nil {
		return NewAmbariError(404, "Component %s not found", componentName)
	}

	// Delete component on all host
	for _, hostComponent := range component.HostComponents {

		err := c.DeleteHostComponent(hostComponent.HostComponentInfo.ClusterName, hostComponent.HostComponentInfo.Hostname, hostComponent.HostComponentInfo.ComponentName)
		if err != nil {
			return err
		}
	}

	// Finnaly delete the component
	path := fmt.Sprintf("/clusters/%s/services/%s/components/%s", clusterName, serviceName, componentName)
	resp, err := c.Client().R().Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response to delete service: ", resp)
	if resp.StatusCode() >= 300 {
		return NewAmbariError(resp.StatusCode(), resp.Status())
	}

	return nil

}

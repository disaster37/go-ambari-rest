package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/component-resources.md

const (
	COMPONENT_CLIENT = "CLIENT"
)

func (s *ServiceComponent) String() string {
	json, _ := json.Marshal(s)
	return string(json)
}

func (c *AmbariClient) CreateComponent(component *ServiceComponent) (*ServiceComponent, error) {

	if component == nil {
		panic("Component can't be nil")
	}

	log.Debug("Component: %s", component.String())

	path := fmt.Sprintf("/clusters/%s/services/%s/components/%s", component.ServiceComponentInfo.ClusterName, component.ServiceComponentInfo.ServiceName, component.ServiceComponentInfo.ComponentName)
	resp, err := c.Client().R().Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	component, err = c.Component(component.ServiceComponentInfo.ClusterName, component.ServiceComponentInfo.ServiceName, component.ServiceComponentInfo.ComponentName)
	if err != nil {
		return nil, err
	}
	if component == nil {
		return nil, NewAmbariError(500, "Can't get component that just created")
	}

	log.Debug("Return component: %s", component)

	return component, nil

}

func (c *AmbariClient) Component(clusterName string, serviceName string, componentName string) (*ServiceComponent, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if serviceName == "" {
		panic("ServiceName can't be empty")
	}
	if componentName == "" {
		panic("ComponentName can't be empty")
	}

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
	component := &ServiceComponent{}
	err = json.Unmarshal(resp.Body(), component)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return component: %s", component)

	return component, nil
}

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

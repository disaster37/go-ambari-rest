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
	ClusterName   string `json:"cluster_name,omitempty"`
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
	hostComponent.HostComponentInfo.State = SERVICE_INIT
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

	// Install host component
	hostComponent.HostComponentInfo.State = SERVICE_INSTALLED
	request := &Request{
		RequestInfo: &RequestInfo{
			Context: fmt.Sprintf("Install component %s on %s from API", hostComponent.HostComponentInfo.ComponentName, hostComponent.HostComponentInfo.Hostname),
		},
		Body: hostComponent,
	}
	requestTask, err := c.SendRequestHostComponent(request)
	if err != nil {
		return nil, err
	}

	if requestTask != nil {
		for requestTask.RequestTaskInfo.ProgressPercent < 100 {

			requestTask, err = c.Request(hostComponent.HostComponentInfo.ClusterName, requestTask.RequestTaskInfo.Id)
			if err != nil {
				return nil, err
			}
			if requestTask == nil {
				return nil, NewAmbariError(404, "Request with Id %d not found", requestTask.RequestTaskInfo.Id)
			}

			time.Sleep(10 * time.Second)
		}

		// Check the status
		if requestTask.RequestTaskInfo.Status != REQUEST_COMPLETED {
			return nil, NewAmbariError(500, "Request %d failed with status %s, task completed %d, task aborded %d, task failed %d", requestTask.RequestTaskInfo.Id, requestTask.RequestTaskInfo.Status, requestTask.RequestTaskInfo.CompletedTask, requestTask.RequestTaskInfo.AbordedTask, requestTask.RequestTaskInfo.FailedTask)
		}
	}

	// Finnaly get the host component
	hostComponent, err = c.HostComponent(hostComponent.HostComponentInfo.ClusterName, hostComponent.HostComponentInfo.Hostname, hostComponent.HostComponentInfo.ComponentName)
	if err != nil {
		return nil, err
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

// SendRequestHostComponent permit to start / stop host component on ambari with message displayed on operation task in Ambari UI
// It keep only State field when it send the request
// It return RequestTask if all work fine
// It return nil if no request is created
// It return error if something wrong when it call the API
func (c *AmbariClient) SendRequestHostComponent(request *Request) (*RequestTask, error) {

	if request == nil {
		panic("Request can't be nil")
	}
	log.Debug("Request: ", request)
	hostComponent := request.Body.(*HostComponent)
	hostComponentTemp := &HostComponent{
		HostComponentInfo: &HostComponentInfo{
			State: hostComponent.HostComponentInfo.State,
		},
	}
	request.Body = hostComponentTemp

	log.Debug("Sended Request: ", request)

	path := fmt.Sprintf("/clusters/%s/hosts/%s/host_components/%s", hostComponent.HostComponentInfo.ClusterName, hostComponent.HostComponentInfo.Hostname, hostComponent.HostComponentInfo.ComponentName)
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Put(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response when send request: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}
	if len(resp.Body()) == 0 {
		return nil, nil
	}
	requestTask := &RequestTask{}
	err = json.Unmarshal(resp.Body(), requestTask)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return request: %s", requestTask)

	return requestTask, err
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
	request := &Request{
		RequestInfo: &RequestInfo{
			Context: fmt.Sprintf("Stop component %s on %s from API", componentName, hostname),
		},
		Body: hostComponent,
	}
	requestTask, err := c.SendRequestHostComponent(request)
	if err != nil {
		return nil, err
	}

	if requestTask != nil {
		for requestTask.RequestTaskInfo.ProgressPercent < 100 {

			requestTask, err = c.Request(clusterName, requestTask.RequestTaskInfo.Id)
			if err != nil {
				return nil, err
			}
			if requestTask == nil {
				return nil, NewAmbariError(404, "Request with Id %d not found", requestTask.RequestTaskInfo.Id)
			}

			time.Sleep(10 * time.Second)
		}

		// Check the status
		if requestTask.RequestTaskInfo.Status != REQUEST_COMPLETED {
			return nil, NewAmbariError(500, "Request %d failed with status %s, task completed %d, task aborded %d, task failed %d", requestTask.RequestTaskInfo.Id, requestTask.RequestTaskInfo.Status, requestTask.RequestTaskInfo.CompletedTask, requestTask.RequestTaskInfo.AbordedTask, requestTask.RequestTaskInfo.FailedTask)
		}
	}

	// Finnaly get the host component
	hostComponent, err = c.HostComponent(clusterName, hostname, componentName)
	if err != nil {
		return nil, err
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
	if component.ComponentInfo.Category == COMPONENT_CLIENT {
		log.Debugf("Component %s is client, it can't start", componentName)
		return hostComponent, nil
	}

	// Update the host components
	hostComponent.HostComponentInfo.State = SERVICE_STARTED
	request := &Request{
		RequestInfo: &RequestInfo{
			Context: fmt.Sprintf("Start component %s on %s from API", componentName, hostname),
		},
		Body: hostComponent,
	}
	requestTask, err := c.SendRequestHostComponent(request)
	if err != nil {
		return nil, err
	}

	if requestTask != nil {
		for requestTask.RequestTaskInfo.ProgressPercent < 100 {

			requestTask, err = c.Request(clusterName, requestTask.RequestTaskInfo.Id)
			if err != nil {
				return nil, err
			}
			if requestTask == nil {
				return nil, NewAmbariError(404, "Request with Id %d not found", requestTask.RequestTaskInfo.Id)
			}

			time.Sleep(10 * time.Second)
		}

		// Check the status
		if requestTask.RequestTaskInfo.Status != REQUEST_COMPLETED {
			return nil, NewAmbariError(500, "Request %d failed with status %s, task completed %d, task aborded %d, task failed %d", requestTask.RequestTaskInfo.Id, requestTask.RequestTaskInfo.Status, requestTask.RequestTaskInfo.CompletedTask, requestTask.RequestTaskInfo.AbordedTask, requestTask.RequestTaskInfo.FailedTask)
		}
	}

	// Finnaly get the host component
	hostComponent, err = c.HostComponent(clusterName, hostname, componentName)
	if err != nil {
		return nil, err
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

//This file permit to manage service in Ambari API
// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/services-service.md

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	SERVICE_STARTED       = "STARTED"
	SERVICE_STOPPED       = "INSTALLED"
	SERVICE_INSTALLED     = "INSTALLED"
	SERVICE_UNKNOWN       = "UNKNOWN"
	SERVICE_INIT          = "INIT"
	MAINTENANCE_STATE_ON  = "ON"
	MAINTENANCE_STATE_OFF = "OFF"
)

// Service object
type Service struct {
	ServiceInfo *ServiceInfo `json:"ServiceInfo"`
	Components  []Component  `json:"components,omitempty"`
}
type ServiceInfo struct {
	ClusterName      string `json:"cluster_name,omitempty"`
	ServiceName      string `json:"service_name,omitempty"`
	State            string `json:"state,omitempty"`
	RepositoryId     int    `json:"desired_repository_version_id,omitempty"`
	MaintenanceState string `json:"maintenance_state,omitempty"`
}

// String permit to get service object as Json string
func (s *Service) String() string {
	json, _ := json.Marshal(s)
	return string(json)
}

// ClearBeforeSave permit to clean service before save or update it
func (s *Service) CleanBeforeSave() {
	s.Components = nil
}

// CreateService permit to create new service
// The service is created in INIT state
// It return the service if all work fine
// It return error if something wrong when it call the API
func (c *AmbariClient) CreateService(service *Service) (*Service, error) {

	if service == nil {
		panic("Service can't be nil")
	}
	log.Debugf("Service: %s", service.String())

	service.CleanBeforeSave()
	service.ServiceInfo.State = SERVICE_INIT

	path := fmt.Sprintf("/clusters/%s/services/%s", service.ServiceInfo.ClusterName, service.ServiceInfo.ServiceName)
	jsonData, err := json.Marshal(service)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	service, err = c.Service(service.ServiceInfo.ClusterName, service.ServiceInfo.ServiceName)
	if err != nil {
		return nil, err
	}

	log.Debugf("Return service: %s", service)

	return service, nil

}

// Service permit to get service by is name
// It return Service if is found
// It return nil is service is not found
// It return error if something wrong with the API call
func (c *AmbariClient) Service(clusterName string, serviceName string) (*Service, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if serviceName == "" {
		panic("ServiceName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("ServiceName: ", serviceName)

	path := fmt.Sprintf("/clusters/%s/services/%s", clusterName, serviceName)
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
	service := &Service{}
	err = json.Unmarshal(resp.Body(), service)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return service: %s", service)

	return service, nil
}

// UpdateService permit to update an existing sevice like service state
// It return updated Service if all work fine
// It return error if something wrong when it call the API
func (c *AmbariClient) UpdateService(service *Service) (*Service, error) {

	if service == nil {
		panic("Service can't be nil")
	}
	log.Debug("Service: ", service)
	service.CleanBeforeSave()

	path := fmt.Sprintf("/clusters/%s/services/%s", service.ServiceInfo.ClusterName, service.ServiceInfo.ServiceName)
	jsonData, err := json.Marshal(service)
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

	// Get the service
	service, err = c.Service(service.ServiceInfo.ClusterName, service.ServiceInfo.ServiceName)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, NewAmbariError(500, "Can't get service that just updated")
	}

	log.Debugf("Return service: %s", service.String())

	return service, err

}

// DeleteService permit to delete an existing service
// Before to delete service, it need to stop service and then delete all components
// It return error if service is not found
func (c *AmbariClient) DeleteService(clusterName string, serviceName string) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if serviceName == "" {
		panic("ServiceName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("ServiceName: ", serviceName)

	// Stop service before to delete it
	_, err := c.StopService(clusterName, serviceName, false, true)
	if err != nil {
		return err
	}

	// Get service and remove all components
	service, err := c.Service(clusterName, serviceName)
	for _, component := range service.Components {
		err := c.DeleteComponent(clusterName, serviceName, component.ComponentInfo.ComponentName)
		if err != nil {
			return err
		}
	}

	// Finnaly delete the service
	path := fmt.Sprintf("/clusters/%s/services/%s", clusterName, serviceName)
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

// SendRequestService permit to start / stop service on ambari with message displayed on operation task in Ambari UI
// It keep only State and MaintenanceState field when it send the request
// It return RequestTask if all work fine
// It return nil if no request is created
// It return error if something wrong when it call the API
func (c *AmbariClient) SendRequestService(request *Request) (*RequestTask, error) {

	if request == nil {
		panic("Request can't be nil")
	}
	log.Debug("Request: ", request)
	service := request.Body.(*Service)
	serviceTemp := &Service{
		ServiceInfo: &ServiceInfo{
			State:            service.ServiceInfo.State,
			MaintenanceState: service.ServiceInfo.MaintenanceState,
		},
	}
	request.Body = serviceTemp

	log.Debug("Sended Request: ", request)

	path := fmt.Sprintf("/clusters/%s/services/%s", service.ServiceInfo.ClusterName, service.ServiceInfo.ServiceName)
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

// InstallService permit to start the service installation
// It must have the service setting and component associated to host to work
// It return service if all work fine
// It return error if something wrong when it call the API
func (c *AmbariClient) InstallService(service *Service) (*Service, error) {
	if service == nil {
		panic("Service can't be nil")
	}
	log.Debug("Service: ", service)

	// Check if service is already installed
	if service.ServiceInfo.State == SERVICE_INSTALLED {
		log.Debugf("The service %s is already installed", service.ServiceInfo.ServiceName)
		return service, nil
	}

	// Install service and wait
	service.ServiceInfo.State = SERVICE_INSTALLED
	service, err := c.UpdateService(service)
	if err != nil {
		return nil, err
	}
	for service.ServiceInfo.State != SERVICE_INSTALLED {
		time.Sleep(10 * time.Second)
		service, err = c.Service(service.ServiceInfo.ClusterName, service.ServiceInfo.ServiceName)
		if err != nil {
			return nil, err
		}
	}

	return service, nil
}

// StartService start the HDP service
// If you set disableMaintenanceMode to true, it will start the service even if maintenance state was ON
// If not, it will not stay that the service run because of it can't run it
// If disableMaintenanceMode is set to true, it will disable maintenance state before start the service
func (c *AmbariClient) StartService(clusterName string, serviceName string, disableMaintenanceMode bool) (*Service, error) {
	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if serviceName == "" {
		panic("ServiceName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("ServiceName: ", serviceName)
	log.Debug("DisableMaintenanceMode: ", disableMaintenanceMode)

	// Get the service
	service, err := c.Service(clusterName, serviceName)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, NewAmbariError(404, "Service %s not found in cluster %s", serviceName, clusterName)
	}

	// Check if service is already started
	if service.ServiceInfo.State == SERVICE_STARTED {
		log.Debugf("Service %s is already started", service.ServiceInfo.ServiceName)
		return service, nil
	}

	service.ServiceInfo.State = SERVICE_STARTED
	if disableMaintenanceMode == true {
		service.ServiceInfo.MaintenanceState = MAINTENANCE_STATE_OFF
	}
	request := &Request{
		RequestInfo: &RequestInfo{
			Context: fmt.Sprintf("Start service %s from API", serviceName),
		},
		Body: service,
	}
	requestTask, err := c.SendRequestService(request)
	if err != nil {
		return nil, err
	}

	// Wait the end of the request if request has been created
	if requestTask != nil {

		// Wait the end of the request
		err = requestTask.Wait(c, clusterName)
		if err != nil {
			return nil, err
		}

		// Check the status
		if requestTask.RequestTaskInfo.Status != REQUEST_COMPLETED {
			return nil, NewAmbariError(500, "Request %d failed with status %s, task completed %d, task aborded %d, task failed %d", requestTask.RequestTaskInfo.Id, requestTask.RequestTaskInfo.Status, requestTask.RequestTaskInfo.CompletedTask, requestTask.RequestTaskInfo.AbordedTask, requestTask.RequestTaskInfo.FailedTask)
		}
	}

	// Finnaly get the service
	service, err = c.Service(clusterName, serviceName)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// StopService stop the HDP service
// If force is set to true, It will disable maintenance state before stop the service
// It not, It will not stay that the service stop because of it can't stop it.
// If enableMaintenanceMode is set to true, it will enable maintenance state after stopped the service
func (c *AmbariClient) StopService(clusterName string, serviceName string, enableMaintenanceMode bool, force bool) (*Service, error) {
	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if serviceName == "" {
		panic("ServiceName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("ServiceName: ", serviceName)
	log.Debug("EnableMaintenanceMode: ", enableMaintenanceMode)
	log.Debug("Force: ", force)

	// Get the service
	service, err := c.Service(clusterName, serviceName)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, NewAmbariError(404, "Service %s not found in cluster %s", serviceName, clusterName)
	}

	// Check if service is already stopped
	if service.ServiceInfo.State == SERVICE_STOPPED {
		log.Debugf("The service %s is already stopped", service.ServiceInfo.ServiceName)
		return service, nil
	}

	// Stop service and wait
	service.ServiceInfo.State = SERVICE_STOPPED
	if force == true {
		service.ServiceInfo.MaintenanceState = MAINTENANCE_STATE_OFF
	}
	request := &Request{
		RequestInfo: &RequestInfo{
			Context: fmt.Sprintf("Stop service %s from API", serviceName),
		},
		Body: service,
	}
	requestTask, err := c.SendRequestService(request)
	if err != nil {
		return nil, err
	}

	if requestTask != nil {

		// Wait the end of the request
		err = requestTask.Wait(c, clusterName)
		if err != nil {
			return nil, err
		}

		// Check the status
		if requestTask.RequestTaskInfo.Status != REQUEST_COMPLETED {
			return nil, NewAmbariError(500, "Request %d failed with status %s, task completed %d, task aborded %d, task failed %d", requestTask.RequestTaskInfo.Id, requestTask.RequestTaskInfo.Status, requestTask.RequestTaskInfo.CompletedTask, requestTask.RequestTaskInfo.AbordedTask, requestTask.RequestTaskInfo.FailedTask)
		}

	}

	// Enable maintenance state if needed
	if enableMaintenanceMode == true {
		service.ServiceInfo.MaintenanceState = MAINTENANCE_STATE_ON
		service, err = c.UpdateService(service)
		if err != nil {
			return nil, err
		}
	}

	// Finnaly get the service
	service, err = c.Service(clusterName, serviceName)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// StopAllServices stop all services in HDP cluster.
// If enableMaintenanceMode is set to true, it will put all services in maintenance state after stopped all services.
// If force is set to true, it will remove maintenance state in all services before stop all services. In this way, it will stop all services.
// It returns error if the cluster not exist or if API call failed
func (c *AmbariClient) StopAllServices(cluster *Cluster, enableMaintenanceMode bool, force bool) error {
	if cluster == nil {
		panic("Cluster can't be nil")
	}
	log.Debug("Cluster: ", cluster)
	log.Debug("EnableMaintenanceMode: ", enableMaintenanceMode)
	log.Debug("Force: ", force)

	// Stop all services
	service := &Service{
		ServiceInfo: &ServiceInfo{
			State: SERVICE_STOPPED,
		},
	}
	if force == true {
		service.ServiceInfo.MaintenanceState = MAINTENANCE_STATE_OFF
		log.Debugf("Disable maintenance state before stop all services")
	}
	request := &Request{
		RequestInfo: &RequestInfo{
			Context: "Stop all services from API",
		},
		Body: service,
	}
	path := fmt.Sprintf("/clusters/%s/services", cluster.ClusterInfo.ClusterName)
	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}
	resp, err := c.Client().R().SetBody(jsonData).Put(path)
	if err != nil {
		return err
	}
	log.Debug("Response to stop all services: ", resp)
	if resp.StatusCode() >= 300 {
		return NewAmbariError(resp.StatusCode(), resp.Status())
	}

	if len(resp.Body()) == 0 {
		log.Debugf("All service already stopped")
		return nil
	}
	requestTask := &RequestTask{}
	err = json.Unmarshal(resp.Body(), requestTask)
	if err != nil {
		return err
	}
	log.Debugf("Return request: %s", requestTask)

	// Wait the end of the request
	err = requestTask.Wait(c, cluster.ClusterInfo.ClusterName)
	if err != nil {
		return err
	}

	// Check the status
	if requestTask.RequestTaskInfo.Status != REQUEST_COMPLETED {
		return NewAmbariError(500, "Request %d failed with status %s, task completed %d, task aborded %d, task failed %d", requestTask.RequestTaskInfo.Id, requestTask.RequestTaskInfo.Status, requestTask.RequestTaskInfo.CompletedTask, requestTask.RequestTaskInfo.AbordedTask, requestTask.RequestTaskInfo.FailedTask)
	}

	// Put all services in maintenance state if needed
	if enableMaintenanceMode == true {

		log.Debugf("Enable maintenance state after stop all services")
		service := &Service{
			ServiceInfo: &ServiceInfo{
				MaintenanceState: MAINTENANCE_STATE_ON,
			},
		}

		jsonData, err := json.Marshal(service)
		if err != nil {
			return err
		}
		resp, err := c.Client().R().SetBody(jsonData).Put(path)
		if err != nil {
			return err
		}
		log.Debug("Response to put all services in maintenance state: ", resp)
		if resp.StatusCode() >= 300 {
			return NewAmbariError(resp.StatusCode(), resp.Status())
		}
	}

	return nil

}

// StartAllServices start all services in HDP cluster.
// If disableMaintenanceMode is set to true, it will remove the meaintenance state in all services before start all services. In this way, all services will start.
// It return error if cluster not exist or if API call failed.
func (c *AmbariClient) StartAllServices(cluster *Cluster, disableMaintenanceMode bool) error {
	if cluster == nil {
		panic("Cluster can't be nil")
	}
	log.Debug("Cluster: ", cluster)

	// Start all services
	service := &Service{
		ServiceInfo: &ServiceInfo{
			State: SERVICE_STARTED,
		},
	}
	if disableMaintenanceMode == true {
		service.ServiceInfo.MaintenanceState = MAINTENANCE_STATE_OFF
		log.Debugf("Disable maintenance state in all services before start them")
	}
	request := &Request{
		RequestInfo: &RequestInfo{
			Context: "Start all services from API",
		},
		Body: service,
	}
	path := fmt.Sprintf("/clusters/%s/services", cluster.ClusterInfo.ClusterName)
	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}
	resp, err := c.Client().R().SetBody(jsonData).Put(path)
	if err != nil {
		return err
	}
	log.Debug("Response to start all services: ", resp)
	if resp.StatusCode() >= 300 {
		return NewAmbariError(resp.StatusCode(), resp.Status())
	}
	if len(resp.Body()) == 0 {
		log.Debugf("All service already started")
		return nil
	}
	requestTask := &RequestTask{}
	err = json.Unmarshal(resp.Body(), requestTask)
	if err != nil {
		return err
	}
	log.Debugf("Return request: %s", requestTask)

	// Wait the end of the request
	err = requestTask.Wait(c, cluster.ClusterInfo.ClusterName)
	if err != nil {
		return err
	}

	// Check the status
	if requestTask.RequestTaskInfo.Status != REQUEST_COMPLETED {
		return NewAmbariError(500, "Request %d failed with status %s, task completed %d, task aborded %d, task failed %d", requestTask.RequestTaskInfo.Id, requestTask.RequestTaskInfo.Status, requestTask.RequestTaskInfo.CompletedTask, requestTask.RequestTaskInfo.AbordedTask, requestTask.RequestTaskInfo.FailedTask)
	}

	return nil

}

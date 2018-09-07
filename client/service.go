package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/services-service.md

const (
	SERVICE_STARTED       = "STARTED"
	SERVICE_STOPPED       = "INSTALLED"
	SERVICE_INSTALLED     = "INSTALLED"
	SERVICE_UNKNOWN       = "UNKNOWN"
	SERVICE_INIT          = "INIT"
	MAINTENANCE_STATE_ON  = "ON"
	MAINTENANCE_STATE_OFF = "OFF"
)

type Service struct {
	ServiceInfo       *ServiceInfo       `json:"ServiceInfo"`
	ServiceComponents []ServiceComponent `json:"components,omitempty"`
}

type ServiceInfo struct {
	ClusterName      string `json:"cluster_name,omitempty"`
	ServiceName      string `json:"service_name,omitempty"`
	State            string `json:"state,omitempty"`
	RepositoryId     int    `json:"desired_repository_version_id,omitempty"`
	DesiredState     string `json:"desired_state,omitempty"`
	MaintenanceState string `json:"maintenance_state,omitempty"`
}

type ServiceComponent struct {
	ServiceComponentInfo *ServiceComponentInfo `json:"ServiceComponentInfo"`
}

type ServiceComponentInfo struct {
	ClusterName   string `json:"cluster_name,omitempty"`
	ServiceName   string `json:"service_name,omitempty"`
	ComponentName string `json:"component_name,omitempty"`
	State         string `json:"state,omitempty"`
	Category      string `json:"category,omitempty"`
}

func (s *Service) String() string {
	json, _ := json.Marshal(s)
	return string(json)
}

func (s *Service) CleanBeforeSave() {
	s.ServiceComponents = nil
	s.ServiceInfo.DesiredState = ""
}

func (c *AmbariClient) CreateService(service *Service) (*Service, error) {

	if service == nil {
		panic("Service can't be nil")
	}

	log.Debug("Service: %s", service.String())

	// Check if service already exist
	serviceTemp, err := c.Service(service.ServiceInfo.ClusterName, service.ServiceInfo.ServiceName)
	if err != nil {
		return nil, err
	}
	if serviceTemp != nil && (serviceTemp.ServiceInfo.State == SERVICE_INSTALLED || serviceTemp.ServiceInfo.State == SERVICE_STARTED || serviceTemp.ServiceInfo.State == SERVICE_STOPPED) {
		log.Debugf("Service %s is already installed")
		return serviceTemp, nil
	}

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

	log.Debug("Return service: %s", service)

	return service, nil

}

func (c *AmbariClient) Service(clusterName string, serviceName string) (*Service, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if serviceName == "" {
		panic("ServiceName can't be empty")
	}

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

func (c *AmbariClient) DeleteService(clusterName string, serviceName string) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if serviceName == "" {
		panic("ServiceName can't be empty")
	}

	// Stop service before to delete it
	_, err := c.StopService(clusterName, serviceName, false, true)
	if err != nil {
		return nil
	}

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

func (c *AmbariClient) InstallService(service *Service) (*Service, error) {
	if service == nil {
		panic("Service can't be nil")
	}
	log.Debug("Service: ", service)

	// Check if service is already installed
	if service.ServiceInfo.State == SERVICE_INSTALLED && service.ServiceInfo.DesiredState == SERVICE_INSTALLED {
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
		time.Sleep(5 * time.Second)
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
	if service.ServiceInfo.State == SERVICE_STARTED && service.ServiceInfo.DesiredState == SERVICE_STARTED {
		log.Debugf("Service %s is already started", service.ServiceInfo.ServiceName)
		return service, nil
	}

	service.ServiceInfo.State = SERVICE_STARTED
	if disableMaintenanceMode == true {
		service.ServiceInfo.MaintenanceState = MAINTENANCE_STATE_OFF
	}
	service, err = c.UpdateService(service)
	if err != nil {
		return nil, err
	}
	for service.ServiceInfo.MaintenanceState == MAINTENANCE_STATE_OFF && service.ServiceInfo.State != SERVICE_STARTED {
		time.Sleep(5 * time.Second)
		service, err = c.Service(service.ServiceInfo.ClusterName, service.ServiceInfo.ServiceName)
		if err != nil {
			return nil, err
		}
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
	if service.ServiceInfo.State == SERVICE_STOPPED && service.ServiceInfo.DesiredState == SERVICE_STOPPED {
		log.Debugf("The service %s is already stopped", service.ServiceInfo.ServiceName)
		return service, nil
	}

	// Stop service and wait
	service.ServiceInfo.State = SERVICE_STOPPED
	if force == true {
		service.ServiceInfo.MaintenanceState = MAINTENANCE_STATE_OFF
	}
	service, err = c.UpdateService(service)
	if err != nil {
		return nil, err
	}
	for service.ServiceInfo.MaintenanceState == MAINTENANCE_STATE_OFF && service.ServiceInfo.State != SERVICE_STOPPED {
		time.Sleep(5 * time.Second)
		service, err = c.Service(service.ServiceInfo.ClusterName, service.ServiceInfo.ServiceName)
		if err != nil {
			return nil, err
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
	path := fmt.Sprintf("/clusters/%s/services", cluster.Cluster.ClusterName)
	jsonData, err := json.Marshal(service)
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

	// Wait all service is stopped
	currentServicesCheck := make([]Service, len(cluster.Services))
	copy(currentServicesCheck, cluster.Services)
	serviceToCheck := make([]Service, 0, 1)
	for len(currentServicesCheck) > 0 {
		for _, service := range currentServicesCheck {
			serviceTemp, err := c.Service(cluster.Cluster.ClusterName, service.ServiceInfo.ServiceName)
			if err != nil {
				return err
			}
			if serviceTemp.ServiceInfo.State != SERVICE_STOPPED {
				serviceToCheck = append(serviceToCheck, service)
				log.Debugf("Wait service %s is stopped", service.ServiceInfo.ServiceName)
			} else {
				log.Infof("Service %s is stopped", service.ServiceInfo.ServiceName)
			}
		}

		currentServicesCheck = make([]Service, len(serviceToCheck))
		copy(currentServicesCheck, serviceToCheck)
		serviceToCheck = make([]Service, 0, 1)
		time.Sleep(10 * time.Second)
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
	path := fmt.Sprintf("/clusters/%s/services", cluster.Cluster.ClusterName)
	jsonData, err := json.Marshal(service)
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

	// Wait all service is sstarted
	currentServicesCheck := make([]Service, len(cluster.Services))
	copy(currentServicesCheck, cluster.Services)
	serviceToCheck := make([]Service, 0, 1)
	for len(currentServicesCheck) > 0 {
		for _, service := range currentServicesCheck {
			serviceTemp, err := c.Service(cluster.Cluster.ClusterName, service.ServiceInfo.ServiceName)
			if err != nil {
				return err
			}
			if serviceTemp.ServiceInfo.State != SERVICE_STARTED {
				serviceToCheck = append(serviceToCheck, service)
				log.Debugf("Wait service %s is started", service.ServiceInfo.ServiceName)
			} else {
				log.Infof("Service %s is started", service.ServiceInfo.ServiceName)
			}
		}

		currentServicesCheck = make([]Service, len(serviceToCheck))
		copy(currentServicesCheck, serviceToCheck)
		serviceToCheck = make([]Service, 0, 1)
		time.Sleep(10 * time.Second)
	}

	return nil

}

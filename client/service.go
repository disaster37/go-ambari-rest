package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/services-service.md

type Service struct {
	ServiceInfo       *ServiceInfo       `json:"ServiceInfo"`
	ServiceComponents []ServiceComponent `json:"components,omitempty"`
}

type ServiceInfo struct {
	ClusterName  string  `json:"cluster_name,omitempty"`
	ServiceName  string  `json:"service_name,omitempty"`
	State        string  `json:"state,omitempty"`
	RepositoryId int     `json:"desired_repository_version_id,omitempty"`
	DesiredState *string `json:"desired_state,omitempty"`
}

type ServiceComponent struct {
	ServiceComponentInfo *ServiceComponentInfo `json:"ServiceComponentInfo"`
}

type ServiceComponentInfo struct {
	ClusterName   string `json:"cluster_name,omitempty"`
	ServiceName   string `json:"service_name,omitempty"`
	ComponentName string `json:"component_name,omitempty"`
}

func (s *Service) String() string {
	json, _ := json.Marshal(s)
	return string(json)
}

func (s *Service) clean() {
	s.ServiceComponents = nil
	s.ServiceInfo.DesiredState = nil
}

func (c *AmbariClient) CreateService(service *Service) (*Service, error) {

	if service == nil {
		panic("Service can't be nil")
	}

	log.Debug("Service: %s", service.String())
	service.clean()
	service.ServiceInfo.State = "INIT"

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
	if service == nil {
		return nil, NewAmbariError(500, "Can't get service that just created")
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
	service.clean()

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

	service.ServiceInfo.State = "INSTALLED"

	service, err := c.UpdateService(service)
	return service, err
}

func (c *AmbariClient) StartService(service *Service) (*Service, error) {
	if service == nil {
		panic("Service can't be nil")
	}
	log.Debug("Service: ", service)

	service.ServiceInfo.State = "STARTED"

	service, err := c.UpdateService(service)
	return service, err
}

func (c *AmbariClient) StopService(service *Service) (*Service, error) {
	if service == nil {
		panic("Service can't be nil")
	}
	log.Debug("Service: ", service)

	service.ServiceInfo.State = "INSTALLED"

	service, err := c.UpdateService(service)
	return service, err
}

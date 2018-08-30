package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/host-resources.md

type Host struct {
	HostInfo *HostInfo `json:"Hosts"`
}

type HostInfo struct {
	ClusterName      string          `json:"cluster_name,omitempty"`
	Hostname         string          `json:"host_name,omitempty"`
	MaintenanceState string          `json:"maintenance_state,omitempty"`
	Rack             string          `json:"rack_info,omitempty"`
	HostComponents   []HostComponent `json:"host_components,omitempty"`
}

type HostBlueprint struct {
	Blueprint string `json:"blueprint,omitempty"`
	HostGroup string `json:"host_group,omitempty"`
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
	path := fmt.Sprintf("/clusters/%s/hosts/%s", host.HostInfo.ClusterName, host.HostInfo.Hostname)
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
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	host, err = c.UpdateHost(host)
	if err != nil {
		return nil, err
	}

	host, err = c.HostOnCluster(host.HostInfo.ClusterName, host.HostInfo.Hostname)
	if err != nil {
		return nil, err
	}
	if host == nil {
		return nil, NewAmbariError(500, "Can't get host that just created")
	}

	log.Debug("Return host: %s", host)

	return host, nil

}

func (c *AmbariClient) HostOnCluster(clusterName string, hostname string) (*Host, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("HostName can't be empty")
	}

	path := fmt.Sprintf("/clusters/%s/hosts/%s", clusterName, hostname)

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
	host := &Host{}
	err = json.Unmarshal(resp.Body(), host)
	if err != nil {
		return nil, err
	}
	log.Debug("Return host: %s", host)

	return host, nil
}

func (c *AmbariClient) Host(hostname string) (*Host, error) {

	if hostname == "" {
		panic("HostName can't be empty")
	}

	path := fmt.Sprintf("/hosts/%s", hostname)

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
	host := &Host{}
	err = json.Unmarshal(resp.Body(), host)
	if err != nil {
		return nil, err
	}
	log.Debug("Return host: %s", host)

	return host, nil
}

func (c *AmbariClient) UpdateHost(host *Host) (*Host, error) {

	if host == nil {
		panic("Host can't be nil")
	}
	log.Debug("Host: ", host)

	path := fmt.Sprintf("/clusters/%s/hosts/%s", host.HostInfo.ClusterName, host.HostInfo.Hostname)
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
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	// Get the Host
	host, err = c.HostOnCluster(host.HostInfo.ClusterName, host.HostInfo.Hostname)
	if err != nil {
		return nil, err
	}
	if host == nil {
		return nil, NewAmbariError(500, "Can't get host that just updated")
	}

	log.Debug("Return host: %s", host.String())

	return host, err

}

func (c *AmbariClient) DeleteHost(clusterName string, hostname string) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("Hostname can't be empty")
	}
	path := fmt.Sprintf("/clusters/%s/hosts/%s", clusterName, hostname)

	resp, err := c.Client().R().Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response to delete host: ", resp)
	if resp.StatusCode() >= 300 {
		return NewAmbariError(resp.StatusCode(), resp.Status())
	}

	return nil

}

// Permit to register new host on cluster and associate it to existant host group in blueprint
func (c *AmbariClient) RegisterHostOnCluster(clusterName string, hostname string, blueprintName string, role string) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("Hostname can't be empty")
	}
	if blueprintName == "" {
		panic("BlueprintName can't be empty")
	}
	if role == "" {
		panic("Role can't be empty")
	}

	// Check if host exist
	host, err := c.Host(hostname)
	if err != nil {
		return err
	}
	if host == nil {
		return NewAmbariError(404, "Host %s not found", hostname)
	}
	log.Debugf("Host %s found", hostname)

	// Check if cluster exist
	cluster, err := c.Cluster(clusterName)
	if err != nil {
		return err
	}
	if cluster == nil {
		return NewAmbariError(404, "Cluster %s not found", clusterName)
	}
	log.Debugf("Cluster %s found", clusterName)

	// Check if blueprint exit
	blueprint, err := c.Blueprint(blueprintName)
	if err != nil {
		return nil
	}
	if blueprint == nil {
		return NewAmbariError(404, "Blueprint %s not found", blueprintName)
	}
	log.Debugf("Blueprint %s found", blueprintName)

	// Check if role exist on blueprint
	hostGroupFound := false
	for _, hostGroup := range blueprint.HostGroups {
		if hostGroup.Name == role {
			hostGroupFound = true
			break
		}
	}
	if hostGroupFound == false {
		return NewAmbariError(404, "Role %s not found in blueprint %s", role, blueprintName)
	}
	log.Debugf("Role %s found in blueprint %s", role, blueprintName)

	// Associate host to blueprint role
	path := fmt.Sprintf("/clusters/%s/hosts/%s", clusterName, hostname)
	data := &HostBlueprint{
		Blueprint: blueprintName,
		HostGroup: role,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := c.Client().R().SetBody(jsonData).Post(path)
	if err != nil {
		return err
	}
	log.Debug("Response to update: ", resp)
	if resp.StatusCode() >= 300 {
		return NewAmbariError(resp.StatusCode(), resp.Status())
	}

	return nil

}

func (c *AmbariClient) StopAllComponents(clusterName string, hostname string) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("Hostname can't be empty")
	}

	// Check if host exist
	host, err := c.HostOnCluster(clusterName, hostname)
	if err != nil {
		return err
	}
	if host == nil {
		return NewAmbariError(404, "Host %s not found in cluster %s", hostname, clusterName)
	}
	log.Debugf("Host %s found in cluster %s", hostname, clusterName)

	// Loop over host components to stop them
	for _, hostRole := range host.HostInfo.HostComponents {
		hostComponent, err := c.HostComponent(clusterName, hostname, hostRole.HostComponentInfo.ComponentName)
		if err != nil {
			return err
		}

		_, err = c.StopHostComponent(hostComponent)
		if err != nil {
			return err
		}
	}

	return nil

}

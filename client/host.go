// This file permit to manage host in Ambari API
// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/host-resources.md

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Host object
type Host struct {
	HostInfo       *HostInfo       `json:"Hosts"`
	HostComponents []HostComponent `json:"host_components,omitempty"`
}
type Hosts struct {
	Items []Host `json:"items,omitempty"`
}
type HostInfo struct {
	ClusterName      string `json:"cluster_name,omitempty"`
	Hostname         string `json:"host_name,omitempty"`
	MaintenanceState string `json:"maintenance_state,omitempty"`
	Rack             string `json:"rack_info,omitempty"`
}
type HostBlueprint struct {
	Blueprint string `json:"blueprint,omitempty"`
	HostGroup string `json:"host_group,omitempty"`
}

// String return host oject as Json string
func (h *Host) String() string {
	json, _ := json.Marshal(h)
	return string(json)
}

// CleanBeforeSave permit to remove some attribute before save or update host
func (h *Host) CleanBeforeSave() {
	h.HostComponents = make([]HostComponent, 0, 0)
}

// CreateHost permit to create host (attach existing Ambari host on existing cluster)
// It return host if all work fine
// It return error if something wrong when it call the API
func (c *AmbariClient) CreateHost(host *Host) (*Host, error) {

	if host == nil {
		panic("Host can't be nil")
	}
	log.Debug("Host: %s", host.String())

	host.CleanBeforeSave()
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

// HostOnCluster permit to get host from cluster
// It return host object if host is found on cluster
// It return nil if host not found on cluster
// It return error if somethink wrong when it cal the API
func (c *AmbariClient) HostOnCluster(clusterName string, hostname string) (*Host, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("HostName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("Hostname: ", hostname)

	path := fmt.Sprintf("/clusters/%s/hosts/%s", clusterName, hostname)
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

// HostsOnCluster permit to get all hosts in cluster
// It return slice of host (the slice can't be empty if there are no host)
// It return error if something wrong in API call
func (c *AmbariClient) HostsOnCluster(clusterName string) ([]Host, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)

	path := fmt.Sprintf("/clusters/%s/hosts", clusterName)
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
	hosts := &Hosts{}
	err = json.Unmarshal(resp.Body(), hosts)
	if err != nil {
		return nil, err
	}
	log.Debug("Return host: %s", hosts)

	return hosts.Items, nil
}

// Host permit to get host from hostname
// It return host if is found
// It return nil if is not found
// It return error if something wrong when it call the API
func (c *AmbariClient) Host(hostname string) (*Host, error) {

	if hostname == "" {
		panic("HostName can't be empty")
	}
	log.Debug("Hostname: ", hostname)

	path := fmt.Sprintf("/hosts/%s", hostname)
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

// Hosts permit to get all ambari agent hosts
// It return slice of hosts (slice can be empty if there are no ambari agent)
// It return error if something wrong when it call the API
func (c *AmbariClient) Hosts() ([]Host, error) {

	path := fmt.Sprintf("/hosts")

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
	hosts := &Hosts{}
	err = json.Unmarshal(resp.Body(), hosts)
	if err != nil {
		return nil, err
	}
	log.Debug("Return host: %s", hosts)

	return hosts.Items, nil
}

// UpdateHost permit to update host like maintenance state
// It return updated host objcet if all work fine
// It return error if something wrong when it call the API
func (c *AmbariClient) UpdateHost(host *Host) (*Host, error) {

	if host == nil {
		panic("Host can't be nil")
	}
	log.Debug("Host: ", host)

	host.CleanBeforeSave()
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

// Delete host permit to delete host on clusterS
// It need to delete all component hosted on host before to delete the host
func (c *AmbariClient) DeleteHost(clusterName string, hostname string) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("Hostname can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("Hostname: ", hostname)

	// Check if host exist on cluster
	host, err := c.HostOnCluster(clusterName, hostname)
	if err != nil {
		return nil
	}
	if host == nil {
		return NewAmbariError(404, "Host %s not found in cluster %s", hostname, clusterName)
	}

	// Stop All components hosted in host before delete it
	err = c.StopAllComponentsInHost(clusterName, hostname, false, true)
	if err != nil {
		return nil
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

// RegisterHostOnCluster permit to register new host on cluster and associate it to existant host group in blueprint
// It return host if all work fine
// It return error if host is not found by Ambari or if cluster is not found or if blueprint is not found or if role not exist in blueprint or something wrong when it call the API
func (c *AmbariClient) RegisterHostOnCluster(clusterName string, hostname string, blueprintName string, role string) (*Host, error) {

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
	log.Debug("ClusterName: ", clusterName)
	log.Debug("Hostname: ", hostname)
	log.Debug("BlueprintName: ", blueprintName)
	log.Debug("Role: ", role)

	// Check if host exist
	host, err := c.Host(hostname)
	if err != nil {
		return nil, err
	}
	if host == nil {
		return nil, NewAmbariError(404, "Host %s not found", hostname)
	}
	log.Debugf("Host %s found", hostname)

	// Check if cluster exist
	cluster, err := c.Cluster(clusterName)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, NewAmbariError(404, "Cluster %s not found", clusterName)
	}
	log.Debugf("Cluster %s found", clusterName)

	// Check if blueprint exit
	blueprint, err := c.Blueprint(blueprintName)
	if err != nil {
		return nil, err
	}
	if blueprint == nil {
		return nil, NewAmbariError(404, "Blueprint %s not found", blueprintName)
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
		return nil, NewAmbariError(404, "Role %s not found in blueprint %s", role, blueprintName)
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
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to update: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	// Finnaly load the host
	host, err = c.HostOnCluster(clusterName, hostname)
	if err != nil {
		return nil, err
	}

	return host, nil

}

// StopAllComponentsInHost stop all components in host in arbitrary order
// if enableMaintenanceMode is set to true, it will enable maintenance state in host after stop all ressources
// if force is set to true, it will remove maintenance state in host before stop all ressources
func (c *AmbariClient) StopAllComponentsInHost(clusterName string, hostname string, enableMaintenanceMode bool, force bool) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("Hostname can't be empty")
	}

	log.Debug("ClusterName: ", clusterName)
	log.Debug("Hostname: ", hostname)
	log.Debug("EnableMaintenanceMode: ", enableMaintenanceMode)
	log.Debug("Force: ", force)

	// Check if host exist
	host, err := c.HostOnCluster(clusterName, hostname)
	if err != nil {
		return err
	}
	if host == nil {
		return NewAmbariError(404, "Host %s not found in cluster %s", hostname, clusterName)
	}
	log.Debugf("Host %s found in cluster %s", hostname, clusterName)

	// Disable maintenance state in host if needed
	if force == true && host.HostInfo.MaintenanceState != MAINTENANCE_STATE_OFF {
		host.HostInfo.MaintenanceState = MAINTENANCE_STATE_OFF
		host, err = c.UpdateHost(host)
		if err != nil {
			return err
		}
		log.Debugf("Maintenace state is disable on host %s", hostname)
	}

	// Stop all components in host and wait
	for _, hostComponent := range host.HostComponents {
		_, err := c.StopHostComponent(clusterName, hostname, hostComponent.HostComponentInfo.ComponentName)
		if err != nil {
			return err
		}
		log.Infof("Component %s is stopped", hostComponent.HostComponentInfo.ComponentName)
	}

	// Enable host maintenance if needed
	if enableMaintenanceMode == true {
		host.HostInfo.MaintenanceState = MAINTENANCE_STATE_ON
		host, err = c.UpdateHost(host)
		if err != nil {
			return err
		}

		log.Debugf("Maintenace state is enable on host %s", hostname)
	}

	return nil

}

// StartAllComponentsInHost permit to start all components in arbitrary order
// If disableMaintenanceMode is set to true, it will remove maintenance mode on host before start all components
// If maintenanceState is set to on, it will do nothink
// It return error if host is not found or something wrong when it call the API
func (c *AmbariClient) StartAllComponentsInHost(clusterName string, hostname string, disableMaintenanceMode bool) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("Hostname can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	log.Debug("Hostname: ", hostname)
	log.Debug("DisableMaintenanceMode: ", disableMaintenanceMode)

	// Check if host exist
	host, err := c.HostOnCluster(clusterName, hostname)
	if err != nil {
		return err
	}
	if host == nil {
		return NewAmbariError(404, "Host %s not found in cluster %s", hostname, clusterName)
	}
	log.Debugf("Host %s found in cluster %s", hostname, clusterName)

	// Disable maintenance state in host if needed
	if disableMaintenanceMode == true && host.HostInfo.MaintenanceState != MAINTENANCE_STATE_OFF {
		host.HostInfo.MaintenanceState = MAINTENANCE_STATE_OFF
		host, err = c.UpdateHost(host)
		if err != nil {
			return err
		}
		log.Debugf("Maintenace state is disable on host %s", hostname)
	}

	// Start all components in host and wait
	for _, hostComponent := range host.HostComponents {

		_, err = c.StartHostComponent(clusterName, hostname, hostComponent.HostComponentInfo.ComponentName)
		if err != nil {
			return err
		}
		log.Infof("Component %s is started", hostComponent.HostComponentInfo.ComponentName)
	}

	return nil

}

// DeleteAllComponentsInHost stop and delete all components in host in arbitrary order
// if disableMaintenanceMode is set to true, it will remove maintenance state in host before to delete all ressources
func (c *AmbariClient) DeleteAllComponentsInHost(clusterName string, hostname string, disableMaintenanceMode bool) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if hostname == "" {
		panic("Hostname can't be empty")
	}

	log.Debug("ClusterName: ", clusterName)
	log.Debug("Hostname: ", hostname)
	log.Debug("DisableMaintenanceMode: ", disableMaintenanceMode)

	// Check if host exist
	host, err := c.HostOnCluster(clusterName, hostname)
	if err != nil {
		return err
	}
	if host == nil {
		return NewAmbariError(404, "Host %s not found in cluster %s", hostname, clusterName)
	}
	log.Debugf("Host %s found in cluster %s", hostname, clusterName)

	// Disable maintenance state in host if needed
	if disableMaintenanceMode == true && host.HostInfo.MaintenanceState != MAINTENANCE_STATE_OFF {
		host.HostInfo.MaintenanceState = MAINTENANCE_STATE_OFF
		host, err = c.UpdateHost(host)
		if err != nil {
			return err
		}
		log.Debugf("Maintenace state is disable on host %s", hostname)
	}

	// Stop and delete all components in host and wait
	for _, hostComponent := range host.HostComponents {
		_, err := c.StopHostComponent(clusterName, hostname, hostComponent.HostComponentInfo.ComponentName)
		if err != nil {
			return err
		}
		log.Infof("Component %s is stopped", hostComponent.HostComponentInfo.ComponentName)
		err = c.DeleteHostComponent(clusterName, hostname, hostComponent.HostComponentInfo.ComponentName)
		if err != nil {
			return err
		}
		log.Infof("Component %s is deleted", hostComponent.HostComponentInfo.ComponentName)
	}

	return nil

}

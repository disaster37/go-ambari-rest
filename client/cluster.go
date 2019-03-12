// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/cluster-resources.md
// This file permit to manager cluster item on Ambari API

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Cluster item
type Cluster struct {
	ClusterInfo       *ClusterInfo                 `json:"Clusters"`
	Services          []Service                    `json:"services,omitempty"`
	DesiredConfigs    map[string]Configuration     `json:"desired_config,omitempty"`
	SessionAttributes map[string]map[string]string `json:"session_attributes,omitempty"`
}
type ClusterInfo struct {
	ClusterId    int64  `json:"cluster_id,omitempty"`
	ClusterName  string `json:"cluster_name"`
	Version      string `json:"version,omitempty"`
	SecurityType string `json:"security_type,omitempty"`
}

// String permit to return cluster object as Json string
func (c *Cluster) String() string {
	json, _ := json.Marshal(c)
	return string(json)
}

func (c *Cluster) CleanBeforeSave() {
	c.Services = nil
	c.ClusterInfo = &ClusterInfo{
		SecurityType: c.ClusterInfo.SecurityType,
		ClusterName:  c.ClusterInfo.ClusterName,
	}
}

// Create cluster eprmit to create new HDP cluster on Ambari
// It return the cluster object if all work fine
// It return error if something wrong when it call the API
func (c *AmbariClient) CreateCluster(cluster *Cluster) (*Cluster, error) {

	if cluster == nil {
		panic("Cluster can't be nil")
	}
	log.Debug("Cluster: ", cluster)

	// Create the Cluster
	path := fmt.Sprintf("/clusters/%s", cluster.ClusterInfo.ClusterName)
	jsonData, err := json.Marshal(cluster)
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

	// Get the cluster
	cluster, err = c.Cluster(cluster.ClusterInfo.ClusterName)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, NewAmbariError(500, "Can't get cluster that just created")
	}

	return cluster, err

}

// CreateClusterFromTemplate permit to create new cluster from template file. It use the blueprint object to create automatically a cluster with topologie.
// It return the cluster object if all work fine
// It return error if something wrong when it call the API
func (c *AmbariClient) CreateClusterFromTemplate(name string, jsonClusterTemplate string) (*Cluster, error) {

	if name == "" {
		panic("Name can't be empty")
	}
	if jsonClusterTemplate == "" {
		panic("JsonClusterTemplate can't be empty")
	}
	var clusterJson interface{}
	err := json.Unmarshal([]byte(jsonClusterTemplate), &clusterJson)
	if err != nil {
		return nil, err
	}
	log.Debugf("Name: %s", name)
	log.Debugf("JsonClusterTemplate: %s", jsonClusterTemplate)

	// Create the Cluster
	path := fmt.Sprintf("/clusters/%s", name)
	resp, err := c.Client().R().SetBody(jsonClusterTemplate).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	// Get the cluster
	cluster, err := c.Cluster(name)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, NewAmbariError(500, "Can't get cluster that just created")
	}

	return cluster, err
}

// Cluster permit to return cluster object from is name
// It return cluster object if found
// It return nil if cluster is not found
// It return error if something wrong when it call the API
func (c *AmbariClient) Cluster(clusterName string) (*Cluster, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	path := fmt.Sprintf("/clusters/%s", clusterName)

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
	cluster := &Cluster{}
	err = json.Unmarshal(resp.Body(), cluster)
	if err != nil {
		return nil, err
	}
	log.Debug("Cluster: ", cluster)

	return cluster, nil
}

// Ambari not support to manage cluster by ID. We need to use clusterName instead.
// So we need to have the old cluster name is the goal to rename it.
// It return cluster if all right fine
// It return error if something wrong when it call the API
func (c *AmbariClient) RenameCluster(oldClusterName string, cluster *Cluster) (*Cluster, error) {

	if oldClusterName == "" {
		panic("OldClusterName can't be nil")
	}
	if cluster == nil {
		panic("Cluster can't be nil")
	}
	log.Debug("OldClusterName: ", oldClusterName)
	log.Debug("Cluster: ", cluster)

	// Update the Cluster
	path := fmt.Sprintf("/clusters/%s", oldClusterName)
	cluster = &Cluster{
		ClusterInfo: &ClusterInfo{
			ClusterName: cluster.ClusterInfo.ClusterName,
		},
	}
	jsonData, err := json.Marshal(cluster)
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

	// Get the cluster
	cluster, err = c.Cluster(cluster.ClusterInfo.ClusterName)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, NewAmbariError(500, "Can't get cluster that just updated")
	}

	log.Debug("Cluster: ", cluster)

	return cluster, err

}

// Ambari not support to manage cluster by ID. We need to use clusterName instead.
// It return cluster if all right fine
// It return error if something wrong when it call the API
func (c *AmbariClient) ManageKerberosOnCluster(cluster *Cluster) (*Cluster, error) {

	if cluster == nil {
		panic("Cluster can't be nil")
	}
	log.Debug("Cluster: ", cluster)

	context := "Disable kerberos from API"
	if cluster.ClusterInfo.SecurityType == "KERBEROS" {
		context = "Enable kerberos from API"
	}

	request := &Request{
		RequestInfo: &RequestInfo{
			Context: context,
		},
		Body: cluster,
	}
	requestTask, err := c.SendRequestCluster(request)
	if err != nil {
		return nil, err
	}

	// Wait the end of the request if request has been created
	if requestTask != nil {

		// Wait the end of the request
		err = requestTask.Wait(c, cluster.ClusterInfo.ClusterName)
		if err != nil {
			return nil, err
		}

		// Check the status
		if requestTask.RequestTaskInfo.Status != REQUEST_COMPLETED {
			return nil, NewAmbariError(500, "Request %d failed with status %s, task completed %d, task aborded %d, task failed %d", requestTask.RequestTaskInfo.Id, requestTask.RequestTaskInfo.Status, requestTask.RequestTaskInfo.CompletedTask, requestTask.RequestTaskInfo.AbordedTask, requestTask.RequestTaskInfo.FailedTask)
		}
	}

	// Finnaly get the cluster
	cluster, err = c.Cluster(cluster.ClusterInfo.ClusterName)
	if err != nil {
		return nil, err
	}

	return cluster, err

}

// DeleteCluster permit to delete existing cluster
// It need to delete all services and delete all hosts before to delete the cluster
// It return error if cluster not exist of something wrong when it call the API
func (c *AmbariClient) DeleteCluster(clusterName string) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)

	// Check if cluster exist
	cluster, err := c.Cluster(clusterName)
	if err != nil {
		return err
	}
	if cluster == nil {
		return NewAmbariError(404, "Cluster %s not found", clusterName)
	}

	// Force to stop all services
	err = c.StopAllServices(cluster, false, true)
	if err != nil {
		return err
	}

	// Remove all services
	for _, service := range cluster.Services {
		err = c.DeleteService(service.ServiceInfo.ClusterName, service.ServiceInfo.ServiceName)
		if err != nil {
			return err
		}
	}

	// Remove all hosts
	hosts, err := c.HostsOnCluster(clusterName)
	if err != nil {
		return err
	}
	for _, host := range hosts {
		err := c.DeleteHost(clusterName, host.HostInfo.Hostname)
		if err != nil {
			return err
		}
	}

	path := fmt.Sprintf("/clusters/%s", clusterName)
	resp, err := c.Client().R().Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response to delete cluster: ", resp)
	if resp.StatusCode() >= 300 {
		return NewAmbariError(resp.StatusCode(), resp.Status())
	}

	return nil

}

// SendRequestCluster permit to start / stop service on ambari with message displayed on operation task in Ambari UI
// It keep only Kerberos field when it send the request
// It return RequestTask if all work fine
// It return nil if no request is created
// It return error if something wrong when it call the API
func (c *AmbariClient) SendRequestCluster(request *Request) (*RequestTask, error) {

	if request == nil {
		panic("Request can't be nil")
	}
	log.Debug("Request: ", request)
	cluster := request.Body.(*Cluster)
	clusterTemp := &Cluster{
		ClusterInfo: &ClusterInfo{
			SecurityType: cluster.ClusterInfo.SecurityType,
			ClusterName:  cluster.ClusterInfo.ClusterName,
		},
		SessionAttributes: cluster.SessionAttributes,
	}
	request.Body = clusterTemp

	log.Debug("Sended Request: ", request)

	path := fmt.Sprintf("/clusters/%s", cluster.ClusterInfo.ClusterName)
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

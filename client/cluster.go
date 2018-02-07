package client

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/cluster-resources.md

type Cluster struct {
	Response
	Cluster *ClusterInfo `json:"Clusters"`
}

type ClusterInfo struct {
	ClusterId   int64  `json:"cluster_id,omitempty"`
	ClusterName string `json:"cluster_name"`
	Version     string `json:"version,omitempty"`
}

func (c *Cluster) String() string {
	json, _ := json.Marshal(c)
	return string(json)
}

func (c *AmbariClient) CreateCluster(cluster *Cluster) (*Cluster, error) {

	if cluster == nil {
		panic("Cluster can't be nil")
	}

	// Create the Cluster
	path := fmt.Sprintf("/clusters/%s", cluster.Cluster.ClusterName)
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
		return nil, errors.New(resp.Status())
	}

	// Get the cluster
	cluster, err = c.Cluster(cluster.Cluster.ClusterName)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, errors.New("Can't get cluster that just created")
	}

	return cluster, err

}

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
	log.Debug("Name: %s", name)
	log.Debug("JsonClusterTemplate: %s", jsonClusterTemplate)

	// Create the Cluster
	path := fmt.Sprintf("/clusters/%s", name)
	resp, err := c.Client().R().SetBody(jsonClusterTemplate).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, errors.New(resp.Status())
	}

	// Get the cluster
	cluster, err := c.Cluster(name)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, errors.New("Can't get cluster that just created")
	}

	return cluster, err
}

// Get cluster by ID is not supported by ambari api
func (c *AmbariClient) Cluster(clusterName string) (*Cluster, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	path := fmt.Sprintf("/clusters/%s", clusterName)

	// Get the privilege
	resp, err := c.Client().R().Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to get: ", resp)
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
// CHange AD version by this way is not supported. We need to use upgrade API to to that.
func (c *AmbariClient) UpdateCluster(oldClusterName string, cluster *Cluster) (*Cluster, error) {

	if oldClusterName == "" {
		panic("OldClusterName can't be nil")
	}
	if cluster == nil {
		panic("Cluster can't be nil")
	}

	log.Debug("Cluster: ", cluster)

	// Update the Cluster
	path := fmt.Sprintf("/clusters/%s", oldClusterName)
	cluster = &Cluster{
		Cluster: &ClusterInfo{
			ClusterName: cluster.Cluster.ClusterName,
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
		return nil, errors.New(resp.Status())
	}

	// Get the cluster
	cluster, err = c.Cluster(cluster.Cluster.ClusterName)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, errors.New("Can't get cluster that just updated")
	}

	log.Debug("Cluster: ", cluster)

	return cluster, err

}

func (c *AmbariClient) DeleteCluster(clusterName string) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}

	path := fmt.Sprintf("/clusters/%s", clusterName)
	resp, err := c.Client().R().Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response to delete cluster: ", resp)
	if resp.StatusCode() >= 300 {
		return errors.New(resp.Status())
	}

	return nil

}

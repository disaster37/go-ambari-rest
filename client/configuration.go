// This file permit to manage configuratiopn from Ambari API
// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/configuration.md

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Object item
type Configuration struct {
	Type       string            `json:"type,omitempty"`
	Tag        string            `json:"tag,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}
type DesiredConfig struct {
	DesiredConfig *Configuration `json:"desired_config,omitempty"`
}
type RequestAddConfig struct {
	Cluster *DesiredConfig `json:"Clusters,omitempty"`
}

// String permit to return Configuration as Json string
func (c *Configuration) String() string {
	json, _ := json.Marshal(c)
	return string(json)
}

// CreateConfigurationOnCluster permit to add new service confoguration on cluster
// It return cluster object if all right fine
// It return error if something wrong
func (c *AmbariClient) CreateConfigurationOnCluster(clusterName string, configuration *Configuration) (*Cluster, error) {
	if clusterName == "" {
		panic("ClusterName can't be empty")
	}

	if configuration == nil {
		panic("Configuration can't be empty")
	}

	log.Debugf("ClusterName: %s", clusterName)
	log.Debugf("Configuration: %s", configuration)

	// Create the configuration
	path := fmt.Sprintf("/clusters/%s", clusterName)
	data := &RequestAddConfig{
		Cluster: &DesiredConfig{
			DesiredConfig: configuration,
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Put(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	// Get the cluster
	cluster, err := c.Cluster(clusterName)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, NewAmbariError(500, "Can't get cluster that just updated")
	}

	return cluster, err

}

package client

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

// Test the constructor
func (s *ClientTestSuite) TestCluster() {

	// Create cluster
	cluster := &Cluster{
		Cluster: &ClusterInfo{
			Version:     "HDP-2.6",
			ClusterName: "test2",
		},
	}
	cluster, err := s.client.CreateCluster(cluster)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cluster)
	if cluster != nil {
		assert.NotEqual(s.T(), "", cluster.Cluster.ClusterId)
		assert.Equal(s.T(), "test2", cluster.Cluster.ClusterName)
		assert.Equal(s.T(), "HDP-2.6", cluster.Cluster.Version)
	}

	// Get cluster
	cluster, err = s.client.Cluster("test2")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cluster)
	if cluster != nil {
		assert.NotEqual(s.T(), "", cluster.Cluster.ClusterId)
		assert.Equal(s.T(), "test2", cluster.Cluster.ClusterName)
		assert.Equal(s.T(), "HDP-2.6", cluster.Cluster.Version)
	}

	// Update cluster
	if cluster != nil {
		cluster.Cluster.ClusterName = "test3"
		cluster, err = s.client.UpdateCluster("test2", cluster)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), cluster)
		if cluster != nil {
			assert.NotEqual(s.T(), "", cluster.Cluster.ClusterId)
			assert.Equal(s.T(), "test3", cluster.Cluster.ClusterName)
			assert.Equal(s.T(), "HDP-2.6", cluster.Cluster.Version)
		}
	}

	// Delete cluster
	err = s.client.DeleteCluster("test3")
	assert.NoError(s.T(), err)

	// Create cluster with blueprint
	// Create blueprint
	b, err := ioutil.ReadFile("../fixtures/blueprint.json")
	if err != nil {
		panic(err)
	}
	blueprintJson := string(b)
	_, err = s.client.CreateBlueprint("test", blueprintJson)
	if err != nil {
		panic(err)
	}
	//Create cluster from template
	b, err = ioutil.ReadFile("../fixtures/cluster-template.json")
	if err != nil {
		panic(err)
	}
	templateJson := string(b)
	cluster, err = s.client.CreateClusterFromTemplate("test2", templateJson)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cluster)

	// Delete the cluster
	err = s.client.DeleteCluster("test2")
	assert.NoError(s.T(), err)

}

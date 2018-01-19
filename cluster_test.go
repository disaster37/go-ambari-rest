package client

import (
	"github.com/stretchr/testify/assert"
)

// Test the constructor
func (s *ClientTestSuite) TestCluster() {

	// Create cluster
	cluster := &Cluster{
		Cluster: &ClusterInfo{
			Version:     "HDP-2.6",
			ClusterName: "test",
		},
	}
	cluster, err := s.client.CreateCluster(cluster)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cluster)
	if cluster != nil {
		assert.NotEqual(s.T(), "", cluster.Cluster.ClusterId)
		assert.Equal(s.T(), "test", cluster.Cluster.ClusterName)
		assert.Equal(s.T(), "HDP-2.6", cluster.Cluster.Version)
	}

	// Get cluster
	cluster, err = s.client.Cluster("test")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cluster)
	if cluster != nil {
		assert.NotEqual(s.T(), "", cluster.Cluster.ClusterId)
		assert.Equal(s.T(), "test", cluster.Cluster.ClusterName)
		assert.Equal(s.T(), "HDP-2.6", cluster.Cluster.Version)
	}

	// Update cluster
	cluster.Cluster.ClusterName = "test2"
	cluster, err = s.client.UpdateCluster("test", cluster)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cluster)
	if cluster != nil {
		assert.NotEqual(s.T(), "", cluster.Cluster.ClusterId)
		assert.Equal(s.T(), "test2", cluster.Cluster.ClusterName)
		assert.Equal(s.T(), "HDP-2.6", cluster.Cluster.Version)
	}

	// Delete cluster
	err = s.client.DeleteCluster("test2")
	assert.NoError(s.T(), err)

}

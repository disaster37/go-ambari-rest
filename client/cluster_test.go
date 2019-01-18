package client

import (
	"github.com/stretchr/testify/assert"
)

func (s *ClientTestSuite) TestCluster() {

	// Create cluster
	cluster := &Cluster{
		ClusterInfo: &ClusterInfo{
			Version:     "HDP-2.6",
			ClusterName: "test2",
		},
	}
	cluster, err := s.client.CreateCluster(cluster)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cluster)
	if cluster != nil {
		assert.NotEqual(s.T(), "", cluster.ClusterInfo.ClusterId)
		assert.Equal(s.T(), "test2", cluster.ClusterInfo.ClusterName)
		assert.Equal(s.T(), "HDP-2.6", cluster.ClusterInfo.Version)
	}

	// Get cluster
	cluster, err = s.client.Cluster("test2")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cluster)
	if cluster != nil {
		assert.NotEqual(s.T(), "", cluster.ClusterInfo.ClusterId)
		assert.Equal(s.T(), "test2", cluster.ClusterInfo.ClusterName)
		assert.Equal(s.T(), "HDP-2.6", cluster.ClusterInfo.Version)
	}

	// Rename cluster
	if cluster != nil {
		cluster.ClusterInfo.ClusterName = "test3"
		cluster, err = s.client.RenameCluster("test2", cluster)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), cluster)
		if cluster != nil {
			assert.NotEqual(s.T(), "", cluster.ClusterInfo.ClusterId)
			assert.Equal(s.T(), "test3", cluster.ClusterInfo.ClusterName)
			assert.Equal(s.T(), "HDP-2.6", cluster.ClusterInfo.Version)
		}
	}

	// Delete cluster
	err = s.client.DeleteCluster("test3")
	assert.NoError(s.T(), err)

	// Create cluster with blueprint and delete them
	// It's already tested on client_test to create fresh cluster

}

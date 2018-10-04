package client

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"time"
)

// Test the constructor
func (s *ClientTestSuite) TestHost() {

	//Wait some time that host join after cluster deletion
	time.Sleep(60)

	// Get host
	host, err := s.client.Host("ambari-agent")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), host)
	if host != nil {
		assert.Equal(s.T(), "", host.HostInfo.ClusterName)
		assert.Equal(s.T(), "ambari-agent", host.HostInfo.Hostname)
		assert.Equal(s.T(), "/default-rack", host.HostInfo.Rack)
		assert.Equal(s.T(), "", host.HostInfo.MaintenanceState)
	}

	// Create host
	host = &Host{
		HostInfo: &HostInfo{
			ClusterName: "test",
			Hostname:    "ambari-agent",
		},
	}
	host, err = s.client.CreateHost(host)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), host)
	if host != nil {
		assert.Equal(s.T(), "test", host.HostInfo.ClusterName)
		assert.Equal(s.T(), "ambari-agent", host.HostInfo.Hostname)
		assert.Equal(s.T(), "/default-rack", host.HostInfo.Rack)
		assert.Equal(s.T(), "OFF", host.HostInfo.MaintenanceState)
	}

	// Get host on cluster
	host, err = s.client.HostOnCluster("test", "ambari-agent")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), host)
	if host != nil {
		assert.Equal(s.T(), "test", host.HostInfo.ClusterName)
		assert.Equal(s.T(), "ambari-agent", host.HostInfo.Hostname)
		assert.Equal(s.T(), "/default-rack", host.HostInfo.Rack)
		assert.Equal(s.T(), "OFF", host.HostInfo.MaintenanceState)
	}

	// Update host
	if host != nil {
		host.HostInfo.MaintenanceState = "ON"
		host.HostInfo.Rack = "/B1"
		host, err = s.client.UpdateHost(host)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), host)
		if host != nil {
			assert.Equal(s.T(), "test", host.HostInfo.ClusterName)
			assert.Equal(s.T(), "ambari-agent", host.HostInfo.Hostname)
			assert.Equal(s.T(), "/B1", host.HostInfo.Rack)
			assert.Equal(s.T(), "ON", host.HostInfo.MaintenanceState)
		}
	}

	// Delete host
	if host != nil {
		host.HostInfo.MaintenanceState = "OFF"
		s.client.UpdateHost(host)
		err = s.client.DeleteHost("test", "ambari-agent")
		assert.NoError(s.T(), err)
	}

	// Create blueprint
	b, err := ioutil.ReadFile("../fixtures/blueprint.json")
	if err != nil {
		panic(err)
	}
	blueprintJson := string(b)
	_, err = s.client.CreateBlueprint("testHost", blueprintJson)
	if err != nil {
		panic(err)
	}
	//Create cluster from template
	b, err = ioutil.ReadFile("../fixtures/cluster-template.json")
	if err != nil {
		panic(err)
	}
	templateJson := string(b)

	_, err = s.client.CreateClusterFromTemplate("testHost", templateJson)
	if err != nil {
		panic(err)
	}
	_, err = s.client.RegisterHostOnCluster("testHost", "ambari-agent3", "testHost", "host_group_1")
	assert.NoError(s.T(), err)

}

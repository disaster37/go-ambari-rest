package client

// Need to add ambari-agent on containers
/*
import (
	"github.com/stretchr/testify/assert"
)

// Test the constructor
func (s *ClientTestSuite) TestHost() {

	// Create host
	host := &Host{
		HostInfo: &HostInfo{
			ClusterName: "test",
			Hostname:    "master1",
			Rack:        "/B1",
		},
	}
	host, err := s.client.CreateHost(host)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), host)
	if host != nil {
		assert.Equal(s.T(), "test", host.HostInfo.ClusterName)
		assert.Equal(s.T(), "master1", host.HostInfo.Hostname)
		assert.Equal(s.T(), "/B1", host.HostInfo.Rack)
		assert.Equal(s.T(), "OFF", host.HostInfo.MaintenanceState)
	}

	// Get host
	host, err = s.client.Host("test", "master1")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), host)
	if host != nil {
		assert.Equal(s.T(), "test", host.HostInfo.ClusterName)
		assert.Equal(s.T(), "master1", host.HostInfo.Hostname)
		assert.Equal(s.T(), "/B1", host.HostInfo.Rack)
		assert.Equal(s.T(), "OFF", host.HostInfo.MaintenanceState)
	}

	// Update host
	host.HostInfo.MaintenanceState = "ON"
	host, err = s.client.UpdateHost(host)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), host)
	if host != nil {
		assert.Equal(s.T(), "test", host.HostInfo.ClusterName)
		assert.Equal(s.T(), "master1", host.HostInfo.Hostname)
		assert.Equal(s.T(), "/B1", host.HostInfo.Rack)
		assert.Equal(s.T(), "ON", host.HostInfo.MaintenanceState)
	}

	// Delete host
	err = s.client.DeleteHost("test", "master1")
	assert.NoError(s.T(), err)

}
*/

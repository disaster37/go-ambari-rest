package client

// Need to manage service
/*
import (
	"github.com/stretchr/testify/assert"
	"time"
)

// Test the constructor
func (s *ClientTestSuite) TestHostComponent() {

	// Wait heartbeat before add host
	time.Sleep(time.Second * 30)
	host := &Host{
		HostInfo: &HostInfo{
			ClusterName: "test",
			Hostname:    "master1",
			Rack:        "/B1",
		},
	}
	host, err := s.client.CreateHost(host)
	if err != nil {
		//panic(err)
	}

	// Create hostComponent
	hostComponent := &HostComponent{
		HostComponentInfo: &HostComponentInfo{
			ClusterName:   "test",
			ComponentName: "DATANODE",
			Hostname:      "master1",
		},
	}
	hostComponent, err = s.client.CreateHostComponent(hostComponent)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), "test", hostComponent.HostComponentInfo.ClusterName)
		assert.Equal(s.T(), "DATANODE", hostComponent.HostComponentInfo.ComponentName)
		assert.Equal(s.T(), "master1", hostComponent.HostComponentInfo.Hostname)
		assert.Equal(s.T(), "INSTALLED", hostComponent.HostComponentInfo.DesiredState)
		assert.NotEqual(s.T(), "", hostComponent.HostComponentInfo.State)
	}

	// Get hostComponent
	hostComponent, err = s.client.HostComponent("test", "master1", "DATANODE")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), "test", hostComponent.HostComponentInfo.ClusterName)
		assert.Equal(s.T(), "DATANODE", hostComponent.HostComponentInfo.ComponentName)
		assert.Equal(s.T(), "master1", hostComponent.HostComponentInfo.Hostname)
		assert.Equal(s.T(), "INSTALLED", hostComponent.HostComponentInfo.DesiredState)
		assert.NotEqual(s.T(), "", hostComponent.HostComponentInfo.State)
	}

	// Update hostComponent
	hostComponent.HostComponentInfo.State = "STARTED"
	hostComponent, err = s.client.UpdateHostComponent(hostComponent)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), "test", hostComponent.HostComponentInfo.ClusterName)
		assert.Equal(s.T(), "DATANODE", hostComponent.HostComponentInfo.ComponentName)
		assert.Equal(s.T(), "master1", hostComponent.HostComponentInfo.Hostname)
		assert.NotEqual(s.T(), "", hostComponent.HostComponentInfo.DesiredState)
		assert.NotEqual(s.T(), "", hostComponent.HostComponentInfo.State)
	}

	// Delete hostComponent
	err = s.client.DeleteHostComponent("test", "master1", "DATANODE")
	assert.NoError(s.T(), err)

}
*/

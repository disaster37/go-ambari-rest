package client

import (
	"github.com/stretchr/testify/assert"
)

func (s *ClientTestSuite) TestHostComponent() {

	// Re create hostComponent ZOOKEEPER_CLIENT on ambari-agent2 that we previously remove on component test
	hostComponent := &HostComponent{
		HostComponentInfo: &HostComponentInfo{
			ClusterName:   "test",
			ComponentName: "ZOOKEEPER_CLIENT",
			Hostname:      "ambari-agent2",
			ServiceName:   "ZOOKEEPER",
		},
	}
	hostComponent, err := s.client.CreateHostComponent(hostComponent)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), "test", hostComponent.HostComponentInfo.ClusterName)
		assert.Equal(s.T(), "ZOOKEEPER_CLIENT", hostComponent.HostComponentInfo.ComponentName)
		assert.Equal(s.T(), "ambari-agent2", hostComponent.HostComponentInfo.Hostname)
		assert.Equal(s.T(), SERVICE_INSTALLED, hostComponent.HostComponentInfo.State)
	}

	// Get hostComponent
	hostComponent, err = s.client.HostComponent("test", "ambari-agent2", "ZOOKEEPER_SERVER")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), "test", hostComponent.HostComponentInfo.ClusterName)
		assert.Equal(s.T(), "ZOOKEEPER_SERVER", hostComponent.HostComponentInfo.ComponentName)
		assert.Equal(s.T(), "ambari-agent2", hostComponent.HostComponentInfo.Hostname)
		assert.Equal(s.T(), SERVICE_STARTED, hostComponent.HostComponentInfo.State)
	}

	// Update hostComponent
	hostComponent.HostComponentInfo.State = SERVICE_STOPPED
	hostComponent, err = s.client.UpdateHostComponent(hostComponent)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), "test", hostComponent.HostComponentInfo.ClusterName)
		assert.Equal(s.T(), "ZOOKEEPER_SERVER", hostComponent.HostComponentInfo.ComponentName)
		assert.Equal(s.T(), "ambari-agent2", hostComponent.HostComponentInfo.Hostname)
		assert.Equal(s.T(), SERVICE_STOPPED, hostComponent.HostComponentInfo.DesiredState)
	}
	s.WaitClusterTask()

	// Start component
	hostComponent, err = s.client.StartHostComponent("test", "ambari-agent2", "ZOOKEEPER_SERVER")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), SERVICE_STARTED, hostComponent.HostComponentInfo.State)
	}

	// Stop component
	hostComponent, err = s.client.StopHostComponent("test", "ambari-agent2", "ZOOKEEPER_SERVER")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), SERVICE_STOPPED, hostComponent.HostComponentInfo.State)
	}

	// Delete hostComponent
	err = s.client.DeleteHostComponent("test", "ambari-agent2", "ZOOKEEPER_CLIENT")
	assert.NoError(s.T(), err)
	hostComponent, err = s.client.HostComponent("test", "ambari-agent2", "ZOOKEEPER_CLIENT")
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), hostComponent)

}

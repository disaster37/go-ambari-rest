package client

import (
	"github.com/stretchr/testify/assert"
)

// Test the constructor
func (s *ClientTestSuite) TestHostComponent() {

	// Create hostComponent
	hostComponent := &HostComponent{
		HostComponentInfo: &HostComponentInfo{
			ClusterName:   "sihm-test",
			ComponentName: "DATANODE",
			HostName:      "087e327a3ec0",
		},
	}

	hostComponent, err := s.client.CreateHostComponent(hostComponent)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), "sihm-test", hostComponent.HostComponentInfo.ClusterName)
		assert.Equal(s.T(), "DATANODE", hostComponent.HostComponentInfo.ComponentName)
		assert.Equal(s.T(), "087e327a3ec0", hostComponent.HostComponentInfo.HostName)
		assert.Equal(s.T(), "INSTALLED", hostComponent.HostComponentInfo.DesiredState)
		assert.NotEqual(s.T(), "", hostComponent.HostComponentInfo.State)
	}

	// Get hostComponent
	hostComponent, err = s.client.HostComponent("sihm-test", "087e327a3ec0", "DATANODE")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), "sihm-test", hostComponent.HostComponentInfo.ClusterName)
		assert.Equal(s.T(), "DATANODE", hostComponent.HostComponentInfo.ComponentName)
		assert.Equal(s.T(), "087e327a3ec0", hostComponent.HostComponentInfo.HostName)
		assert.Equal(s.T(), "INSTALLED", hostComponent.HostComponentInfo.DesiredState)
		assert.NotEqual(s.T(), "", hostComponent.HostComponentInfo.State)
	}

	// Update hostComponent
	hostComponent.HostComponentInfo.State = "STARTED"
	hostComponent, err = s.client.UpdateHostComponent(hostComponent)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), "sihm-test", hostComponent.HostComponentInfo.ClusterName)
		assert.Equal(s.T(), "DATANODE", hostComponent.HostComponentInfo.ComponentName)
		assert.Equal(s.T(), "087e327a3ec0", hostComponent.HostComponentInfo.HostName)
		assert.NotEqual(s.T(), "", hostComponent.HostComponentInfo.DesiredState)
		assert.NotEqual(s.T(), "", hostComponent.HostComponentInfo.State)
	}

	// Delete hostComponent
	err = s.client.DeleteHostComponent("sihm-test", "087e327a3ec0", "DATANODE")
	assert.NoError(s.T(), err)

}

package client

import (
	"github.com/stretchr/testify/assert"
)

func (s *ClientTestSuite) TestComponent() {

	// Get component
	component, err := s.client.Component("test", "ZOOKEEPER", "ZOOKEEPER_SERVER")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), component)
	if component != nil {
		assert.Equal(s.T(), "test", component.ComponentInfo.ClusterName)
		assert.Equal(s.T(), "ZOOKEEPER", component.ComponentInfo.ServiceName)
		assert.Equal(s.T(), "ZOOKEEPER_SERVER", component.ComponentInfo.ComponentName)
		assert.Equal(s.T(), "MASTER", component.ComponentInfo.Category)
		assert.Equal(s.T(), SERVICE_STARTED, component.ComponentInfo.State)
	}

	// Delete component
	err = s.client.DeleteComponent("test", "ZOOKEEPER", "ZOOKEEPER_CLIENT")
	assert.NoError(s.T(), err)
	component, err = s.client.Component("test", "ZOOKEEPER", "ZOOKEEPER_CLIENT")
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), component)

	// Recreate component for other test
	component = &Component{
		ComponentInfo: &ComponentInfo{
			ComponentName: "ZOOKEEPER_CLIENT",
			ServiceName:   "ZOOKEEPER",
			ClusterName:   "test",
		},
	}
	component, err = s.client.CreateComponent(component)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), component)
	if component != nil {
		assert.Equal(s.T(), "test", component.ComponentInfo.ClusterName)
		assert.Equal(s.T(), "ZOOKEEPER", component.ComponentInfo.ServiceName)
		assert.Equal(s.T(), "ZOOKEEPER_CLIENT", component.ComponentInfo.ComponentName)
	}

}

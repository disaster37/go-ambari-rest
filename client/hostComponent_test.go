package client

/*
import (
	"github.com/stretchr/testify/assert"
)

// Test the constructor
func (s *ClientTestSuite) TestHostComponent() {

	//Afect new host on cluster
	host := &Host{
		HostInfo: &HostInfo{
			ClusterName: "test",
			Hostname:    "ambari-agent3",
			Rack:        "/B1",
		},
	}
	host, err := s.client.CreateHost(host)
	if err != nil {
		panic(err)
	}
	// Create new service
	service := &Service{
		ServiceInfo: &ServiceInfo{
			ClusterName:  "test",
			ServiceName:  "ZOOKEEPER",
			RepositoryId: 1,
		},
	}
	_, err = s.client.CreateService(service)
	if err != nil {
		panic(err)
	}
	// Create new component
	component := &ServiceComponent{
		ServiceComponentInfo: &ServiceComponentInfo{
			ClusterName:   "test",
			ServiceName:   "ZOOKEEPER",
			ComponentName: "ZOOKEEPER_CLIENT",
		},
	}
	_, err = s.client.CreateComponent(component)
	if err != nil {
		panic(err)
	}

	// Create hostComponent
	hostComponent := &HostComponent{
		HostComponentInfo: &HostComponentInfo{
			ClusterName:   "test",
			ComponentName: "ZOOKEEPER_CLIENT",
			Hostname:      "ambari-agent3",
		},
	}
	hostComponent, err = s.client.CreateHostComponent(hostComponent)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), "test", hostComponent.HostComponentInfo.ClusterName)
		assert.Equal(s.T(), "ZOOKEEPER_CLIENT", hostComponent.HostComponentInfo.ComponentName)
		assert.Equal(s.T(), "ambari-agent3", hostComponent.HostComponentInfo.Hostname)
		assert.Equal(s.T(), "INSTALLED", hostComponent.HostComponentInfo.DesiredState)
		assert.NotEqual(s.T(), "", hostComponent.HostComponentInfo.State)
	}

	// Get hostComponent
	hostComponent, err = s.client.HostComponent("test", "ambari-agent3", "ZOOKEEPER_CLIENT")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hostComponent)
	if hostComponent != nil {
		assert.Equal(s.T(), "test", hostComponent.HostComponentInfo.ClusterName)
		assert.Equal(s.T(), "ZOOKEEPER_CLIENT", hostComponent.HostComponentInfo.ComponentName)
		assert.Equal(s.T(), "ambari-agent3", hostComponent.HostComponentInfo.Hostname)
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
		assert.Equal(s.T(), "ZOOKEEPER_CLIENT", hostComponent.HostComponentInfo.ComponentName)
		assert.Equal(s.T(), "ambari-agent3", hostComponent.HostComponentInfo.Hostname)
		assert.NotEqual(s.T(), "", hostComponent.HostComponentInfo.DesiredState)
		assert.NotEqual(s.T(), "", hostComponent.HostComponentInfo.State)
	}

	// Delete hostComponent
	err = s.client.DeleteHostComponent("test", "ambari-agent3", "ZOOKEEPER_CLIENT")
	assert.NoError(s.T(), err)

}

*/

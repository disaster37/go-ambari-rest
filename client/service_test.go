package client

import (
	"github.com/stretchr/testify/assert"
)

func (s *ClientTestSuite) TestService() {

	// Create service
	service := &Service{
		ServiceInfo: &ServiceInfo{
			ClusterName:  "test",
			ServiceName:  "ZOOKEEPER",
			RepositoryId: 1,
		},
	}
	service, err := s.client.CreateService(service)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), service)
	if service != nil {
		assert.Equal(s.T(), "test", service.ServiceInfo.ClusterName)
		assert.Equal(s.T(), "ZOOKEEPER", service.ServiceInfo.ServiceName)
		assert.Equal(s.T(), "UNKNOWN", service.ServiceInfo.State)
		assert.Equal(s.T(), 0, len(service.ServiceComponents))
	}

	// Get service
	service, err = s.client.Service("test", "ZOOKEEPER")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), service)
	if service != nil {
		assert.Equal(s.T(), "test", service.ServiceInfo.ClusterName)
		assert.Equal(s.T(), "ZOOKEEPER", service.ServiceInfo.ServiceName)
		assert.Equal(s.T(), "UNKNOWN", service.ServiceInfo.State)
		assert.Equal(s.T(), 0, len(service.ServiceComponents))
	}

	// Update service
	service.ServiceInfo.State = "INSTALLED"
	service, err = s.client.UpdateService(service)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), service)

	// Install service
	service, err = s.client.InstallService(service)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), service)

	// Start service
	service, err = s.client.StartService(service)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), service)

	// Stop service
	service, err = s.client.StopService(service)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), service)

	// Delete service
	//err = s.client.DeleteService("test", "ZOOKEEPER")
	//assert.NoError(s.T(), err)

}

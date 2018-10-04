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
		assert.Equal(s.T(), SERVICE_UNKNOWN, service.ServiceInfo.State)
		assert.Equal(s.T(), 0, len(service.ServiceComponents))
	}

	// Get service
	service, err = s.client.Service("test", "ZOOKEEPER")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), service)
	if service != nil {
		assert.Equal(s.T(), "test", service.ServiceInfo.ClusterName)
		assert.Equal(s.T(), "ZOOKEEPER", service.ServiceInfo.ServiceName)
		assert.Equal(s.T(), SERVICE_UNKNOWN, service.ServiceInfo.State)
		assert.Equal(s.T(), 0, len(service.ServiceComponents))
	}

	if service != nil {

		// Delete service
		err = s.client.DeleteService("test", "ZOOKEEPER")
		assert.NoError(s.T(), err)
	}

	// To test update, install, start, stop and delete service, we need to provide some configuration
	// @TODO

}

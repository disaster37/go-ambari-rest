package client

import (
	"github.com/stretchr/testify/assert"
)

func (s *ClientTestSuite) TestService() {

	// Create service
	service := &Service{
		ServiceInfo: &ServiceInfo{
			ClusterName:  "test",
			ServiceName:  "SMARTSENSE",
			RepositoryId: 1,
		},
	}
	service, err := s.client.CreateService(service)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), service)
	if service != nil {
		assert.Equal(s.T(), "test", service.ServiceInfo.ClusterName)
		assert.Equal(s.T(), "SMARTSENSE", service.ServiceInfo.ServiceName)
		assert.Equal(s.T(), SERVICE_UNKNOWN, service.ServiceInfo.State)
		assert.Equal(s.T(), 0, len(service.Components))
	}

	// Get service
	service, err = s.client.Service("test", "SMARTSENSE")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), service)
	if service != nil {
		assert.Equal(s.T(), "test", service.ServiceInfo.ClusterName)
		assert.Equal(s.T(), "SMARTSENSE", service.ServiceInfo.ServiceName)
		assert.Equal(s.T(), SERVICE_UNKNOWN, service.ServiceInfo.State)
		assert.Equal(s.T(), 0, len(service.Components))
	}

	// Stop service
	service, err = s.client.StopService("test", "ZOOKEEPER", false, false)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), service)
	if service != nil {
		assert.Equal(s.T(), SERVICE_STOPPED, service.ServiceInfo.State)
	}

	// Start service
	service, err = s.client.StartService("test", "ZOOKEEPER", false)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), service)
	if service != nil {
		assert.Equal(s.T(), SERVICE_STARTED, service.ServiceInfo.State)
	}

	// Stop all services
	cluster, err := s.client.Cluster("test")
	if err != nil {
		panic(err)
	}
	err = s.client.StopAllServices(cluster, false, false)
	assert.NoError(s.T(), err)
	service, err = s.client.Service("test", "ZOOKEEPER")
	if err != nil {
		panic(err)
	}
	assert.Equal(s.T(), SERVICE_STOPPED, service.ServiceInfo.State)

	//Start all services
	err = s.client.StartAllServices(cluster, false)
	assert.NoError(s.T(), err)
	service, err = s.client.Service("test", "ZOOKEEPER")
	if err != nil {
		panic(err)
	}
	assert.Equal(s.T(), SERVICE_STARTED, service.ServiceInfo.State)

	// Delete service
	err = s.client.DeleteService("test", "ZOOKEEPER")
	assert.NoError(s.T(), err)
	service, err = s.client.Service("test", "ZOOKEEPER")
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), service)

	// Install service (we can't install service without add components and host components)
	// @TODO

}

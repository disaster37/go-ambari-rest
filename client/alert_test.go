package client

import (
	"github.com/stretchr/testify/assert"
)

func (s *ClientTestSuite) TestAlert() {

	// Get alerts in cluster
	alerts, err := s.client.AlertsInCluster("test")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), alerts)

	// Get alerts in service
	alerts, err = s.client.AlertsInService("test", "ZOOKEEPER")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), alerts)

	// Get Alerts in hosts
	alerts, err = s.client.AlertsInHost("test", "ambari-agent2")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), alerts)

}

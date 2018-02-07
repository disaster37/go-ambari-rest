package client

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"testing"
)

type ClientTestSuite struct {
	suite.Suite
	client *AmbariClient
}

func (s *ClientTestSuite) SetupSuite() {
	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

	s.client = New("http://ambari-server:8080/api/v1", "admin", "admin")
	s.client.DisableVerifySSL()

	// Remove cluster
	s.client.DeleteCluster("test")

	// Create freash cluster
	cluster := &Cluster{
		Cluster: &ClusterInfo{
			Version:     "HDP-2.6",
			ClusterName: "test",
		},
	}
	cluster, err := s.client.CreateCluster(cluster)
	if err != nil {
		panic(err)
	}

}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

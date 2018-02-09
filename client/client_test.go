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

	// Create repository to create cluster
	repository := &Repository{
		RepositoryVersion: &RepositoryVersion{
			Version:      "2.6.4.0",
			Name:         "HDP-2.6.4.0",
			StackName:    "HDP",
			StackVersion: "2.6",
		},
		OS: []OS{
			OS{
				OSInfo: &OSInfo{
					Type: "redhat7",
				},
				RepositoriesData: []RepositoryData{
					RepositoryData{
						RepositoryInfo: &RepositoryInfo{
							Id:      "HDP",
							Name:    "HDP",
							BaseUrl: "http://public-repo-1.hortonworks.com/HDP/centos7/2.x/updates/2.6.4.0",
						},
					},
					RepositoryData{
						RepositoryInfo: &RepositoryInfo{
							Id:      "HDP-UTILS",
							Name:    "HDP-UTILS",
							BaseUrl: "http://public-repo-1.hortonworks.com/HDP-UTILS-1.1.0.22/repos/centos7",
						},
					},
				},
			},
		},
	}
	s.client.CreateRepository(repository)

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

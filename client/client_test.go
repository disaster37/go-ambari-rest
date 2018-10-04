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

	// Remove cluster if already exist
	cluster, err := s.client.Cluster("test")
	if err != nil {
		panic(err)
	}
	if cluster != nil {
		err = s.client.DeleteCluster("test")
		if err != nil {
			panic(err)
		}
	}

	// Remove blueprint if already exist
	blueprint, err := s.client.Blueprint("test")
	if err != nil {
		panic(err)
	}
	if blueprint != nil {
		err = s.client.DeleteBlueprint("test")
		if err != nil {
			panic(err)
		}
	}

	// Create repository if not exist
	repository, err := s.client.Repository("HDP", "2.6", 1)
	if err != nil {
		panic(err)
	}
	if repository == nil {
		repository := &Repository{
			RepositoryVersion: &RepositoryVersion{
				Version:      "2.6.4.0-91",
				Name:         "HDP-2.6.4.0",
				StackName:    "HDP",
				StackVersion: "2.6",
			},
			OS: []OS{
				OS{
					OSInfo: &OSInfo{
						Type:              "redhat7",
						ManagedRepository: true,
					},
					RepositoriesData: []RepositoryData{
						RepositoryData{
							RepositoryInfo: &RepositoryInfo{
								Id:      "HDP-2.6.4.0",
								Name:    "HDP",
								BaseUrl: "http://public-repo-1.hortonworks.com/HDP/centos7/2.x/updates/2.6.4.0",
							},
						},
						RepositoryData{
							RepositoryInfo: &RepositoryInfo{
								Id:      "HDP-UTILS-1.1.0.22",
								Name:    "HDP-UTILS",
								BaseUrl: "http://public-repo-1.hortonworks.com/HDP-UTILS-1.1.0.22/repos/centos7",
							},
						},
					},
				},
			},
		}
		_, err = s.client.CreateRepository(repository)
		if err != nil {
			panic(err)
		}
	}

	// Create freash cluster
	cluster = &Cluster{
		Cluster: &ClusterInfo{
			Version:     "HDP-2.6",
			ClusterName: "test",
		},
	}
	_, err = s.client.CreateCluster(cluster)
	if err != nil {
		panic(err)
	}

	// Add minimal config for cluster
	config := &Configuration{
		Tag:  "INITIAL",
		Type: "cluster-env",
		Properties: map[string]string{
			"cluster_name": "test",
			"stack_id":     "HDP-2.6",
		},
	}
	_, err = s.client.CreateConfigurationOnCluster("test", config)
	if err != nil {
		panic(err)
	}

}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

package client

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io/ioutil"
	"testing"
	"time"
)

type ClientTestSuite struct {
	suite.Suite
	client *AmbariClient
}

func (s *ClientTestSuite) SetupSuite() {

	// Init logger
	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

	// Init client
	s.client = New("http://ambari-server:8080/api/v1", "admin", "admin")
	s.client.DisableVerifySSL()

	// Create repository
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
	_, err := s.client.CreateRepository(repository)
	if err != nil {
		panic(err)
	}

	// Create blueprint
	b, err := ioutil.ReadFile("../fixtures/blueprint.json")
	if err != nil {
		panic(err)
	}
	blueprintJson := string(b)
	_, err = s.client.CreateBlueprint("test", blueprintJson)
	if err != nil {
		panic(err)
	}

	// Wait agent join the cluster
	isAgentJoinCluster := false
	for isAgentJoinCluster == false {
		host, err := s.client.Host("ambari-agent")
		if err != nil {
			panic(err)
		}
		if host != nil {
			isAgentJoinCluster = true
		} else {
			time.Sleep(5)
		}
	}
	isAgentJoinCluster = false
	for isAgentJoinCluster == false {
		host, err := s.client.Host("ambari-agent")
		if err != nil {
			panic(err)
		}
		if host != nil {
			isAgentJoinCluster = true
		} else {
			time.Sleep(5)
		}
	}
	isAgentJoinCluster = false
	for isAgentJoinCluster == false {
		host, err := s.client.Host("ambari-agent3")
		if err != nil {
			panic(err)
		}
		if host != nil {
			isAgentJoinCluster = true
		} else {
			time.Sleep(5)
		}
	}

	// Create  a real fresh cluster with blueprint
	b, err = ioutil.ReadFile("../fixtures/cluster-template.json")
	if err != nil {
		panic(err)
	}
	templateJson := string(b)
	_, err = s.client.CreateClusterFromTemplate("test", templateJson)
	if err != nil {
		panic(err)
	}

}

func (s *ClientTestSuite) SetupTest() {

	// Wait all task before run the next test
	s.WaitClusterTask()

}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (s *ClientTestSuite) WaitClusterTask() {

	path := "/clusters/test/requests"

	isRequestTask := true
	for isRequestTask {
		resp, err := s.client.Client().R().SetQueryString("Requests/progress_percent!=100").Get(path)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode() >= 300 {
			panic(fmt.Sprintf("Problems when request current task: %d - %s", resp.StatusCode(), resp.String()))
		}
		requestsTask := &RequestsTask{}
		err = json.Unmarshal(resp.Body(), requestsTask)
		if err != nil {
			panic(err)
		}

		if len(requestsTask.Items) == 0 {
			isRequestTask = false
		} else {
			time.Sleep(5 * time.Second)
		}

	}

}

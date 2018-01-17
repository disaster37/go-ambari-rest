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

	s.client = New("https://10.221.78.60:5010/api/v1", "admin", "admin")
	s.client.DisableVerifySSL()

}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

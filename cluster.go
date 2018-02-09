package main

import (
	"github.com/disaster37/go-ambari-rest/client"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
)

func createCluster(c *cli.Context) error {

	clientAmbari, err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if c.String("cluster-name") == "" {
		return cli.NewExitError("You must set cluster-name parameter", 1)
	}
	if c.String("blueprint-file") == "" {
		return cli.NewExitError("You must set blueprint-file parameter", 1)
	}
	if c.String("hosts-template-file") == "" {
		return cli.NewExitError("You must set hosts-template-file parameter", 1)
	}

	// Read the Json files
	b, err := ioutil.ReadFile(c.String("blueprint-file"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	blueprintJson := string(b)
	log.Debug("BlueprintJson: ", blueprintJson)
	b, err = ioutil.ReadFile(c.String("hosts-template-file"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	hostsTemplateJson := string(b)
	log.Debug("HostsTemplateJson: ", hostsTemplateJson)

	// Create the blueprint
	_, err = clientAmbari.CreateBlueprint(c.String("cluster-name"), blueprintJson)
	if err != nil {
		ambariError := err.(client.AmbariError)
		if ambariError.Code != 409 {
			return cli.NewExitError(ambariError.Message, 1)
		}
	}

	// Create the cluster
	_, err = clientAmbari.CreateClusterFromTemplate(c.String("cluster-name"), hostsTemplateJson)
	if err != nil {
		ambariError := err.(client.AmbariError)
		if ambariError.Code != 409 {
			return cli.NewExitError(ambariError.Message, 1)
		}
	}

	log.Info("Repository created successfully")
	return nil

}

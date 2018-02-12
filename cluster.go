package main

import (
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

	// Check if blueprint already exist
	blueprint, err := clientAmbari.Blueprint(c.String("cluster-name"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if blueprint == nil {
		// Create the blueprint
		_, err = clientAmbari.CreateBlueprint(c.String("cluster-name"), blueprintJson)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		log.Info("Create blueprint successfully")
	} else {
		log.Info("Blueprint already exist, skip.")
	}

	// Check if cluster already exist
	cluster, err := clientAmbari.Cluster(c.String("cluster-name"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if cluster == nil {
		// Create the cluster
		_, err = clientAmbari.CreateClusterFromTemplate(c.String("cluster-name"), hostsTemplateJson)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		log.Info("Cluster created successfully, look /var/log/ambari-server/ambari-server.log about potential topologie error")
	} else {
		log.Info("Cluster already exist, skip")
	}

	return nil

}

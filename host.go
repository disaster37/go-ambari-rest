package main

import (
	"github.com/disaster37/go-ambari-rest/client"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

func addHostInCluster(c *cli.Context) error {

	clientAmbari, err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if c.String("cluster-name") == "" {
		return cli.NewExitError("You must set cluster-name parameter", 1)
	}
	if c.String("blueprint-name") == "" {
		return cli.NewExitError("You must set blueprint-name parameter", 1)
	}
	if c.String("hostname") == "" {
		return cli.NewExitError("You must set hostname parameter", 1)
	}
	if c.String("role") == "" {
		return cli.NewExitError("You must set role parameter", 1)
	}

	// Register host in cluster
	err = clientAmbari.RegisterHostOnCluster(c.String("cluster-name"), c.String("hostname"), c.String("blueprint-name"), c.String("role"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	log.Infof("Successfully add new host %s in cluster %s with role %", c.String("hostname"), c.String("cluster-name"), c.String("role"))

	// Set the rack for host
	if c.String("rack") != "" {
		host, err := clientAmbari.HostOnCluster(c.String("cluster-name"), c.String("hostname"))
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if host == nil {
			return cli.NewExitError(client.NewAmbariError(404, "Host %s not found when try to set the rack", c.String("hostname")), 1)
		}

		host.HostInfo.Rack = c.String("rack")
		host, err = clientAmbari.UpdateHost(host)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		log.Infof("Successfully set rack %s to host %s", c.String("rack"), c.String("hostname"))

	}

	return nil
}

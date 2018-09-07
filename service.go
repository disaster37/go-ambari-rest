package main

import (
	"github.com/disaster37/go-ambari-rest/client"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

func stopServiceInCluster(c *cli.Context) error {

	clientAmbari, err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if c.String("cluster-name") == "" {
		return cli.NewExitError("You must set cluster-name parameter", 1)
	}
	if c.String("service-name") == "" {
		return cli.NewExitError("You must set service-name parameter", 1)
	}

	// Stop the service
	_, err = clientAmbari.StopService(c.String("cluster-name"), c.String("service-name"), c.Bool("enable-maintenance"), c.Bool("force"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	log.Infof("Successfully stop service %s in cluster %s with enable maintenance mode to %t", c.String("service-name"), c.String("cluster-name"), c.Bool("enable-maintenance"))

	return nil
}

func startServiceInCluster(c *cli.Context) error {

	clientAmbari, err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if c.String("cluster-name") == "" {
		return cli.NewExitError("You must set cluster-name parameter", 1)
	}
	if c.String("service-name") == "" {
		return cli.NewExitError("You must set service-name parameter", 1)
	}

	// Stop the service
	_, err = clientAmbari.StartService(c.String("cluster-name"), c.String("service-name"), c.Bool("disable-maintenance"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	log.Infof("Successfully start service %s in cluster %s with  disable maintenance mode to %t", c.String("service-name"), c.String("cluster-name"), c.Bool("disable-maintenance"))

	return nil
}

func stopAllServicesInCluster(c *cli.Context) error {

	clientAmbari, err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if c.String("cluster-name") == "" {
		return cli.NewExitError("You must set cluster-name parameter", 1)
	}

	// Get the cluster
	cluster, err := clientAmbari.Cluster(c.String("cluster-name"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if cluster == nil {
		return cli.NewExitError(client.NewAmbariError(404, "Cluster %s not found", c.String("cluster-name")), 1)
	}

	// Stop all the services
	err = clientAmbari.StopAllServices(cluster, c.Bool("enable-maintenance"), c.Bool("force"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	log.Infof("Successfully stop all services in cluster %s with enable maintenance mode to %t", c.String("cluster-name"), c.Bool("enable-maintenance"))

	return nil
}

func startAllServicesInCluster(c *cli.Context) error {

	clientAmbari, err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if c.String("cluster-name") == "" {
		return cli.NewExitError("You must set cluster-name parameter", 1)
	}

	// Get the cluster
	cluster, err := clientAmbari.Cluster(c.String("cluster-name"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if cluster == nil {
		return cli.NewExitError(client.NewAmbariError(404, "Cluster %s not found", c.String("cluster-name")), 1)
	}

	// Start all the services
	err = clientAmbari.StartAllServices(cluster, c.Bool("disable-maintenance"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	log.Infof("Successfully start all services in cluster %s", c.String("cluster-name"))

	return nil
}

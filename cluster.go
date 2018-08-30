package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"time"
)

type ClusterTemplate struct {
	HostGroups []struct {
		Name  string `json:"name"`
		Hosts []Host `json:"hosts"`
	} `json:"host_groups"`
}

type Host struct {
	FQDN string `json:"fqdn"`
}

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

		// Before create cluster, we need wait all node join to macro substitution work fine
		clusterTemplate := &ClusterTemplate{}
		err = json.Unmarshal([]byte(hostsTemplateJson), clusterTemplate)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		nbNodes := 0
		for _, hostGroup := range clusterTemplate.HostGroups {
			nbNodes = nbNodes + len(hostGroup.Hosts)
		}

		log.Infof("Wait all nodes (%d) join Ambari server to avoid hostgroup substitution ...", nbNodes)
		loop := true
		for loop == true {
			loop = false
			for idx, hostGroup := range clusterTemplate.HostGroups {
				tempHosts := make([]Host, 0, len(hostGroup.Hosts))
				for _, hostTemp := range hostGroup.Hosts {
					// Check if host already here
					host, err := clientAmbari.Host(hostTemp.FQDN)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					if host == nil {
						// Wait host join
						loop = true
						tempHosts = append(tempHosts, hostTemp)
						log.Infof("Host %s not yet join the cluster, continuous to wait...", hostTemp.FQDN)
						time.Sleep(10 * time.Second)
					} else {
						log.Infof("Host %s already join the cluster", hostTemp.FQDN)
					}
				}
				clusterTemplate.HostGroups[idx].Hosts = tempHosts
			}
		}
		log.Info("All nodes have join the Ambari server.")

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

package main

import (
	"github.com/disaster37/go-ambari-rest/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/altsrc"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/urfave/cli.v1"
	"os"
)

var debug bool
var ambariURL string
var ambariLogin string
var ambariPassword string

func main() {

	// Logger setting
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)

	// CLI settings
	app := cli.NewApp()
	app.Usage = "Manage Ambari on cli interface"
	app.Version = "develop"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "Load configuration from `FILE`",
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "ambari-url",
			Usage:       "The Ambari base URL (with api version)",
			EnvVar:      "AMBARI_URL",
			Destination: &ambariURL,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "ambari-login",
			Usage:       "The Ambari admin login",
			EnvVar:      "AMBARI_LOGIN",
			Destination: &ambariLogin,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "ambari-password",
			Usage:       "The Ambari admin password",
			EnvVar:      "AMBARI_PASSWORD",
			Destination: &ambariPassword,
		}),
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Display debug output",
			Destination: &debug,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "create-or-update-repository",
			Usage: "Create or update repository",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "repository-file",
					Usage: "The full path of repository file",
				},
				cli.BoolFlag{
					Name:  "use-spacewalk",
					Usage: "Permit to use spacewalk. Default is false.",
				},
			},
			Action: createRepository,
		},
		{
			Name:  "create-cluster-if-not-exist",
			Usage: "Create new cluster if not exist with blueprint and hosts template",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name you should to create",
				},
				cli.StringFlag{
					Name:  "blueprint-file",
					Usage: "The full path of blueprint file",
				},
				cli.StringFlag{
					Name:  "hosts-template-file",
					Usage: "The full path of hosts template file",
				},
			},
			Action: createCluster,
		},
		{
			Name:  "create-or-update-privileges",
			Usage: "Create or update privileges",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "privileges-file",
					Usage: "The full path of privileges file",
				},
			},
			Action: createPrivileges,
		},
		{
			Name:  "add-host-in-cluster",
			Usage: "Add new host in existing cluster",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name to add host",
				},
				cli.StringFlag{
					Name:  "blueprint-name",
					Usage: "The blueprint name to use to affect host role in cluster",
				},
				cli.StringFlag{
					Name:  "hostname",
					Usage: "The hostname to add on cluster",
				},
				cli.StringFlag{
					Name:  "role",
					Usage: "The role of the host in cluster",
				},
				cli.StringFlag{
					Name:  "rack",
					Usage: "The rack name",
				},
			},
			Action: addHostInCluster,
		},
		{
			Name:  "stop-service",
			Usage: "Stop service and wait service is stopped",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name where to stop the service",
				},
				cli.StringFlag{
					Name:  "service-name",
					Usage: "The service name to stop",
				},
				cli.BoolFlag{
					Name:  "enable-maintenance",
					Usage: "Put service in maintenance mode before stop the service",
				},
				cli.BoolFlag{
					Name:  "force",
					Usage: "Remove maintenance state in service before stop them",
				},
			},
			Action: stopServiceInCluster,
		},
		{
			Name:  "start-service",
			Usage: "Start service and wait service is started",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name where to stop the service",
				},
				cli.StringFlag{
					Name:  "service-name",
					Usage: "The service name to stop",
				},
				cli.BoolFlag{
					Name:  "disable-maintenance",
					Usage: "Put service in maintenance mode OFF after start the service",
				},
			},
			Action: startServiceInCluster,
		},
		{
			Name:  "stop-all-services",
			Usage: "Stop all services and wait all services are stopped",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name where to stop all services",
				},
				cli.BoolFlag{
					Name:  "enable-maintenance",
					Usage: "Put all services in maintenance state after stop all services",
				},
				cli.BoolFlag{
					Name:  "force",
					Usage: "Remove maintenance state in all services before stop them",
				},
			},
			Action: stopAllServicesInCluster,
		},
		{
			Name:  "start-all-services",
			Usage: "Start all services and wait all services are started",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name where to stop all services",
				},
				cli.BoolFlag{
					Name:  "disable-maintenance",
					Usage: "Remove all mainetnance state in all services before start them",
				},
			},
			Action: startAllServicesInCluster,
		},
		{
			Name:  "stop-all-components-in-host",
			Usage: "Stop all components in host and wait all components are stopped",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name where to stop all components",
				},
				cli.StringFlag{
					Name:  "hostname",
					Usage: "The hostname where to stop all components",
				},
				cli.BoolFlag{
					Name:  "enable-maintenance",
					Usage: "Enable maintenance state in host after stop components",
				},
				cli.BoolFlag{
					Name:  "force",
					Usage: "Disable maintenance state in host before stop components",
				},
			},
			Action: stopAllComponentsInHost,
		},
		{
			Name:  "start-all-components-in-host",
			Usage: "Start all components in host and wait all components are started",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name where to start all components",
				},
				cli.StringFlag{
					Name:  "hostname",
					Usage: "The hostname where to start all components",
				},
				cli.BoolFlag{
					Name:  "disable-maintenance",
					Usage: "Disable maintenance state in host before start components",
				},
			},
			Action: startAllComponentsInHost,
		},
		{
			Name:  "start-component-in-host",
			Usage: "Start component in host and wait component are started",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name where to start component",
				},
				cli.StringFlag{
					Name:  "hostname",
					Usage: "The hostname where to start component",
				},
				cli.StringFlag{
					Name:  "component-name",
					Usage: "The component name to start",
				},
			},
			Action: startComponentInHost,
		},
		{
			Name:  "stop-component-in-host",
			Usage: "Stop component in host and wait component are stopped",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name where to stop component",
				},
				cli.StringFlag{
					Name:  "hostname",
					Usage: "The hostname where to stop component",
				},
				cli.StringFlag{
					Name:  "component-name",
					Usage: "The component name to stop",
				},
			},
			Action: stopComponentInHost,
		},
	}

	app.Before = func(c *cli.Context) error {
		if c.String("config") != "" {
			before := altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewYamlSourceFromFlagFunc("config"))
			return before(c)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Check the global parameter
func manageGlobalParameters() (*client.AmbariClient, error) {
	if debug == true {
		log.SetLevel(log.DebugLevel)
	}

	if ambariURL == "" {
		return nil, errors.New("You must set --ambari-url parameter")
	}

	if ambariLogin == "" {
		return nil, errors.New("You must set --ambari-login parameter")
	}
	if ambariPassword == "" {
		return nil, errors.New("You must set --ambari-password parameter")
	}

	client := client.New(ambariURL, ambariLogin, ambariPassword)
	client.DisableVerifySSL()

	return client, nil
}

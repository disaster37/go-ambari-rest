package main

import (
	"github.com/disaster37/go-ambari-rest/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
			Name:        "ambari-url",
			Usage:       "The Ambari base URL (with api version)",
			Destination: &ambariURL,
		},
		cli.StringFlag{
			Name:        "ambari-login",
			Usage:       "The Ambari admin login",
			Destination: &ambariLogin,
		},
		cli.StringFlag{
			Name:        "ambari-password",
			Usage:       "The Ambari admin password",
			Destination: &ambariPassword,
		},
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
	}

	app.Run(os.Args)
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

	return client, nil
}

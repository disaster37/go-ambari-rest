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
		{
			Name:  "configure-kerberos",
			Usage: "Install and configure Kerberos on HDP cluster",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name where to enable kerberos",
				},
				cli.StringFlag{
					Name:  "kdc-type",
					Usage: "The kdc type to use (active-directory, mit-kdc or ipa)",
					Value: "active-directory",
				},
				cli.BoolFlag{
					Name:  "disable-manage-identities",
					Usage: "Manage Kerberos principals and keytabs manually",
				},

				cli.StringFlag{
					Name:  "kdc-hosts",
					Usage: "A comma separated list of KDC host. Optionnaly a port number may be included",
				},
				cli.StringFlag{
					Name:  "realm",
					Usage: "The default realm to use when creating service principal",
				},
				cli.StringFlag{
					Name:  "ldap-url",
					Usage: "The URL to Active Directory LDAP server. Only needed if the KDC type is Active Directory",
				},
				cli.StringFlag{
					Name:  "container-dn",
					Usage: "The DN of the container used store service principals. Only needed if you use Active Directory",
				},
				cli.StringFlag{
					Name:  "domains",
					Usage: "A comma separated list of domain names used to map server host names to the REALM name. It's optionnal",
				},
				cli.StringFlag{
					Name:  "admin-server-host",
					Usage: "The host for KDC Kerberos administrative host. Optionnaly the port number can be included",
				},
				cli.StringFlag{
					Name:  "principal-name",
					Usage: "Admin principal used to create principals and export keytabs",
				},
				cli.StringFlag{
					Name:  "principal-password",
					Usage: "Admin principal password",
				},
				cli.BoolFlag{
					Name:  "persist-credential",
					Usage: "Store admin credential. Need to enable password encryption before that",
				},
				cli.BoolFlag{
					Name:  "disable-install-packages",
					Usage: "Disable the installation of Kerberos client package",
				},
				cli.StringFlag{
					Name:  "executable-search-paths",
					Usage: "A comma delimited list of search paths used to find Kerberos utilities",
					Value: "/usr/bin, /usr/kerberos/bin, /usr/sbin, /usr/lib/mit/bin, /usr/lib/mit/sbin",
				},
				cli.StringFlag{
					Name:  "encryption-type",
					Usage: "The supported list of session key encryption types that should be returned by the KDC",
					Value: "aes des3-cbc-sha1 rc4 des-cbc-md5",
				},
				cli.Int64Flag{
					Name:  "password-length",
					Usage: "The password length",
					Value: 20,
				},
				cli.Int64Flag{
					Name:  "password-min-lowercase-letters",
					Usage: "The minimal lowercase letters to compose password",
					Value: 1,
				},
				cli.Int64Flag{
					Name:  "password-min-uppercase-letters",
					Usage: "The minimal uppercase letters to compose password",
					Value: 1,
				},
				cli.Int64Flag{
					Name:  "password-min-digits",
					Usage: "The minimal digits to compose password",
					Value: 1,
				},
				cli.Int64Flag{
					Name:  "password-min-punctuation",
					Usage: "The minimal punctuation to compose password",
					Value: 1,
				},
				cli.Int64Flag{
					Name:  "password-min-whitespace",
					Usage: "The minimal whitespace to compose password",
					Value: 0,
				},
				cli.StringFlag{
					Name:  "check-principal-name",
					Usage: "The principal name to use when executing Kerberos service check",
					Value: "${cluster_name|toLower()}-${short_date}",
				},
				cli.BoolFlag{
					Name:  "enable-case-insensitive-username-rules",
					Usage: "Force principal names to resolv to lowercase local usernames in auth-to-local rules",
				},
				cli.BoolFlag{
					Name:  "disable-manage-auth-to-local",
					Usage: "Don't manage the Hadoop auth-to-local rules by Ambari",
				},
				cli.BoolFlag{
					Name:  "disable-create-ambari-principal",
					Usage: "Don't create principal and keytab by Ambari",
				},
				cli.StringFlag{
					Name:  "master-kdc-host",
					Usage: "The master KDC host in master/slave KDC deployment",
				},
				cli.StringFlag{
					Name:  "preconfigure-services",
					Usage: "Preconfigure service. Possible value are NONE, DEFAULT or ALL.",
					Value: "DEFAULT",
				},
				cli.StringFlag{
					Name:  "ad-create-attributes-template",
					Usage: "A velocity template to use when create service principals in Active Directory.",
					Value: "\n{\n  \"objectClass\": [\"top\", \"person\", \"organizationalPerson\", \"user\"],\n  \"cn\": \"$principal_name\",\n  #if( $is_service )\n  \"servicePrincipalName\": \"$principal_name\",\n  #end\n  \"userPrincipalName\": \"$normalized_principal\",\n  \"unicodePwd\": \"$password\",\n  \"accountExpires\": \"0\",\n  \"userAccountControl\": \"66048\"\n}",
				},
				cli.BoolFlag{
					Name:  "disable-manage-krb5-conf",
					Usage: "Don't manage krb5.conf by Ambari",
				},
				cli.StringFlag{
					Name:  "krb5-conf-directory",
					Usage: "The krb5.conf coonfiguration directory",
					Value: "/etc",
				},
				cli.StringFlag{
					Name:  "krb5-conf-template",
					Usage: "The krb5.conf template",
					Value: "[libdefaults]\n  renew_lifetime = 7d\n  forwardable= true\n  default_realm = {{realm|upper()}}\n  ticket_lifetime = 24h\n  dns_lookup_realm = false\n  dns_lookup_kdc = false\n  #default_tgs_enctypes = {{encryption_types}}\n  #default_tkt_enctypes ={{encryption_types}}\n\n{% if domains %}\n[domain_realm]\n{% for domain in domains.split(',') %}\n  {{domain}} = {{realm|upper()}}\n{% endfor %}\n{%endif %}\n\n[logging]\n  default = FILE:/var/log/krb5kdc.log\nadmin_server = FILE:/var/log/kadmind.log\n  kdc = FILE:/var/log/krb5kdc.log\n\n[realms]\n  {{realm}} = {\n    admin_server = {{admin_server_host|default(kdc_host, True)}}\n    kdc = {{kdc_host}}\n }\n\n{# Append additional realm declarations below #}\n",
				},
			},
			Action: addKerberos,
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

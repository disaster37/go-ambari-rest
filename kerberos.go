package main

import (
	"fmt"
	"github.com/disaster37/go-ambari-rest/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"strconv"
	"time"
)

const (
	KERBEROS_SERVICE     = "KERBEROS"
	KERBEROS_COMPONENT   = "KERBEROS_CLIENT"
	ALIAS_KDC_CREDENTIAL = "kdc.admin.credential"
)

func addKerberos(c *cli.Context) error {

	clientAmbari, err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if c.String("cluster-name") == "" {
		return cli.NewExitError("You must set cluster-name parameter", 1)
	}
	if c.String("kdc-type") == "" {
		return cli.NewExitError("You must set kdc-type parameter", 1)
	}
	if c.String("realm") == "" {
		return cli.NewExitError("You must set realm", 1)
	}
	if c.String("executable-search-paths") == "" {
		return cli.NewExitError("executable-search-paths", 1)
	}
	if c.String("kdc-type") == "active-directory" {
		if c.String("ldap-url") == "" {
			return cli.NewExitError("You must set ldap-url parameter", 1)
		}
		if c.String("container-dn") == "" {
			return cli.NewExitError("You must set container-dn parameter", 1)
		}
		if c.String("ad-create-attributes-template") == "" {
			return cli.NewExitError("You must set ad-create-attributes-template parameter", 1)
		}
	}
	if c.Bool("disable-manage-identities") != true {
		if c.String("principal-name") == "" {
			return cli.NewExitError("You must set principal-name parameter", 1)
		}
		if c.String("principal-password") == "" {
			return cli.NewExitError("You must set principal-password parameter", 1)
		}
		if c.String("encryption-type") == "" {
			return cli.NewExitError("You must set encryption-type parameter", 1)
		}
		if c.String("kdc-hosts") == "" {
			return cli.NewExitError("You must set kdc-hosts parameter", 1)
		}
		if c.String("check-principal-name") == "" {
			return cli.NewExitError("check-principal-name", 1)
		}
		if c.String("preconfigure-services") == "" {
			return cli.NewExitError("preconfigure-services", 1)
		}
		if c.String("admin-server-host") == "" {
			return cli.NewExitError("You must set admin-server-host parameter", 1)
		}
		if c.String("krb5-conf-directory") == "" {
			return cli.NewExitError("You must set krb5-conf-directory parameter", 1)
		}
		if c.String("krb5-conf-template") == "" {
			return cli.NewExitError("You must set krb5-conf-template parameter", 1)
		}
	}
	manageKrb5Conf := "true"
	if c.Bool("disable-manage-krb5-conf") {
		manageKrb5Conf = "false"
	}
	manageIdentities := "true"
	if c.Bool("disable-manage-identities") {
		manageIdentities = "false"
	}
	installPackages := "true"
	if c.Bool("disable-install-packages") {
		installPackages = "false"
	}
	createAmbariPrincipal := "true"
	if c.Bool("disable-create-ambari-principal") {
		createAmbariPrincipal = "false"
	}
	caseInsensitiveUsernameRules := "false"
	if c.Bool("enable-case-insensitive-username-rules") {
		caseInsensitiveUsernameRules = "true"
	}
	manageAuthToLocal := "true"
	if c.Bool("disable-manage-auth-to-local") {
		manageAuthToLocal = "false"
	}
	persistCredential := client.CREDENTIAL_TEMPORARY
	if c.Bool("persist-credential") {
		persistCredential = client.CREDENTIAL_PERSISTED
	}

	// Get cluster object
	cluster, err := clientAmbari.Cluster(c.String("cluster-name"))
	if err != nil {
		return err
	}
	if cluster == nil {
		return errors.New(fmt.Sprintf("Cluster %s not found", c.String("cluster-name")))
	}

	// Add Kerberos service if needed
	serviceKerberos, err := clientAmbari.Service(c.String("cluster-name"), KERBEROS_SERVICE)
	if err != nil {
		return err
	}
	if serviceKerberos == nil {
		// Add kerberos service
		serviceKerberos = &client.Service{
			ServiceInfo: &client.ServiceInfo{
				ClusterName: c.String("cluster-name"),
				ServiceName: KERBEROS_SERVICE,
			},
		}
		serviceKerberos, err = clientAmbari.CreateService(serviceKerberos)
		if err != nil {
			return err
		}
		log.Infof("%s service is created", KERBEROS_SERVICE)
	} else {
		log.Infof("%s service is already exist", KERBEROS_SERVICE)
	}

	// Add KERBEROS_CLIENT components if needed
	componentKerberosClient, err := clientAmbari.Component(c.String("cluster-name"), KERBEROS_SERVICE, KERBEROS_COMPONENT)
	if err != nil {
		return err
	}
	if componentKerberosClient == nil {
		componentKerberosClient = &client.Component{
			ComponentInfo: &client.ComponentInfo{
				ClusterName:   c.String("cluster-name"),
				ServiceName:   KERBEROS_SERVICE,
				ComponentName: KERBEROS_COMPONENT,
			},
		}
		componentKerberosClient, err = clientAmbari.CreateComponent(componentKerberosClient)
		if err != nil {
			return err
		}
		log.Infof("%s component is created", KERBEROS_COMPONENT)
	} else {
		log.Infof("%s component is already exist", KERBEROS_COMPONENT)
	}

	// Add Kerberos service settings
	t := time.Now()
	tag := fmt.Sprintf("version_%s", t.Format("2006-01-02_15:04:05"))
	configurationKerberosService := &client.Configuration{
		Type: "krb5-conf",
		Tag:  tag,
		Properties: map[string]string{
			"domains":          c.String("domains"),
			"manage_krb5_conf": manageKrb5Conf,
			"conf_dir":         c.String("krb5-conf-directory"),
			"content":          c.String("krb5-conf-template"),
		},
	}
	_, err = clientAmbari.CreateConfigurationOnCluster(c.String("cluster-name"), configurationKerberosService)
	if err != nil {
		return err
	}
	log.Info("Setting 'krb-conf' is created in cluster")
	configurationKerberosEnv := &client.Configuration{
		Type: "kerberos-env",
		Tag:  tag,
		Properties: map[string]string{
			"kdc_type":                        c.String("kdc-type"),
			"manage_identities":               manageIdentities,
			"install_packages":                installPackages,
			"encryption_types":                c.String("encryption-type"),
			"realm":                           c.String("realm"),
			"kdc_hosts":                       c.String("kdc-hosts"),
			"master_kdc":                      c.String("master-kdc-host"),
			"admin_server_host":               c.String("admin-server-host"),
			"executable_search_paths":         c.String("executable-search-paths"),
			"password_length":                 strconv.FormatInt(c.Int64("password-length"), 10),
			"password_min_lowercase_letters":  strconv.FormatInt(c.Int64("password-min-lowercase-letters"), 10),
			"password_min_uppercase_letters":  strconv.FormatInt(c.Int64("password-min-uppercase-letters"), 10),
			"password_min_digits":             strconv.FormatInt(c.Int64("password-min-digits"), 10),
			"password_min_punctuation":        strconv.FormatInt(c.Int64("password-min-punctuation"), 10),
			"password_min_whitespace":         strconv.FormatInt(c.Int64("password-min-whitespace"), 10),
			"service_check_principal_name":    c.String("check-principal-name"),
			"case_insensitive_username_rules": caseInsensitiveUsernameRules,
			"create_ambari_principal":         createAmbariPrincipal,
			"container_dn":                    c.String("container-dn"),
			"ad_create_attributes_template":   c.String("ad-create-attributes-template"),
			"ldap_url":                        c.String("ldap-url"),
			"manage_auth_to_local":            manageAuthToLocal,
			"kdc_create_attributes":           c.String("kdc-create-attributes"),
			"preconfigure_services":           c.String("preconfigure-services"),
		},
	}
	_, err = clientAmbari.CreateConfigurationOnCluster(c.String("cluster-name"), configurationKerberosEnv)
	if err != nil {
		return err
	}
	log.Info("Setting 'kerberos-env' is created in cluster")

	// Add KERBEROS_CLIENT components on all nodes if needed
	hosts, err := clientAmbari.HostsOnCluster(c.String("cluster-name"))
	if err != nil {
		return err
	}
	log.Debugf("Found %d hosts in cluster", len(hosts))
	for _, host := range hosts {
		hostComponent, err := clientAmbari.HostComponent(c.String("cluster-name"), host.HostInfo.Hostname, KERBEROS_COMPONENT)
		if err != nil {
			return err
		}
		if hostComponent == nil {
			hostComponent := &client.HostComponent{
				HostComponentInfo: &client.HostComponentInfo{
					ClusterName:   c.String("cluster-name"),
					ServiceName:   KERBEROS_SERVICE,
					ComponentName: KERBEROS_COMPONENT,
					Hostname:      host.HostInfo.Hostname,
				},
			}
			hostComponent, err := clientAmbari.CreateHostComponent(hostComponent)
			if err != nil {
				return err
			}
			log.Infof("Component %s is associated to host %s", KERBEROS_COMPONENT, host.HostInfo.Hostname)
		} else {
			log.Infof("Component %s is already associated to host %s", KERBEROS_COMPONENT, host.HostInfo.Hostname)
		}

	}

	// Install service/component Kerberos on all nodes
	serviceKerberos, err = clientAmbari.InstallService(serviceKerberos)
	if err != nil {
		return err
	}
	log.Info("Service KERBEROS is installed")

	// Stop all services
	err = clientAmbari.StopAllServices(cluster, false, true)
	if err != nil {
		return nil
	}
	log.Info("All services are stopped")

	// Create or update Kerberos credential
	credential, err := clientAmbari.Credential(c.String("cluster-name"), ALIAS_KDC_CREDENTIAL)
	if err != nil {
		return err
	}
	if credential == nil {
		credential = &client.Credential{
			CredentialInfo: &client.CredentialInfo{
				Principal:   c.String("principal-name"),
				Key:         c.String("principal-password"),
				Type:        persistCredential,
				Alias:       ALIAS_KDC_CREDENTIAL,
				ClusterName: c.String("cluster-name"),
			},
		}
		_, err = clientAmbari.CreateCredential(credential)
		if err != nil {
			return err
		}
		log.Infof("Create credential %s on Ambari", ALIAS_KDC_CREDENTIAL)
	} else {
		credential = &client.Credential{
			CredentialInfo: &client.CredentialInfo{
				Alias:       ALIAS_KDC_CREDENTIAL,
				Principal:   c.String("principal-name"),
				Key:         c.String("principal-password"),
				Type:        persistCredential,
				ClusterName: c.String("cluster-name"),
			},
		}
		_, err = clientAmbari.UpdateCredential(credential)
		if err != nil {
			return err
		}
		log.Infof("Update credential %s on Ambari", ALIAS_KDC_CREDENTIAL)
	}

	// Enabling Kerberos
	cluster.ClusterInfo.SecurityType = "KERBEROS"

	cluster.SessionAttributes = map[string]map[string]string{
		"kerberos_admin": map[string]string{
			"principal": c.String("principal-name"),
			"password":  c.String("principal-password"),
		},
	}
	cluster, err = clientAmbari.UpdateCluster(cluster)
	if err != nil {
		return err
	}
	log.Info("Kerberos is enabled")

	/*

		// Start all services
		err = clientAmbari.StartAllServices(cluster, false)
		if err != nil {
			return err
		}
		log.Info("All services are started")
	*/

	return nil

}

package main

import (
	"encoding/json"
	"github.com/disaster37/go-ambari-rest/client"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
)

type Privileges struct {
	ClusterName string      `json:"clusterName"`
	Privileges  []Privilege `json:"privileges"`
}

type Privilege struct {
	Permission string `json:"permission"`
	Type       string `json:"type"`
	Name       string `json:"name"`
}

func createPrivileges(c *cli.Context) error {

	clientAmbari, err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if c.String("privileges-file") == "" {
		return cli.NewExitError("You must set --privileges-file parameter", 1)
	}

	// Read the Json file
	b, err := ioutil.ReadFile(c.String("privileges-file"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	privilegesJson := string(b)
	log.Debug("Privileges: ", privilegesJson)
	privileges := &Privileges{}
	err = json.Unmarshal(b, privileges)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	//Loop over privileges
	for _, privilegeItem := range privileges.Privileges {
		// Check if privilege already exist
		privilege, err := clientAmbari.SearchPrivilege(privileges.ClusterName, privilegeItem.Permission, privilegeItem.Name, privilegeItem.Type)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		privilegeTarget := &client.Privilege{
			PrivilegeInfo: &client.PrivilegeInfo{
				PermissionName: privilegeItem.Permission,
				PrincipalName:  privilegeItem.Name,
				PrincipalType:  privilegeItem.Type,
			},
		}

		if privilege == nil {
			// Create new privilege
			_, err = clientAmbari.CreatePrivilege(privileges.ClusterName, privilegeTarget)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			log.Infof("Create privilege %s / %s successfully", privilegeItem.Name, privilegeItem.Permission)
		} else {
			//Update privilege
			privilegeTarget.PrivilegeInfo.PrivilegeId = privilege.PrivilegeInfo.PrivilegeId
			_, err = clientAmbari.UpdatePrivilege(privileges.ClusterName, privilegeTarget)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			log.Infof("Update privilege %s / %s successfully", privilegeItem.Name, privilegeItem.Permission)
		}
	}

	return nil

}

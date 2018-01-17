package client

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Privilege struct {
	Response
	PrivilegeInfo *PrivilegeInfo `json:"PrivilegeInfo"`
}

type PrivilegesResponse struct {
	Response
	Items []Privilege `json:"items"`
}

type PrivilegeInfo struct {
	PrivilegeId     int64  `json:"privilege_id,omitempty"`
	PermissionLabel string `json:"permission_label,omitempty"`
	PermissionName  string `json:"permission_name"`
	PrincipalName   string `json:"principal_name"`
	PrincipalType   string `json:"principal_type"`
}

func (c *AmbariClient) Privilege(clusterName string, id int64) (*Privilege, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}

	path := fmt.Sprintf("/clusters/%s/privileges/%d", clusterName, id)

	resp, err := c.Client().R().Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Result : ", resp)

	privilege := &Privilege{}
	err = json.Unmarshal(resp.Body(), privilege)
	if err != nil {
		return nil, err
	}

	log.Debug("Privilege: ", privilege)

	return privilege, nil

}

func (c *AmbariClient) CreatePrivilege(clusterName string, privilege *Privilege) (*Privilege, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if privilege == nil {
		panic("Privilege can't be nil")
	}

	// Create the privilege
	path := fmt.Sprintf("/clusters/%s/privileges", clusterName)
	jsonData, err := json.Marshal(privilege)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, errors.New(resp.Status())
	}

	// Get the privilege
	resp, err = c.Client().R().SetQueryParams(map[string]string{
		"PrivilegeInfo/permission_name": privilege.PrivilegeInfo.PermissionName,
		"PrivilegeInfo/principal_name":  privilege.PrivilegeInfo.PrincipalName,
		"PrivilegeInfo/principal_type":  privilege.PrivilegeInfo.PrincipalType,
	}).Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to get: ", resp)
	privilegeResponses := &PrivilegesResponse{}
	err = json.Unmarshal(resp.Body(), privilegeResponses)
	if err != nil {
		return nil, err
	}
	log.Debug("PrivilegesResponse: ", privilegeResponses)

	if len(privilegeResponses.Items) > 0 {
		log.Debug("Privilege: ", privilegeResponses.Items[0])
		return &privilegeResponses.Items[0], nil
	} else {
		return nil, nil
	}

}

func (c *AmbariClient) DeletePrivilege(clusterName string, id int64) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}

	// Create the privilege
	path := fmt.Sprintf("/clusters/%s/privileges/%d", clusterName, id)
	resp, err := c.Client().R().Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response to delete privilege: ", resp)
	if resp.StatusCode() >= 300 {
		return errors.New(resp.Status())
	}

	return nil

}

func (c *AmbariClient) UpdatePrivilege(clusterName string, privilege *Privilege) (*Privilege, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	if privilege == nil {
		panic("Privilege can't be nil")
	}

	// Update the privilege
	path := fmt.Sprintf("/clusters/%s/privileges/%d", clusterName, privilege.PrivilegeInfo.PrivilegeId)
	jsonData, err := json.Marshal(privilege)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Put(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to update: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, errors.New(resp.Status())
	}

	return privilege, nil

}

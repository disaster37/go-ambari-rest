// This file permit to manage credential in Ambari cluster
// Ambari documentation:https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/credential-resources.md

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Credential object
type Credential struct {
	CredentialInfo *CredentialInfo `json:"Credential"`
}
type CredentialResponse struct {
	Response
	Items []Credential `json:"items"`
}
type CredentialInfo struct {
	Alias       string `json:"alias,omitempty"`
	ClusterName string `json:"cluster_name,omitempty"`
	Principal   string `json:"principal,omitempty"`
	Key         string `json:"key,omitempty"`
	Type        string `json:"type,omitempty"`
}

const (
	CREDENTIAL_TEMPORARY = "temporary"
	CREDENTIAL_PERSISTED = "persisted"
)

// String return credential object as Json string
func (c *Credential) String() string {
	json, _ := json.Marshal(c)
	return string(json)
}

// Clean object before save or update it
func (c *Credential) CleanBeforeSave() *Credential {

	return &Credential{
		CredentialInfo: &CredentialInfo{
			Principal: c.CredentialInfo.Principal,
			Key:       c.CredentialInfo.Key,
			Type:      c.CredentialInfo.Type,
		},
	}

}

// Credential return existing credential on cluster
// It return the credential if is found
// It return nil if not found
// It return error if something wrong when it call the API
func (c *AmbariClient) Credential(clusterName string, alias string) (*Credential, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	if alias == "" {
		panic("Alias can't be empty")
	}
	log.Debug("Alias: ", alias)

	path := fmt.Sprintf("/clusters/%s/credentials/%s", clusterName, alias)
	resp, err := c.Client().R().Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Result : ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		} else {
			return nil, NewAmbariError(resp.StatusCode(), resp.Status())
		}
	}
	credential := &Credential{}
	err = json.Unmarshal(resp.Body(), credential)
	if err != nil {
		return nil, err
	}

	log.Debug("Credential: ", credential)

	return credential, nil

}

// Credentials return all credential on cluster
// It return the list of credential.
//  If not credential, it return empty list.
// It return error if something wrong when it call the API
func (c *AmbariClient) Credentials(clusterName string) ([]Credential, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)

	path := fmt.Sprintf("/clusters/%s/credentials", clusterName)
	resp, err := c.Client().R().Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Result : ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		} else {
			return nil, NewAmbariError(resp.StatusCode(), resp.Status())
		}
	}
	credentialResponse := &CredentialResponse{}
	err = json.Unmarshal(resp.Body(), credentialResponse)
	if err != nil {
		return nil, err
	}

	log.Debug("Credentials: ", credentialResponse.Items)

	return credentialResponse.Items, nil

}

// CreateCredential permit to create new credential on cluster
// It return the credential if all work fine
// It return error if something wrong when it call the API
func (c *AmbariClient) CreateCredential(credential *Credential) (*Credential, error) {

	if credential == nil {
		panic("Credential can't be null")
	}
	if credential.CredentialInfo.ClusterName == "" {
		panic("ClusterName can't be empty")
	}
	if credential.CredentialInfo.Alias == "" {
		panic("Alias can't be empty")
	}
	log.Debug("Credential: ", credential)

	// Create the credential
	path := fmt.Sprintf("/clusters/%s/credentials/%s", credential.CredentialInfo.ClusterName, credential.CredentialInfo.Alias)

	credentialPayload := credential.CleanBeforeSave()
	log.Debug("Credential payload: ", credentialPayload)
	jsonData, err := json.Marshal(credentialPayload)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	// Get the credential
	credential, err = c.Credential(credential.CredentialInfo.ClusterName, credential.CredentialInfo.Alias)
	if err != nil {
		return nil, err
	}
	if credential == nil {
		return nil, NewAmbariError(500, "Can't get credential that just created")
	}

	return credential, err

}

// DeleteCredential permit to delete existing credential on cluster
// It return error if something wrong when it call the API
func (c *AmbariClient) DeleteCredential(clusterName string, alias string) error {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}
	log.Debug("ClusterName: ", clusterName)
	if alias == "" {
		panic("Alias can't be empty")
	}
	log.Debug("Alias: ", alias)

	path := fmt.Sprintf("/clusters/%s/credentials/%s", clusterName, alias)
	resp, err := c.Client().R().Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response to delete credential: ", resp)
	if resp.StatusCode() >= 300 {
		return NewAmbariError(resp.StatusCode(), resp.Status())
	}

	return nil

}

// UpdateCredential permit to update existing credential
// It return the credential if all work fine
// It return error if something wrong when it call the API
func (c *AmbariClient) UpdateCredential(credential *Credential) (*Credential, error) {

	if credential == nil {
		panic("Credential can't be nil")
	}
	if credential.CredentialInfo.ClusterName == "" {
		panic("ClusterName can't be empty")
	}
	if credential.CredentialInfo.Alias == "" {
		panic("Alias can't be empty")
	}
	log.Debug("Credential: ", credential)

	// Update the credential
	path := fmt.Sprintf("/clusters/%s/credentials/%s", credential.CredentialInfo.ClusterName, credential.CredentialInfo.Alias)
	credentialPayload := credential.CleanBeforeSave()
	log.Debug("Credential payload: ", credentialPayload)
	jsonData, err := json.Marshal(credentialPayload)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client().R().SetBody(jsonData).Put(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to update: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	// Get the credential after update
	credential, err = c.Credential(credential.CredentialInfo.ClusterName, credential.CredentialInfo.Alias)
	if err != nil {
		return nil, err
	}
	if credential == nil {
		return nil, NewAmbariError(500, "Can't get credential that just created")
	}

	return credential, err

}

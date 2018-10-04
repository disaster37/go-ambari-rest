// Ambari documentation: https://cwiki.apache.org/confluence/display/AMBARI/Blueprints
// This file permit to manage blueprint
// Blueprint permit to deploy HDP cluster from deployement plan

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Blueprint Json object
type Blueprint struct {
	Configurations []map[string]map[string]map[string]string `json:"configurations"`
	HostGroups     []HostGroup                               `json:"host_groups"`
	BlueprintInfo  BlueprintInfo                             `json:"Blueprints"`
}
type HostGroup struct {
	Components     []map[string]string                       `json:"components"`
	Configurations []map[string]map[string]map[string]string `json:"configurations"`
	Name           string                                    `json:"name"`
	Cardinality    string                                    `json:"cardinality"`
}
type BlueprintInfo struct {
	Name    string `json:"blueprint_name,omitempty"`
	Stack   string `json:"stack_name"`
	Version string `json:"stack_version"`
}

// String return blueprint object as Json string
func (b *Blueprint) String() string {
	json, _ := json.Marshal(b)
	return string(json)
}

// CreateBlueprint permit to create new blueprint item
// It return blueprint object if all work fine.
// It return error if something wrong when call the API
func (c *AmbariClient) CreateBlueprint(name string, jsonBlueprint string) (*Blueprint, error) {

	if name == "" {
		panic("Name can't be empty")
	}
	if jsonBlueprint == "" {
		panic("JsonBlueprint can't be empty")
	}
	log.Debugf("Name: %s", name)
	log.Debugf("JsonBlueprint: %s", jsonBlueprint)

	var blueprintTest interface{}
	err := json.Unmarshal([]byte(jsonBlueprint), &blueprintTest)
	if err != nil {
		return nil, err
	}

	// Create the BluePrint
	path := fmt.Sprintf("/blueprints/%s", name)
	resp, err := c.Client().R().SetBody(jsonBlueprint).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to create: ", resp)
	if resp.StatusCode() >= 300 {
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	blueprint, err := c.Blueprint(name)
	if err != nil {
		return nil, err
	}
	if blueprint == nil {
		return nil, NewAmbariError(500, "Can't get blueprint that just created")
	}

	log.Debugf("Return blueprint: %s", blueprint)

	return blueprint, nil

}

// Blueprint permit to get blueprint item from is name
// It return blueprint object if exist, else it return nil
// It return error if something wrong when call the API
func (c *AmbariClient) Blueprint(name string) (*Blueprint, error) {

	if name == "" {
		panic("Name can't be empty")
	}
	log.Debug("Name: ", name)

	path := fmt.Sprintf("/blueprints/%s", name)
	resp, err := c.Client().R().Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to get: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		} else {
			return nil, NewAmbariError(resp.StatusCode(), resp.Status())
		}
	}
	blueprint := &Blueprint{}
	err = json.Unmarshal(resp.Body(), blueprint)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return blueprint: %s", blueprint)

	return blueprint, nil
}

// DeleteBlueprint permit to delete blueprint item
// It return error if blueprint item not exist or if something wrong when call the API
func (c *AmbariClient) DeleteBlueprint(name string) error {

	if name == "" {
		panic("Name can't be empty")
	}
	log.Debug("Name: ", name)

	// Check if blueprint exist
	blueprint, err := c.Blueprint(name)
	if err != nil {
		return err
	}
	if blueprint == nil {
		return NewAmbariError(404, "Blueprint %s not found", name)
	}

	path := fmt.Sprintf("/blueprints/%s", name)
	resp, err := c.Client().R().Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response to delete blueprint: ", resp)
	if resp.StatusCode() >= 300 {
		return NewAmbariError(resp.StatusCode(), resp.Status())
	}

	return nil

}

package client

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Ambari documentation: https://cwiki.apache.org/confluence/display/AMBARI/Blueprints

type Blueprint struct {
	Configurations []map[string]map[string]map[string]string `json:"configurations"`
	HostGroups     []HostGroup                               `json:"host_groups"`
	BlueprintInfo  BlueprintInfo                             `json:"Blueprints"`
}

type HostGroup struct {
	Components     []map[string]string                       `json:"components"`
	Configurations []map[string]map[string]map[string]string `json:"configurations"`
	Name           string                                    `json:"name"`
	Cardinality    int                                       `json:"cardinality"`
}

type BlueprintInfo struct {
	Name    string `json:"blueprint_name,omitempty"`
	Stack   string `json:"stack_name"`
	Version string `json:"stack_version"`
}

func (b *Blueprint) String() string {
	json, _ := json.Marshal(b)
	return string(json)
}

// Create new blueprint
func (c *AmbariClient) CreateBlueprint(name string, jsonBlueprint string) (*Blueprint, error) {

	if name == "" {
		panic("Name can't be empty")
	}
	if jsonBlueprint == "" {
		panic("JsonBlueprint can't be empty")
	}
	log.Debug("Name: %s", name)
	log.Debug("JsonBlueprint: %s", jsonBlueprint)
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
		return nil, errors.New(resp.Status())
	}

	blueprint, err := c.Blueprint(name)
	if err != nil {
		return nil, err
	}
	if blueprint == nil {
		return nil, errors.New("Can't get blueprint that just created")
	}

	log.Debug("Return blueprint: %s", blueprint)

	return blueprint, nil

}

// Get blueprint
func (c *AmbariClient) Blueprint(name string) (*Blueprint, error) {

	if name == "" {
		panic("Name can't be empty")
	}

	path := fmt.Sprintf("/blueprints/%s", name)
	resp, err := c.Client().R().Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to get: ", resp)
	blueprint := &Blueprint{}
	err = json.Unmarshal(resp.Body(), blueprint)
	if err != nil {
		return nil, err
	}
	log.Debug("Return blueprint: %s", blueprint)

	return blueprint, nil
}

func (c *AmbariClient) DeleteBlueprint(name string) error {

	if name == "" {
		panic("Name can't be empty")
	}
	path := fmt.Sprintf("/blueprints/%s", name)

	resp, err := c.Client().R().Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response to delete host: ", resp)
	if resp.StatusCode() >= 300 {
		return errors.New(resp.Status())
	}

	return nil

}

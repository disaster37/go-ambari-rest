// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/alerts.md
// This file permit to manager Alert item on Ambari API

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Alert struct {
	AlertInfo *AlertInfo `json:"Alert"`
}

type AlertInfo struct {
	ClusterName      string `json:"cluster_name,omitempty"`
	ServiceName      string `json:"service_name,omitempty"`
	ComponentName    string `json:"component_name,omitempty"`
	Hostname         string `json:"host_name,omitempty"`
	Label            string `json:"label,omitempty"`
	MaintenanceState string `json:"maintenance_state,omitempty"`
	State            string `json:"state,omitempty"`
	Text             string `json:"text,omitempty"`
}

type Alerts struct {
	Items []Alert `json:"items,omitempty"`
}

// String permit to return Alert object as Json string
func (a *Alert) String() string {
	json, _ := json.Marshal(a)
	return string(json)
}

func (c *AmbariClient) AlertsInHost(clusterName string, hostname string) ([]Alert, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}

	if hostname == "" {
		panic("Hostname can't be empty")
	}

	log.Debug("ClusterName: ", clusterName)
	log.Debug("Hostname: ", hostname)

	// Check if host exist
	host, err := c.HostOnCluster(clusterName, hostname)
	if err != nil {
		return nil, err
	}
	if host == nil {
		return nil, NewAmbariError(404, "Host %s not found in cluster", hostname)
	}

	query := "fields=*&Alert/maintenance_state=OFF"
	path := fmt.Sprintf("/clusters/%s/hosts/%s/alerts", clusterName, hostname)
	resp, err := c.Client().R().SetQueryString(query).Get(path)
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
	alertsTemp := &Alerts{}
	err = json.Unmarshal(resp.Body(), alertsTemp)
	if err != nil {
		return nil, err
	}

	// Keep only alert
	alerts, err := c.filterAlerts(alertsTemp.Items)
	if err != nil {
		return nil, err
	}

	log.Debugf("Return alerts: %s", alerts)

	return alerts, nil
}

func (c *AmbariClient) AlertsInService(clusterName string, serviceName string) ([]Alert, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}

	if serviceName == "" {
		panic("Serviceame can't be empty")
	}

	log.Debug("ClusterName: ", clusterName)
	log.Debug("ServiceName: ", serviceName)

	// Check if service exist
	service, err := c.Service(clusterName, serviceName)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, NewAmbariError(404, "Service %s not found", serviceName)
	}

	query := "fields=*&Alert/maintenance_state=OFF"

	path := fmt.Sprintf("/clusters/%s/services/%s/alerts", clusterName, serviceName)
	resp, err := c.Client().R().SetQueryString(query).Get(path)
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
	alertsTemp := &Alerts{}
	err = json.Unmarshal(resp.Body(), alertsTemp)
	if err != nil {
		return nil, err
	}

	// Keep only alert
	alerts, err := c.filterAlerts(alertsTemp.Items)
	if err != nil {
		return nil, err
	}

	log.Debugf("Return alerts: %s", alerts)

	return alerts, nil
}

func (c *AmbariClient) AlertsInCluster(clusterName string) ([]Alert, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}

	log.Debug("ClusterName: ", clusterName)

	path := fmt.Sprintf("/clusters/%s/alerts", clusterName)
	query := "fields=*&Alert/maintenance_state=OFF"

	resp, err := c.Client().R().SetQueryString(query).Get(path)
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
	alertsTemp := &Alerts{}
	err = json.Unmarshal(resp.Body(), alertsTemp)
	if err != nil {
		return nil, err
	}

	// Keep only alert
	alerts, err := c.filterAlerts(alertsTemp.Items)
	if err != nil {
		return nil, err
	}

	log.Debugf("Return alerts: %s", alerts)

	return alerts, nil
}

func (c *AmbariClient) Alerts(clusterName string) ([]Alert, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}

	log.Debug("ClusterName: ", clusterName)

	path := fmt.Sprintf("/clusters/%s/alerts", clusterName)
	query := "fields=*&Alert/maintenance_state=OFF"
	resp, err := c.Client().R().SetQueryString(query).Get(path)
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
	alerts := &Alerts{}
	err = json.Unmarshal(resp.Body(), alerts)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return alerts: %s", alerts)

	return alerts.Items, nil
}

// filterAlerts permet to keep only WARNING and CRITICAL alerts
// It's return []AlertInfo if all work fine
// It's return error if somthing wrong
func (c *AmbariClient) filterAlerts(alerts []Alert) ([]Alert, error) {

	resultAlerts := make([]Alert, 0, 1)

	for _, alert := range alerts {
		if alert.AlertInfo.State == "WARNING" || alert.AlertInfo.State == "CRITICAL" || alert.AlertInfo.State == "UNKNOWN" {
			resultAlerts = append(resultAlerts, alert)
		}
	}

	return resultAlerts, nil

}

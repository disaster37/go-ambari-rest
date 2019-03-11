// This file permit to manage Request in Ambari API
// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/request-resources.md

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const (
	REQUEST_FAILED    = "FAILED"
	REQUEST_ACCEPTED  = "ACCEPTED"
	REQUEST_COMPLETED = "COMPLETED"
	REQUEST_ABORDED   = "ABORDED"
)

type RequestTask struct {
	RequestTaskInfo *RequestTaskInfo `json:"Requests,omitempty"`
}

type RequestTaskInfo struct {
	Id              int     `json:"id,omitempty"`
	CompletedTask   int     `json:"completed_task_count,omitempty"`
	AbordedTask     int     `json:"aborted_task_count,omitempty"`
	FailedTask      int     `json:"failed_task_count,omitempty"`
	TaskCount       int     `json:"task_count,omitempty"`
	ProgressPercent float64 `json:"progress_percent,omitempty"`
	Status          string  `json:"request_status,omitempty"`
	Context         string  `json:"request_context,omitempty"`
}

type RequestsTask struct {
	Items []RequestTask `json:"Items"`
}

// String permit to get Request object as Json string
func (r *RequestTask) String() string {
	json, _ := json.Marshal(r)
	return string(json)
}

// String permit to get Request object as Json string
func (r *RequestsTask) String() string {
	json, _ := json.Marshal(r)
	return string(json)
}

// Request permit to get request by is name
// It return RequestTask if is found
// It return nil is request is not found
// It return error if something wrong with the API call
func (c *AmbariClient) Request(clusterName string, Id int) (*RequestTask, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}

	log.Debug("ClusterName: ", clusterName)
	log.Debug("Id: ", Id)

	path := fmt.Sprintf("/clusters/%s/requests/%d", clusterName, Id)
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
	requestTask := &RequestTask{}
	err = json.Unmarshal(resp.Body(), requestTask)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return requestTask: %s", requestTask)

	return requestTask, nil
}

// Requests permit to get all requests
// It return the list of requestTask
// It return empty list if there are no tasks
// It return error if something wrong with the API call
func (c *AmbariClient) Requests(clusterName string) ([]RequestTask, error) {

	if clusterName == "" {
		panic("ClusterName can't be empty")
	}

	log.Debug("ClusterName: ", clusterName)

	path := fmt.Sprintf("/clusters/%s/requests?fields=*", clusterName)
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
	requestsTask := &RequestsTask{}
	err = json.Unmarshal(resp.Body(), requestsTask)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return requestsTask: %s", requestsTask)

	return requestsTask.Items, nil
}

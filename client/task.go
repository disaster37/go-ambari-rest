// This file permit to manage Request in Ambari API
// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/request-resources.md

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
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
	ClusterName     string  `json:"cluster_name,omitempty"`
}

type RequestsTask struct {
	Items []RequestTask `json:"Items"`
}

// String permit to get Request object as Json string
func (r *RequestTask) String() string {
	json, _ := json.Marshal(r)
	return string(json)
}

// Permit to wait the rerquest task is finished
// It can return error if API call failed
func (r *RequestTask) Wait(c *AmbariClient, clusterName string) error {
	if r.RequestTaskInfo != nil {
		isRun := true
		for isRun {
			requestTask, err := c.Request(clusterName, r.RequestTaskInfo.Id)
			if err != nil {
				return err
			}
			*r = *requestTask
			if r.RequestTaskInfo.ProgressPercent < 100 {
				log.Debugf("Task '%s' (%d) is not yet finished, state is %s (%f %%)", r.RequestTaskInfo.Context, r.RequestTaskInfo.Id, r.RequestTaskInfo.Status, r.RequestTaskInfo.ProgressPercent)
				time.Sleep(10 * time.Second)
			} else {
				isRun = false
			}
		}

		log.Debugf("Task '%s' (%d) is finished with state %s", r.RequestTaskInfo.Context, r.RequestTaskInfo.Id, r.RequestTaskInfo.Status)
	} else {
		log.Debugf("Task is empty...")
	}

	return nil

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

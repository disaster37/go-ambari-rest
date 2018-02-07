package client

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/repository-version-resources.md

type Repository struct {
	Response
	RepositoryVersion RepositoryVersion `json:"RepositoryVersions"`
	OS                []OS              `json:"operating_systems"`
}

type RepositoryVersion struct {
	Id           int    `json:"id,omitempty"`
	Version      string `json:"repository_version"`
	Name         string `json:"display_name"`
	StackName    string `json:"stack_name,omitempty"`
	StackVersion string `json:"stack_version,omitempty"`
}

type OS struct {
	OSInfo           OSInfo           `json:"OperatingSystems"`
	RepositoriesInfo []RepositoryInfo `json:"repositories"`
}

type OSInfo struct {
	Type string `json:"os_type"`
}

type RepositoryInfo struct {
	Id      string `json:"repo_id"`
	Name    string `json:"repo_name"`
	BaseUrl string `json:"base_url"`
}

type RepositoriesResponse struct {
	Response
	Items []Repository `json:"items"`
}

func (r *Repository) String() string {
	json, _ := json.Marshal(r)
	return string(json)
}

func (c *AmbariClient) CreateRepository(repository *Repository) (*Repository, error) {

	if repository == nil {
		panic("Repository can't be nil")
	}

	log.Debug("Repository: %s", repository.String())

	path := fmt.Sprintf("/stacks/%s/versions/%s/repository_versions", repository.RepositoryVersion.StackName, repository.RepositoryVersion.StackVersion)
	jsonData, err := json.Marshal(repository)
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

	repository, err = c.SearchRepository(repository.RepositoryVersion.StackName, repository.RepositoryVersion.StackVersion, repository.RepositoryVersion.Name, repository.RepositoryVersion.Version)
	if err != nil {
		return nil, err
	}
	if repository == nil {
		return nil, errors.New("Can't get repository that just created")
	}

	log.Debug("Return repository: %s", repository)

	return host, nil

}

// Get cluster by ID is not supported by ambari api
func (c *AmbariClient) Repository(stackName string, stackVersion string, repositoryId int) (*Repository, error) {

	if stackName == "" {
		panic("StackName can't be empty")
	}
	if stackVersion == "" {
		panic("StackVersion can't be empty")
	}

	path := fmt.Sprintf("/stacks/%s/versions/%s/repository_versions/%d", stackName, stackVersion, repositoryId)

	// Get the host components
	resp, err := c.Client().R().Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to get: ", resp)
	repository := &Repository{}
	err = json.Unmarshal(resp.Body(), host)
	if err != nil {
		return nil, err
	}
	log.Debug("Return repository: %s", host)

	return repository, nil
}

func (c *AmbariClient) UpdateRepository(repository *Repository) (*Repository, error) {

	if repository == nil {
		panic("Repository can't be nil")
	}
	log.Debug("Repository: ", repository)

	path := fmt.Sprintf("/stacks/%s/versions/%s/repository_versions/%d", repository.RepositoryVersion.StackName, repository.RepositoryVersion.StackVersion, repository.RepositoryVersion.Id)
	jsonData, err := json.Marshal(host)
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

	// Get the Host
	repository, err = c.Repository(repository.RepositoryVersion.StackName, repository.RepositoryVersion.StackVersion, repository.RepositoryVersion.Id)
	if err != nil {
		return nil, err
	}
	if repository == nil {
		return nil, errors.New("Can't get repository that just updated")
	}

	log.Debug("Return repository: %s", repository.String())

	return repository, nil

}

func (c *AmbariClient) DeleteRepository(stackName string, stackVersion string, repositoryId int) error {

	if stackName == "" {
		panic("StackName can't be empty")
	}
	if stackVersion == "" {
		panic("StackVersion can't be empty")
	}
	path := fmt.Sprintf("/stacks/%s/versions/%s/repository_versions/%d", stackName, stackVersion, repositoryId)

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

func (c *AmbariClient) SearchRepository(stackName string, stackVersion string, repositoryName string, repositoryVersion string) (*Repository, error) {

	if stackName == "" {
		panic("StackName can't be empty")
	}
	if stackVersion == "" {
		panic("StackVersion can't be empty")
	}
	if repositoryName == "" {
		panic("RepositoryName can't be empty")
	}
	if repositoryVersion == "" {
		panic("RepositoryVersion can't be empty")
	}

	path := fmt.Sprintf("/stacks/%s/versions/%s/repository_versions", stackName, stackVersion)

	resp, err := c.Client().R().SetQueryParams(map[string]string{
		"RepositoryVersions/repository_version": repositoryVersion,
		"RepositoryVersions/display_name":       repositoryName,
	}).Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response to get: ", resp)
	repositoryResponse := &RepositoriesResponse{}
	err = json.Unmarshal(resp.Body(), repositoryResponse)
	if err != nil {
		return nil, err
	}
	log.Debug("RepositoryResponse: ", repositoryResponse)

	if len(repositoryResponse.Items) > 0 {
		log.Debug("Repository: ", repositoryResponse.Items[0])
		return &repositoryResponse.Items[0], nil
	} else {
		return nil, nil
	}
}

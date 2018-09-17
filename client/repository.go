// This file permit to manage repository on Ambari API
// Ambari documentation: https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/repository-version-resources.md

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Repository object
type Repository struct {
	RepositoryVersion *RepositoryVersion `json:"RepositoryVersions"`
	OS                []OS               `json:"operating_systems"`
}
type RepositoryVersion struct {
	Id           int    `json:"id,omitempty"`
	Version      string `json:"repository_version"`
	Name         string `json:"display_name"`
	StackName    string `json:"stack_name,omitempty"`
	StackVersion string `json:"stack_version,omitempty"`
}
type OS struct {
	Response
	OSInfo           *OSInfo          `json:"OperatingSystems"`
	RepositoriesData []RepositoryData `json:"repositories"`
}
type OSInfo struct {
	Type              string `json:"os_type"`
	ManagedRepository bool   `json:"ambari_managed_repositories, omitempty"`
}
type RepositoryData struct {
	Response
	RepositoryInfo *RepositoryInfo `json:"Repositories"`
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

// String return repository object as Json string
func (r *Repository) String() string {
	json, _ := json.Marshal(r)
	return string(json)
}

// CleanBeforeSave permit to clean some attributes before to save or update repository
func (r *Repository) CleanBeforeSave() {
	// Remove href before send
	for index, os := range r.OS {
		os.Href = nil

		for index2, repositoryData := range os.RepositoriesData {
			repositoryData.Href = nil
			os.RepositoriesData[index2] = repositoryData
		}

		r.OS[index] = os
	}
}

// CreateRepository permit to create new repository
// It return repository if all work fine
// It return error if something wrong when it call the API
func (c *AmbariClient) CreateRepository(repository *Repository) (*Repository, error) {

	if repository == nil {
		panic("Repository can't be nil")
	}
	log.Debug("Repository: %s", repository.String())

	repository.CleanBeforeSave()

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
		return nil, NewAmbariError(resp.StatusCode(), resp.Status())
	}

	repository, err = c.SearchRepository(repository.RepositoryVersion.StackName, repository.RepositoryVersion.StackVersion, repository.RepositoryVersion.Name, repository.RepositoryVersion.Version)
	if err != nil {
		return nil, err
	}
	if repository == nil {
		return nil, NewAmbariError(500, "Can't get repository that just created")
	}

	log.Debug("Return repository: %s", repository)

	return repository, nil

}

// Repository permit to get existing repository
// It return repository if is found
// It return nil if is not found
func (c *AmbariClient) Repository(stackName string, stackVersion string, repositoryId int) (*Repository, error) {

	if stackName == "" {
		panic("StackName can't be empty")
	}
	if stackVersion == "" {
		panic("StackVersion can't be empty")
	}
	log.Debug("StackName: ", stackName)
	log.Debug("StackVersion: ", stackVersion)

	path := fmt.Sprintf("/stacks/%s/versions/%s/repository_versions/%d", stackName, stackVersion, repositoryId)
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
	repository := &Repository{}
	err = json.Unmarshal(resp.Body(), repository)
	if err != nil {
		return nil, err
	}

	log.Debug("Repository : ", repository.String())

	// Get Repositories for each OS
	if repository != nil {
		for index, os := range repository.OS {
			log.Debug("Call ", os.Href)
			resp, err = c.Client().R().Get(*os.Href)
			if err != nil {
				return nil, err
			}
			log.Debug("Response to get repository: ", resp)
			os = OS{}
			err = json.Unmarshal(resp.Body(), &os)
			if err != nil {
				return nil, err
			}
			log.Debug("Return repository: ", os)

			for index2, repositoryData := range os.RepositoriesData {
				log.Debug("Call ", repositoryData.Href)
				resp, err = c.Client().R().Get(*repositoryData.Href)
				if err != nil {
					return nil, err
				}
				log.Debug("Response to get repositoryData: ", resp)
				repositoryData = RepositoryData{}
				err = json.Unmarshal(resp.Body(), &repositoryData)
				if err != nil {
					return nil, err
				}
				log.Debug("Return repositoryData: ", repositoryData)
				os.RepositoriesData[index2] = repositoryData
			}

			repository.OS[index] = os
		}
	}
	log.Debug("Return repository: %s", repository)
	return repository, nil
}

// UpdateRepository permit to update an existing repository like a repo URL
// It return updated repository if all work fine
// It return nil if something wrong when it call the API
func (c *AmbariClient) UpdateRepository(repository *Repository) (*Repository, error) {

	if repository == nil {
		panic("Repository can't be nil")
	}
	log.Debug("Repository: ", repository)

	repository.CleanBeforeSave()

	path := fmt.Sprintf("/stacks/%s/versions/%s/repository_versions/%d", repository.RepositoryVersion.StackName, repository.RepositoryVersion.StackVersion, repository.RepositoryVersion.Id)
	jsonData, err := json.Marshal(repository)
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

	repository, err = c.Repository(repository.RepositoryVersion.StackName, repository.RepositoryVersion.StackVersion, repository.RepositoryVersion.Id)
	if err != nil {
		return nil, err
	}
	if repository == nil {
		return nil, NewAmbariError(500, "Can't get repository that just updated")
	}

	log.Debug("Return repository: %s", repository.String())

	return repository, nil

}

// DeleteRepository permit to remove the repository
// You can't delete the current repository use by the cluster
// It return error if something wrong when it call the API
func (c *AmbariClient) DeleteRepository(stackName string, stackVersion string, repositoryId int) error {

	if stackName == "" {
		panic("StackName can't be empty")
	}
	if stackVersion == "" {
		panic("StackVersion can't be empty")
	}
	log.Debug("StackName: ", stackName)
	log.Debug("StackVersion: ", stackVersion)

	path := fmt.Sprintf("/stacks/%s/versions/%s/repository_versions/%d", stackName, stackVersion, repositoryId)
	resp, err := c.Client().R().Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response to delete host: ", resp)
	if resp.StatusCode() >= 300 {
		return NewAmbariError(resp.StatusCode(), resp.Status())
	}

	return nil

}

// SearchRepository permit to get existing repository by is name
// It return the repository if is found
// It return nil if is not found
// It return error if something wrong when it call the API
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
	log.Debug("StackName: ", stackName)
	log.Debug("StackVersion: ", stackVersion)
	log.Debug("RepositoryName ", repositoryName)
	log.Debug("RepositoryVersion ", repositoryVersion)

	path := fmt.Sprintf("/stacks/%s/versions/%s/repository_versions", stackName, stackVersion)
	resp, err := c.Client().R().SetQueryParams(map[string]string{
		"RepositoryVersions/repository_version": repositoryVersion,
		"RepositoryVersions/display_name":       repositoryName,
	}).Get(path)
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
	repositoryResponse := &RepositoriesResponse{}
	err = json.Unmarshal(resp.Body(), repositoryResponse)
	if err != nil {
		return nil, err
	}
	log.Debug("RepositoryResponse: ", repositoryResponse)

	if len(repositoryResponse.Items) > 0 {
		log.Debug("Repository: ", repositoryResponse.Items[0])

		repository, err := c.Repository(stackName, stackVersion, repositoryResponse.Items[0].RepositoryVersion.Id)
		if err != nil {
			return nil, err
		}

		return repository, nil
	} else {
		return nil, nil
	}
}

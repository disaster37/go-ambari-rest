package main

import (
	"encoding/json"
	"github.com/disaster37/go-ambari-rest/client"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
)

type RepositoryStack struct {
	StackName        string `json:"stackName"`
	StackVersion     string `json:"stackVersion"`
	Name             string `json:"name"`
	Version          string `json:"version"`
	OperatingSystems []OS   `json:"operatingSystems"`
}

type OS struct {
	OsName       string       `json:"osName"`
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	RepositoryId      string `json:"repositoryId"`
	RepositoryName    string `json:"repositoryName"`
	RepositoryBaseUrl string `json:"repositoryBaseUrl"`
}

func createRepository(c *cli.Context) error {

	clientAmbari, err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if c.String("repository-file") == "" {
		return cli.NewExitError("You must set --repository-file parameter", 1)
	}

	// Read the Json file
	b, err := ioutil.ReadFile(c.String("repository-file"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	repositoryStackJson := string(b)
	log.Debug("Repository: ", repositoryStackJson)
	repositoryStack := &RepositoryStack{}
	err = json.Unmarshal(b, repositoryStack)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// Check if repository already exist
	repository, err := clientAmbari.SearchRepository(repositoryStack.StackName, repositoryStack.StackVersion, repositoryStack.Name, repositoryStack.Version)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// Create the target struct
	repositoryTarget := &client.Repository{
		RepositoryVersion: &client.RepositoryVersion{
			Version:      repositoryStack.Version,
			Name:         repositoryStack.Name,
			StackName:    repositoryStack.StackName,
			StackVersion: repositoryStack.StackVersion,
		},
		OS: make([]client.OS, 0, 1),
	}
	for _, osTemp := range repositoryStack.OperatingSystems {
		os := client.OS{
			OSInfo: &client.OSInfo{
				Type: osTemp.OsName,
			},
			RepositoriesData: make([]client.RepositoryData, 0, 2),
		}
		for _, repositoryTemp := range osTemp.Repositories {
			repositoryData := client.RepositoryData{
				RepositoryInfo: &client.RepositoryInfo{
					Id:      repositoryTemp.RepositoryId,
					Name:    repositoryTemp.RepositoryName,
					BaseUrl: repositoryTemp.RepositoryBaseUrl,
				},
			}
			os.RepositoriesData = append(os.RepositoriesData, repositoryData)
		}
		repositoryTarget.OS = append(repositoryTarget.OS, os)
	}

	// Create new repository
	if repository == nil {
		_, err = clientAmbari.CreateRepository(repositoryTarget)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		log.Info("Repository created successfully")
	} else {
		// Update Repository
		repositoryTarget.RepositoryVersion.Id = repository.RepositoryVersion.Id
		_, err = clientAmbari.UpdateRepository(repositoryTarget)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		log.Info("Repository updated successfully")
	}

	return nil

}

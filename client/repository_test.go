package client

import (
	"github.com/stretchr/testify/assert"
)

func (s *ClientTestSuite) TestRepository() {

	// Create repository
	repository := &Repository{
		RepositoryVersion: &RepositoryVersion{
			Version:      "2.6.4.0.1",
			Name:         "HDP-2.6.4.0.1",
			StackName:    "HDP",
			StackVersion: "2.6",
		},
		OS: []OS{
			OS{
				OSInfo: &OSInfo{
					Type:              "redhat7",
					ManagedRepository: true,
				},
				RepositoriesData: []RepositoryData{
					RepositoryData{
						RepositoryInfo: &RepositoryInfo{
							Id:      "HDP",
							Name:    "HDP",
							BaseUrl: "http://public-repo-1.hortonworks.com/HDP/centos7/2.x/updates/2.6.4.0",
						},
					},
					RepositoryData{
						RepositoryInfo: &RepositoryInfo{
							Id:      "HDP-UTILS",
							Name:    "HDP-UTILS",
							BaseUrl: "http://public-repo-1.hortonworks.com/HDP-UTILS-1.1.0.22/repos/centos7",
						},
					},
				},
			},
		},
	}
	repository, err := s.client.CreateRepository(repository)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), repository)
	if repository != nil {
		assert.NotEqual(s.T(), "", repository.RepositoryVersion.Id)
		assert.Equal(s.T(), "2.6.4.0.1", repository.RepositoryVersion.Version)
		assert.Equal(s.T(), "HDP-2.6.4.0.1", repository.RepositoryVersion.Name)
		assert.Equal(s.T(), "HDP", repository.RepositoryVersion.StackName)
		assert.Equal(s.T(), "2.6", repository.RepositoryVersion.StackVersion)
		assert.Equal(s.T(), "redhat7", repository.OS[0].OSInfo.Type)
		assert.Equal(s.T(), true, repository.OS[0].OSInfo.ManagedRepository)
		assert.Equal(s.T(), "HDP", repository.OS[0].RepositoriesData[0].RepositoryInfo.Id)
		assert.Equal(s.T(), "HDP", repository.OS[0].RepositoriesData[0].RepositoryInfo.Name)
		assert.Equal(s.T(), "http://public-repo-1.hortonworks.com/HDP/centos7/2.x/updates/2.6.4.0", repository.OS[0].RepositoriesData[0].RepositoryInfo.BaseUrl)
		assert.Equal(s.T(), "HDP-UTILS", repository.OS[0].RepositoriesData[1].RepositoryInfo.Id)
		assert.Equal(s.T(), "HDP-UTILS", repository.OS[0].RepositoriesData[1].RepositoryInfo.Name)
		assert.Equal(s.T(), "http://public-repo-1.hortonworks.com/HDP-UTILS-1.1.0.22/repos/centos7", repository.OS[0].RepositoriesData[1].RepositoryInfo.BaseUrl)
	}

	// Get Repository
	repository, err = s.client.Repository("HDP", "2.6", repository.RepositoryVersion.Id)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), repository)
	if repository != nil {
		assert.NotEqual(s.T(), "", repository.RepositoryVersion.Id)
		assert.Equal(s.T(), "2.6.4.0.1", repository.RepositoryVersion.Version)
		assert.Equal(s.T(), "HDP-2.6.4.0.1", repository.RepositoryVersion.Name)
		assert.Equal(s.T(), "HDP", repository.RepositoryVersion.StackName)
		assert.Equal(s.T(), "2.6", repository.RepositoryVersion.StackVersion)
		assert.Equal(s.T(), "redhat7", repository.OS[0].OSInfo.Type)
		assert.Equal(s.T(), true, repository.OS[0].OSInfo.ManagedRepository)
		assert.Equal(s.T(), "HDP", repository.OS[0].RepositoriesData[0].RepositoryInfo.Id)
		assert.Equal(s.T(), "HDP", repository.OS[0].RepositoriesData[0].RepositoryInfo.Name)
		assert.Equal(s.T(), "http://public-repo-1.hortonworks.com/HDP/centos7/2.x/updates/2.6.4.0", repository.OS[0].RepositoriesData[0].RepositoryInfo.BaseUrl)
		assert.Equal(s.T(), "HDP-UTILS", repository.OS[0].RepositoriesData[1].RepositoryInfo.Id)
		assert.Equal(s.T(), "HDP-UTILS", repository.OS[0].RepositoriesData[1].RepositoryInfo.Name)
		assert.Equal(s.T(), "http://public-repo-1.hortonworks.com/HDP-UTILS-1.1.0.22/repos/centos7", repository.OS[0].RepositoriesData[1].RepositoryInfo.BaseUrl)
	}

	// Search Repository
	repository, err = s.client.SearchRepository("HDP", "2.6", "HDP-2.6.4.0.1", "2.6.4.0.1")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), repository)
	if repository != nil {
		assert.NotEqual(s.T(), "", repository.RepositoryVersion.Id)
		assert.Equal(s.T(), "2.6.4.0.1", repository.RepositoryVersion.Version)
		assert.Equal(s.T(), "HDP-2.6.4.0.1", repository.RepositoryVersion.Name)
		assert.Equal(s.T(), "HDP", repository.RepositoryVersion.StackName)
		assert.Equal(s.T(), "2.6", repository.RepositoryVersion.StackVersion)
		assert.Equal(s.T(), "redhat7", repository.OS[0].OSInfo.Type)
		assert.Equal(s.T(), true, repository.OS[0].OSInfo.ManagedRepository)
		assert.Equal(s.T(), "HDP", repository.OS[0].RepositoriesData[0].RepositoryInfo.Id)
		assert.Equal(s.T(), "HDP", repository.OS[0].RepositoriesData[0].RepositoryInfo.Name)
		assert.Equal(s.T(), "http://public-repo-1.hortonworks.com/HDP/centos7/2.x/updates/2.6.4.0", repository.OS[0].RepositoriesData[0].RepositoryInfo.BaseUrl)
		assert.Equal(s.T(), "HDP-UTILS", repository.OS[0].RepositoriesData[1].RepositoryInfo.Id)
		assert.Equal(s.T(), "HDP-UTILS", repository.OS[0].RepositoriesData[1].RepositoryInfo.Name)
		assert.Equal(s.T(), "http://public-repo-1.hortonworks.com/HDP-UTILS-1.1.0.22/repos/centos7", repository.OS[0].RepositoriesData[1].RepositoryInfo.BaseUrl)
	}

	// Update repository
	repository.RepositoryVersion.Name = "HDP-2.6.4.0.2"
	repository, err = s.client.UpdateRepository(repository)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), repository)
	if repository != nil {
		assert.NotEqual(s.T(), "", repository.RepositoryVersion.Id)
		assert.Equal(s.T(), "HDP-2.6.4.0.2", repository.RepositoryVersion.Name)
		assert.Equal(s.T(), "HDP", repository.RepositoryVersion.StackName)
		assert.Equal(s.T(), "2.6", repository.RepositoryVersion.StackVersion)
		assert.Equal(s.T(), "redhat7", repository.OS[0].OSInfo.Type)
		assert.Equal(s.T(), true, repository.OS[0].OSInfo.ManagedRepository)
		assert.Equal(s.T(), "HDP", repository.OS[0].RepositoriesData[0].RepositoryInfo.Id)
		assert.Equal(s.T(), "HDP", repository.OS[0].RepositoriesData[0].RepositoryInfo.Name)
		assert.Equal(s.T(), "http://public-repo-1.hortonworks.com/HDP/centos7/2.x/updates/2.6.4.0", repository.OS[0].RepositoriesData[0].RepositoryInfo.BaseUrl)
		assert.Equal(s.T(), "HDP-UTILS", repository.OS[0].RepositoriesData[1].RepositoryInfo.Id)
		assert.Equal(s.T(), "HDP-UTILS", repository.OS[0].RepositoriesData[1].RepositoryInfo.Name)
		assert.Equal(s.T(), "http://public-repo-1.hortonworks.com/HDP-UTILS-1.1.0.22/repos/centos7", repository.OS[0].RepositoriesData[1].RepositoryInfo.BaseUrl)
	}

	// Delete repository
	err = s.client.DeleteRepository("HDP", "2.6", repository.RepositoryVersion.Id)
	assert.NoError(s.T(), err)

	// Create repository with spacewalk
	repository = &Repository{
		RepositoryVersion: &RepositoryVersion{
			Version:      "2.6.4.0.2",
			Name:         "HDP-2.6.4.0.2",
			StackName:    "HDP",
			StackVersion: "2.6",
		},
		OS: []OS{
			OS{
				OSInfo: &OSInfo{
					Type:              "redhat7",
					ManagedRepository: false,
				},
				RepositoriesData: []RepositoryData{
					RepositoryData{
						RepositoryInfo: &RepositoryInfo{
							Id:      "HDP",
							Name:    "HDP",
							BaseUrl: "",
						},
					},
					RepositoryData{
						RepositoryInfo: &RepositoryInfo{
							Id:      "HDP-UTILS",
							Name:    "HDP-UTILS",
							BaseUrl: "",
						},
					},
				},
			},
		},
	}
	repository, err = s.client.CreateRepository(repository)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), repository)
	if repository != nil {
		assert.NotEqual(s.T(), "", repository.RepositoryVersion.Id)
		assert.Equal(s.T(), "2.6.4.0.2", repository.RepositoryVersion.Version)
		assert.Equal(s.T(), "HDP-2.6.4.0.2", repository.RepositoryVersion.Name)
		assert.Equal(s.T(), "HDP", repository.RepositoryVersion.StackName)
		assert.Equal(s.T(), "2.6", repository.RepositoryVersion.StackVersion)
		assert.Equal(s.T(), "redhat7", repository.OS[0].OSInfo.Type)
		assert.Equal(s.T(), false, repository.OS[0].OSInfo.ManagedRepository)
		assert.Equal(s.T(), "HDP", repository.OS[0].RepositoriesData[0].RepositoryInfo.Id)
		assert.Equal(s.T(), "HDP", repository.OS[0].RepositoriesData[0].RepositoryInfo.Name)
		assert.Equal(s.T(), "", repository.OS[0].RepositoriesData[0].RepositoryInfo.BaseUrl)
		assert.Equal(s.T(), "HDP-UTILS", repository.OS[0].RepositoriesData[1].RepositoryInfo.Id)
		assert.Equal(s.T(), "HDP-UTILS", repository.OS[0].RepositoriesData[1].RepositoryInfo.Name)
		assert.Equal(s.T(), "", repository.OS[0].RepositoriesData[1].RepositoryInfo.BaseUrl)
	}

	// Delete repository
	err = s.client.DeleteRepository("HDP", "2.6", repository.RepositoryVersion.Id)
	assert.NoError(s.T(), err)

}

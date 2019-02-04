package client

import (
	"github.com/stretchr/testify/assert"
)

func (s *ClientTestSuite) TestCredential() {

	// Create Credential
	credential := &Credential{
		CredentialInfo: &CredentialInfo{
			Alias:       "kdc.admin.credential",
			ClusterName: "test",
			Principal:   "admin@TEST.LOCAL",
			Key:         "adminadmin",
			Type:        CREDENTIAL_TEMPORARY,
		},
	}
	credential, err := s.client.CreateCredential(credential)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), credential)
	if credential != nil {
		assert.Equal(s.T(), "kdc.admin.credential", credential.CredentialInfo.Alias)
		assert.Equal(s.T(), "test", credential.CredentialInfo.ClusterName)
		assert.Equal(s.T(), CREDENTIAL_TEMPORARY, credential.CredentialInfo.Type)
	}

	// Get credential
	credential, err = s.client.Credential("test", "kdc.admin.credential")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), credential)
	if credential != nil {
		assert.Equal(s.T(), "kdc.admin.credential", credential.CredentialInfo.Alias)
		assert.Equal(s.T(), "test", credential.CredentialInfo.ClusterName)
		assert.Equal(s.T(), CREDENTIAL_TEMPORARY, credential.CredentialInfo.Type)
	}

	// Update credential
	credential.CredentialInfo.Principal = "admin2@TEST.LOCAL"
	credential.CredentialInfo.Key = "adminadmin"
	credential, err = s.client.UpdateCredential(credential)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), credential)

	// Get all credential
	credentials, err := s.client.Credentials("test")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), credentials)
	if credentials != nil {
		assert.NotEmpty(s.T(), credentials)
	}

	// Delete credential
	err = s.client.DeleteCredential("test", "kdc.admin.credential")
	assert.NoError(s.T(), err)
	credential, err = s.client.Credential("test", "kdc.admin.credential")
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), credential)
}

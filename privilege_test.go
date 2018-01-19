package client

import (
	"github.com/stretchr/testify/assert"
)

// Test the constructor
func (s *ClientTestSuite) TestPrivilege() {

	// Create privilege
	privilege := &Privilege{
		PrivilegeInfo: &PrivilegeInfo{
			PermissionName: "CLUSTER.ADMINISTRATOR",
			PrincipalName:  "hm_etl_outils",
			PrincipalType:  "GROUP",
		},
	}
	privilege, err := s.client.CreatePrivilege("sihm-test", privilege)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), privilege)
	assert.NotEqual(s.T(), "", privilege.PrivilegeInfo.PrivilegeId)
	assert.Equal(s.T(), "CLUSTER.ADMINISTRATOR", privilege.PrivilegeInfo.PermissionName)
	assert.Equal(s.T(), "hm_etl_outils", privilege.PrivilegeInfo.PrincipalName)
	assert.Equal(s.T(), "GROUP", privilege.PrivilegeInfo.PrincipalType)

	// Get privilege
	privilege, err = s.client.Privilege("sihm-test", privilege.PrivilegeInfo.PrivilegeId)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), privilege)
	if privilege != nil {
		assert.Equal(s.T(), "CLUSTER.ADMINISTRATOR", privilege.PrivilegeInfo.PermissionName)
		assert.Equal(s.T(), "hm_etl_outils", privilege.PrivilegeInfo.PrincipalName)
		assert.Equal(s.T(), "GROUP", privilege.PrivilegeInfo.PrincipalType)
		assert.NotEqual(s.T(), "", privilege.PrivilegeInfo.PrivilegeId)
	}

	// Update privilege
	privilege.PrivilegeInfo.PermissionName = "CLUSTER.OPERATOR"
	privilege, err = s.client.UpdatePrivilege("sihm-test", privilege)
	assert.NoError(s.T(), err)

	// Delete privilege
	err = s.client.DeletePrivilege("sihm-test", privilege.PrivilegeInfo.PrivilegeId)
	assert.NoError(s.T(), err)

}

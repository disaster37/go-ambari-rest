package client

import (
	"github.com/stretchr/testify/assert"
)

func (s *ClientTestSuite) TestPrivilege() {

	// Create privilege
	privilege := &Privilege{
		PrivilegeInfo: &PrivilegeInfo{
			PermissionName: "CLUSTER.ADMINISTRATOR",
			PrincipalName:  "admin",
			PrincipalType:  "USER",
		},
	}
	privilege, err := s.client.CreatePrivilege("test", privilege)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), privilege)
	if privilege != nil {
		assert.NotEqual(s.T(), "", privilege.PrivilegeInfo.PrivilegeId)
		assert.Equal(s.T(), "CLUSTER.ADMINISTRATOR", privilege.PrivilegeInfo.PermissionName)
		assert.Equal(s.T(), "admin", privilege.PrivilegeInfo.PrincipalName)
		assert.Equal(s.T(), "USER", privilege.PrivilegeInfo.PrincipalType)
	}

	// Get privilege
	privilege, err = s.client.Privilege("test", privilege.PrivilegeInfo.PrivilegeId)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), privilege)
	if privilege != nil {
		assert.Equal(s.T(), "CLUSTER.ADMINISTRATOR", privilege.PrivilegeInfo.PermissionName)
		assert.Equal(s.T(), "admin", privilege.PrivilegeInfo.PrincipalName)
		assert.Equal(s.T(), "USER", privilege.PrivilegeInfo.PrincipalType)
		assert.NotEqual(s.T(), "", privilege.PrivilegeInfo.PrivilegeId)
	}

	// Update privilege
	privilege.PrivilegeInfo.PermissionName = "CLUSTER.OPERATOR"
	privilege, err = s.client.UpdatePrivilege("test", privilege)
	assert.NoError(s.T(), err)

	// Delete privilege
	err = s.client.DeletePrivilege("test", privilege.PrivilegeInfo.PrivilegeId)
	assert.NoError(s.T(), err)

}

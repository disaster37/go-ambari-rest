package client

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

// Test the constructor
func (s *ClientTestSuite) TestBlueprint() {

	// Create blueprint
	b, err := ioutil.ReadFile("../fixtures/blueprint.json")
	if err != nil {
		panic(err)
	}
	blueprintJson := string(b)

	blueprint, err := s.client.CreateBlueprint("test", blueprintJson)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), blueprint)
	if blueprint != nil {
		assert.Equal(s.T(), "test", blueprint.BlueprintInfo.Name)
		assert.Equal(s.T(), "HDP", blueprint.BlueprintInfo.Stack)
		assert.Equal(s.T(), "2.6", blueprint.BlueprintInfo.Version)
		assert.Equal(s.T(), 0, len(blueprint.Configurations))
		assert.Equal(s.T(), 2, len(blueprint.HostGroups))
	}

	// Get blueprint
	blueprint, err = s.client.Blueprint("test")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), blueprint)
	if blueprint != nil {
		assert.Equal(s.T(), "test", blueprint.BlueprintInfo.Name)
		assert.Equal(s.T(), "HDP", blueprint.BlueprintInfo.Stack)
		assert.Equal(s.T(), "2.6", blueprint.BlueprintInfo.Version)
		//		assert.NotEqual(s.T(), 0, len(blueprint.Configurations))
		//		assert.Equal(s.T(), 1, len(blueprint.HostGroups))
	}

	assert.NoError(s.T(), err)

	// Delete blueprint
	err = s.client.DeleteBlueprint("test")
	assert.NoError(s.T(), err)

}

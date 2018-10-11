package client

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

func (s *ClientTestSuite) TestBlueprint() {

	// Create blueprint
	b, err := ioutil.ReadFile("../fixtures/blueprint.json")
	if err != nil {
		panic(err)
	}
	blueprintJson := string(b)

	blueprint, err := s.client.CreateBlueprint("testBlueprint", blueprintJson)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), blueprint)
	if blueprint != nil {
		assert.Equal(s.T(), "testBlueprint", blueprint.BlueprintInfo.Name)
		assert.Equal(s.T(), "HDP", blueprint.BlueprintInfo.Stack)
		assert.Equal(s.T(), "2.6", blueprint.BlueprintInfo.Version)
		assert.Equal(s.T(), 0, len(blueprint.Configurations))
		assert.Equal(s.T(), 2, len(blueprint.HostGroups))
	}

	// Get blueprint
	blueprint, err = s.client.Blueprint("testBlueprint")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), blueprint)
	if blueprint != nil {
		assert.Equal(s.T(), "testBlueprint", blueprint.BlueprintInfo.Name)
		assert.Equal(s.T(), "HDP", blueprint.BlueprintInfo.Stack)
		assert.Equal(s.T(), "2.6", blueprint.BlueprintInfo.Version)
	}

	assert.NoError(s.T(), err)

	// Delete blueprint
	err = s.client.DeleteBlueprint("testBlueprint")
	assert.NoError(s.T(), err)
	blueprint, err = s.client.Blueprint("testBlueprint")
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), blueprint)
}

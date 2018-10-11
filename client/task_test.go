package client

import (
	"github.com/stretchr/testify/assert"
)

func (s *ClientTestSuite) TestTask() {

	// Get task
	requestTask, err := s.client.Request("test", 4)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), requestTask)
	if requestTask != nil {
		assert.Equal(s.T(), 0, requestTask.RequestTaskInfo.AbordedTask)
		assert.Equal(s.T(), 0, requestTask.RequestTaskInfo.FailedTask)
		assert.Equal(s.T(), 2, requestTask.RequestTaskInfo.TaskCount)
		assert.Equal(s.T(), 2, requestTask.RequestTaskInfo.CompletedTask)
		assert.Equal(s.T(), float64(100), requestTask.RequestTaskInfo.ProgressPercent)
		assert.Equal(s.T(), REQUEST_COMPLETED, requestTask.RequestTaskInfo.Status)
	}
}

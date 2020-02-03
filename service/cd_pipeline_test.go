package service

import (
	"edp-admin-console/models"
	"edp-admin-console/service/cd_pipeline"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type CDPipelineRepositoryMock struct {
	mock.Mock
}

func (mock *CDPipelineRepositoryMock) GetCDPipelineByName(pipelineName string) (*models.CDPipelineDTO, error) {
	args := mock.Called(pipelineName)
	return args.Get(0).(*models.CDPipelineDTO), args.Error(1)
}

func TestShouldReturnCDPipeline(t *testing.T) {
	repositoryMock := new(CDPipelineRepositoryMock)
	repositoryMock.On("GetCDPipelineByName", mock.Anything).Return(&models.CDPipelineDTO{}, nil)
	pipelineService := cd_pipeline.CDPipelineService{ICDPipelineRepository: repositoryMock}

	cdPipeline, err := pipelineService.GetCDPipelineByName(mock.Anything)
	assert.Nil(t, err)
	assert.NotNil(t, cdPipeline)
}

func TestShouldReturnErrorCDPipeline(t *testing.T) {
	repositoryMock := new(CDPipelineRepositoryMock)
	repositoryMock.On("GetCDPipelineByName", mock.Anything).Return(&models.CDPipelineDTO{}, errors.New("internal error"))
	pipelineService := cd_pipeline.CDPipelineService{ICDPipelineRepository: repositoryMock}

	cdPipeline, err := pipelineService.GetCDPipelineByName(mock.Anything)
	assert.NotNil(t, err)
	assert.Nil(t, cdPipeline)
}

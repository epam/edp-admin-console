package service

import (
	"edp-admin-console/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type ReleaseBranchRepositoryMock struct {
	mock.Mock
}

func (mock *ReleaseBranchRepositoryMock) GetAllReleaseBranches(appName, edpName string) ([]models.ReleaseBranchView, error) {
	args := mock.Called(edpName)
	return args.Get(0).([]models.ReleaseBranchView), args.Error(1)
}

func (mock *ReleaseBranchRepositoryMock) GetReleaseBranch(appName, branchName, edpName string) (*models.ReleaseBranchView, error) {
	args := mock.Called(appName, edpName)
	return args.Get(0).(*models.ReleaseBranchView), args.Error(1)
}

func TestShouldReturnBranch(t *testing.T) {
	repositoryMock := new(ReleaseBranchRepositoryMock)
	repositoryMock.On("GetReleaseBranch", mock.Anything, mock.Anything, mock.Anything).Return(&models.ReleaseBranchView{}, nil)
	branchService := BranchService{IReleaseBranchRepository: repositoryMock}

	branchEntity, err := branchService.GetReleaseBranch(mock.Anything, mock.Anything)
	assert.Nil(t, err)
	assert.NotNil(t, branchEntity)
}

func TestShouldReturnErrorBranch(t *testing.T) {
	repositoryMock := new(ReleaseBranchRepositoryMock)
	repositoryMock.On("GetReleaseBranch", mock.Anything, mock.Anything, mock.Anything).Return(&models.ReleaseBranchView{}, errors.New("internal error"))
	branchService := BranchService{IReleaseBranchRepository: repositoryMock}

	branchEntity, err := branchService.GetReleaseBranch(mock.Anything, mock.Anything)
	assert.NotNil(t, err)
	assert.Nil(t, branchEntity)
}

func TestShouldReturnAllBranch(t *testing.T) {
	repositoryMock := new(ReleaseBranchRepositoryMock)
	repositoryMock.On("GetAllReleaseBranches", mock.Anything, mock.Anything).Return(createBranchEntities(), nil)
	branchService := BranchService{IReleaseBranchRepository: repositoryMock}

	branchEntity, err := branchService.GetAllReleaseBranches(mock.Anything)
	assert.Nil(t, err)
	assert.NotNil(t, branchEntity)
}

func TestShouldReturnErrorAllBranch(t *testing.T) {
	repositoryMock := new(ReleaseBranchRepositoryMock)
	repositoryMock.On("GetAllReleaseBranches", mock.Anything, mock.Anything).Return(createBranchEntities(), errors.New("internal error"))
	branchService := BranchService{IReleaseBranchRepository: repositoryMock}

	branchEntity, err := branchService.GetAllReleaseBranches(mock.Anything)
	assert.NotNil(t, err)
	assert.Nil(t, branchEntity)
}

func createBranchEntities() []models.ReleaseBranchView {
	return []models.ReleaseBranchView{
		{
			Name:            "fake-name",
			Event:           "fake-event",
			Username:        "fake-username",
			DetailedMessage: "fake-message",
			UpdatedAt:       time.Now(),
		},
	}
}

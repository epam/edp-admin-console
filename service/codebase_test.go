package service

import (
	"edp-admin-console/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

//todo add test for CreateCodebase and GetCodebaseCR
type ApplicationEntityRepositoryMock struct {
	mock.Mock
}

func (mock *ApplicationEntityRepositoryMock) GetAllApplications(filterCriteria models.CodebaseCriteria) ([]models.Application, error) {
	args := mock.Called(filterCriteria)
	return args.Get(0).([]models.Application), args.Error(1)
}

func (mock *ApplicationEntityRepositoryMock) GetApplication(appName string) (*models.ApplicationInfo, error) {
	args := mock.Called(appName)
	return args.Get(0).(*models.ApplicationInfo), args.Error(1)
}

func (mock *ApplicationEntityRepositoryMock) GetAllApplicationsWithReleaseBranches(applicationFilterCriteria models.CodebaseCriteria) ([]models.ApplicationWithReleaseBranch, error) {
	args := mock.Called(applicationFilterCriteria)
	return args.Get(0).([]models.ApplicationWithReleaseBranch), args.Error(1)
}

func TestShouldReturnApplications(t *testing.T) {
	repositoryMock := new(ApplicationEntityRepositoryMock)
	repositoryMock.On("GetCodebasesByCriteria", mock.Anything).Return(createApplications(), nil)
	appService := codebase.CodebaseService{ICodebaseRepository: repositoryMock}

	applications, err := appService.GetCodebasesByCriteria(models.CodebaseCriteria{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(applications), "they should be equal")
}

func TestShouldReturnErrorFromGetAllApplications(t *testing.T) {
	repositoryMock := new(ApplicationEntityRepositoryMock)
	repositoryMock.On("GetCodebasesByCriteria", mock.Anything).Return(createApplications(), errors.New("internal error"))
	appService := codebase.CodebaseService{ICodebaseRepository: repositoryMock}

	applications, err := appService.GetCodebasesByCriteria(models.CodebaseCriteria{})
	assert.Error(t, err)
	assert.Nil(t, applications)
}

func TestShouldReturnApplication(t *testing.T) {
	repositoryMock := new(ApplicationEntityRepositoryMock)
	repositoryMock.On("GetCodebaseByName", mock.Anything, mock.Anything).Return(&models.ApplicationInfo{}, nil)
	appService := codebase.CodebaseService{ICodebaseRepository: repositoryMock}

	application, err := appService.GetCodebaseByName(mock.Anything)
	assert.Nil(t, err)
	assert.NotNil(t, application)
}

func TestShouldReturnErrorFromGetAllApplication(t *testing.T) {
	repositoryMock := new(ApplicationEntityRepositoryMock)
	repositoryMock.On("GetCodebaseByName", mock.Anything, mock.Anything).Return(&models.ApplicationInfo{}, errors.New("internal error"))
	appService := codebase.CodebaseService{ICodebaseRepository: repositoryMock}

	application, err := appService.GetCodebaseByName(mock.Anything)
	assert.NotNil(t, err)
	assert.Nil(t, application)
}

func TestShouldReturnApplicationsWithReleaseBranches(t *testing.T) {
	repositoryMock := new(ApplicationEntityRepositoryMock)
	repositoryMock.On("GetAllCodebasesWithReleaseBranches", mock.Anything).Return([]models.ApplicationWithReleaseBranch{}, errors.New("internal error"))
	appService := codebase.CodebaseService{ICodebaseRepository: repositoryMock}

	applications, err := appService.GetAllCodebasesWithReleaseBranches(models.CodebaseCriteria{})
	assert.NotNil(t, err)
	assert.Nil(t, applications)
}

func TestShouldReturnErrorApplicationsWithReleaseBranches(t *testing.T) {
	repositoryMock := new(ApplicationEntityRepositoryMock)
	repositoryMock.On("GetAllCodebasesWithReleaseBranches", mock.Anything).Return([]models.ApplicationWithReleaseBranch{}, nil)
	appService := codebase.CodebaseService{ICodebaseRepository: repositoryMock}

	applications, err := appService.GetAllCodebasesWithReleaseBranches(models.CodebaseCriteria{})
	assert.Nil(t, err)
	assert.NotNil(t, applications)
}

func createApplications() []models.Application {
	return []models.Application{
		{
			Name:      "testName",
			Language:  "testLang",
			BuildTool: "testBuildTool",
		},
	}
}

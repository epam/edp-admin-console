package service

import (
	"edp-admin-console/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

//todo add test for CreateApp and GetApplicationCR
type ApplicationEntityRepositoryMock struct {
	mock.Mock
}

func (mock *ApplicationEntityRepositoryMock) GetAllApplications(edpName string) ([]models.Application, error) {
	args := mock.Called(edpName)
	return args.Get(0).([]models.Application), args.Error(1)
}

func (mock *ApplicationEntityRepositoryMock) GetApplication(appName string, edpName string) (*models.ApplicationInfo, error) {
	args := mock.Called(appName, edpName)
	return args.Get(0).(*models.ApplicationInfo), args.Error(1)
}

func (mock *ApplicationEntityRepositoryMock) GetAllApplicationsWithReleaseBranches() ([]models.ApplicationWithReleaseBranch, error) {
	args := mock.Called()
	return args.Get(0).([]models.ApplicationWithReleaseBranch), args.Error(1)
}

func TestShouldReturnApplications(t *testing.T) {
	repositoryMock := new(ApplicationEntityRepositoryMock)
	repositoryMock.On("GetAllApplications", mock.Anything).Return(createApplications(), nil)
	appService := ApplicationService{IApplicationRepository: repositoryMock}

	applications, err := appService.GetAllApplications()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(applications), "they should be equal")
}

func TestShouldReturnErrorFromGetAllApplications(t *testing.T) {
	repositoryMock := new(ApplicationEntityRepositoryMock)
	repositoryMock.On("GetAllApplications", mock.Anything).Return(createApplications(), errors.New("internal error"))
	appService := ApplicationService{IApplicationRepository: repositoryMock}

	applications, err := appService.GetAllApplications()
	assert.Error(t, err)
	assert.Nil(t, applications)
}

func TestShouldReturnApplication(t *testing.T) {
	repositoryMock := new(ApplicationEntityRepositoryMock)
	repositoryMock.On("GetApplication", mock.Anything, mock.Anything).Return(&models.ApplicationInfo{}, nil)
	appService := ApplicationService{IApplicationRepository: repositoryMock}

	application, err := appService.GetApplication(mock.Anything)
	assert.Nil(t, err)
	assert.NotNil(t, application)
}

func TestShouldReturnErrorFromGetAllApplication(t *testing.T) {
	repositoryMock := new(ApplicationEntityRepositoryMock)
	repositoryMock.On("GetApplication", mock.Anything, mock.Anything).Return(&models.ApplicationInfo{}, errors.New("internal error"))
	appService := ApplicationService{IApplicationRepository: repositoryMock}

	application, err := appService.GetApplication(mock.Anything)
	assert.NotNil(t, err)
	assert.Nil(t, application)
}

func TestShouldReturnApplicationsWithReleaseBranches(t *testing.T) {
	repositoryMock := new(ApplicationEntityRepositoryMock)
	repositoryMock.On("GetAllApplicationsWithReleaseBranches", mock.Anything).Return([]models.ApplicationWithReleaseBranch{}, errors.New("internal error"))
	appService := ApplicationService{IApplicationRepository: repositoryMock}

	applications, err := appService.GetAllApplicationsWithReleaseBranches()
	assert.NotNil(t, err)
	assert.Nil(t, applications)
}

func TestShouldReturnErrorApplicationsWithReleaseBranches(t *testing.T) {
	repositoryMock := new(ApplicationEntityRepositoryMock)
	repositoryMock.On("GetAllApplicationsWithReleaseBranches", mock.Anything).Return([]models.ApplicationWithReleaseBranch{}, nil)
	appService := ApplicationService{IApplicationRepository: repositoryMock}

	applications, err := appService.GetAllApplicationsWithReleaseBranches()
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

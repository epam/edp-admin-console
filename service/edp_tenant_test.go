package service

import (
	"edp-admin-console/models"
	"errors"
	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func init() {
	beego.AppConfig.Set("adminRole", "administrator")
	beego.AppConfig.Set("developerRole", "developer")
}

type EDPTenantRepositoryMock struct {
	mock.Mock
}

func (mock *EDPTenantRepositoryMock) GetEdpVersionByName(edpName string) (string, error) {
	args := mock.Called(edpName)
	return args.String(0), args.Error(1)
}

func (mock *EDPTenantRepositoryMock) GetTenantByName(edpName string) (*models.EDPTenant, error) {
	args := mock.Called(edpName)
	return args.Get(0).(*models.EDPTenant), args.Error(1)
}

func TestShouldReturnEDPVersionByName(t *testing.T) {
	repositoryMock := new(EDPTenantRepositoryMock)
	repositoryMock.On("GetEdpVersionByName", mock.Anything).Return(mock.Anything, nil)
	edpService := EDPTenantService{IEDPTenantRep: repositoryMock}

	tenant, err := edpService.GetEDPVersion()
	assert.NoError(t, err)
	assert.NotEmpty(t, tenant)
}

func TestShouldReturnErrorFromEDPVersionByName(t *testing.T) {
	repositoryMock := new(EDPTenantRepositoryMock)
	repositoryMock.On("GetEdpVersionByName", mock.Anything).Return(mock.Anything, errors.New("internal error"))
	edpService := EDPTenantService{IEDPTenantRep: repositoryMock}

	tenant, err := edpService.GetEDPVersion()
	assert.Error(t, err)
	assert.Empty(t, tenant)
}

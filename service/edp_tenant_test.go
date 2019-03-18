package service

import (
	"edp-admin-console/models"
	"errors"
	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func init()  {
	beego.AppConfig.Set("adminRole", "administrator")
	beego.AppConfig.Set("developerRole", "developer")
}

type EDPTenantRepositoryMock struct {
	mock.Mock
}

func (mock *EDPTenantRepositoryMock) GetAllEDPTenantsByNames(adminClients []string) ([]*models.EDPTenant, error) {
	args := mock.Called(adminClients)
	return args.Get(0).([]*models.EDPTenant), args.Error(1)
}

func (mock *EDPTenantRepositoryMock) GetEdpVersionByName(edpName string) (string, error) {
	args := mock.Called(edpName)
	return args.String(0), args.Error(1)
}

func (mock *EDPTenantRepositoryMock) GetTenantByName(edpName string) (*models.EDPTenant, error) {
	args := mock.Called(edpName)
	return args.Get(0).(*models.EDPTenant), args.Error(1)
}

func TestShouldReturnTenantByName(t *testing.T) {
	repositoryMock := new(EDPTenantRepositoryMock)
	repositoryMock.On("GetTenantByName", mock.Anything).Return(&models.EDPTenant{}, nil)
	edpService := EDPTenantService{IEDPTenantRep: repositoryMock}

	tenant, err := edpService.GetTenantByName(mock.Anything)
	assert.NoError(t, err)
	assert.NotNil(t, tenant)
}

func TestShouldReturnErrorFromTenantByName(t *testing.T) {
	repositoryMock := new(EDPTenantRepositoryMock)
	repositoryMock.On("GetTenantByName", mock.Anything).Return(&models.EDPTenant{}, errors.New("internal error"))
	edpService := EDPTenantService{IEDPTenantRep: repositoryMock}

	tenant, err := edpService.GetTenantByName(mock.Anything)
	assert.Error(t, err)
	assert.Nil(t, tenant)
}

func TestShouldReturnEDPVersionByName(t *testing.T) {
	repositoryMock := new(EDPTenantRepositoryMock)
	repositoryMock.On("GetEdpVersionByName", mock.Anything).Return(mock.Anything, nil)
	edpService := EDPTenantService{IEDPTenantRep: repositoryMock}

	tenant, err := edpService.GetEDPVersionByName(mock.Anything)
	assert.NoError(t, err)
	assert.NotEmpty(t, tenant)
}

func TestShouldReturnErrorFromEDPVersionByName(t *testing.T) {
	repositoryMock := new(EDPTenantRepositoryMock)
	repositoryMock.On("GetEdpVersionByName", mock.Anything).Return(mock.Anything, errors.New("internal error"))
	edpService := EDPTenantService{IEDPTenantRep: repositoryMock}

	tenant, err := edpService.GetEDPVersionByName(mock.Anything)
	assert.Error(t, err)
	assert.Empty(t, tenant)
}

func TestShouldReturnErrorFromEDPTenantsByNameBecauseInvalidArguments(t *testing.T) {
	repositoryMock := new(EDPTenantRepositoryMock)
	repositoryMock.On("GetAllEDPTenantsByNames", mock.Anything).Return(createEDPTenants(), nil)
	edpService := EDPTenantService{IEDPTenantRep: repositoryMock}

	tenants, err := edpService.GetEDPTenants(map[string][]string{
		"client1": {"developer", "administrator"},
		"client2": {"developer"},
	})
	assert.NoError(t, err)
	assert.Nil(t, tenants, "should be nil as there're all incorrect tenants")
}

func TestShouldReturnAllEDPTenants(t *testing.T) {
	repositoryMock := new(EDPTenantRepositoryMock)
	repositoryMock.On("GetAllEDPTenantsByNames", mock.Anything).Return(createEDPTenants(), nil)
	edpService := EDPTenantService{IEDPTenantRep: repositoryMock}

	tenants, err := edpService.GetEDPTenants(map[string][]string{
		"client1-edp": {"developer", "administrator"},
		"client2-edp": {"developer"},
	})

	assert.NoError(t, err)
	assert.Equal(t,1, len(tenants))
}

func TestShouldReturnErrorFromGetAllEDPTenantsByNames(t *testing.T) {
	repositoryMock := new(EDPTenantRepositoryMock)
	repositoryMock.On("GetAllEDPTenantsByNames", mock.Anything).Return(createEDPTenants(), errors.New("internal error"))
	edpService := EDPTenantService{IEDPTenantRep: repositoryMock}

	tenants, err := edpService.GetEDPTenants(map[string][]string{
		"client1-edp": {"developer", "administrator"},
		"client2-edp": {"developer"},
	})

	assert.Error(t, err)
	assert.NotEqual(t,1, len(tenants))
}

func createEDPTenants() []*models.EDPTenant {
	return []*models.EDPTenant{
		{
			Name:    "testName",
			Version: "v1.0",
		},
	}
}

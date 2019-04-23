package service

import (
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

func TestShouldReturnEDPVersionByName(t *testing.T) {
	repositoryMock := new(EDPTenantRepositoryMock)
	repositoryMock.On("GetEdpVersionByName", mock.Anything).Return(mock.Anything, nil)
	edpService := EDPTenantService{}

	tenant, err := edpService.GetEDPVersion()
	assert.NoError(t, err)
	assert.NotEmpty(t, tenant)
}

func TestShouldReturnErrorFromEDPVersionByName(t *testing.T) {
	repositoryMock := new(EDPTenantRepositoryMock)
	repositoryMock.On("GetEdpVersionByName", mock.Anything).Return(mock.Anything, errors.New("internal error"))
	edpService := EDPTenantService{}

	tenant, err := edpService.GetEDPVersion()
	assert.Error(t, err)
	assert.Empty(t, tenant)
}

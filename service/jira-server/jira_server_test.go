package jira_server

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetJiraServersMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	mjs := new(mock.MockJiraServer)
	js := JiraServer{
		IJiraServer: mjs,
	}

	mjs.On("GetJiraServers").Return(
		[]*query.JiraServer{
			{
				Id:   1,
				Name: "fake-server",
			},
		}, nil)

	jsName, err := js.GetJiraServers()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jsName))
}

func TestGetJiraServersMethod_ShouldReturnError(t *testing.T) {
	mjs := new(mock.MockJiraServer)
	js := JiraServer{
		IJiraServer: mjs,
	}

	mjs.On("GetJiraServers").Return(nil, errors.New("failed"))

	_, err := js.GetJiraServers()
	assert.Error(t, err)
}

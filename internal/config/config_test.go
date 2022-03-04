package config

import (
	"context"
	"testing"

	"edp-admin-console/util/consts"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

func TestSetupConfig(t *testing.T) {
	mustSetAppConfig(t, AuthEnable, "false")
	mustSetAppConfig(t, DiagramPageEnabled, "false")
	mustSetAppConfig(t, vcsIntegrationEnabled, "false")
	mustSetAppConfig(t, platformType, consts.Openshift)
	mustSetAppConfig(t, IntegrationStrategies, IntegrationStrategies)
	mustSetAppConfig(t, BuildTools, BuildTools)
	mustSetAppConfig(t, VersioningTypes, VersioningTypes)
	mustSetAppConfig(t, DeploymentScript, DeploymentScript)
	mustSetAppConfig(t, CiTools, CiTools)
	mustSetAppConfig(t, PerfDataSources, PerfDataSources)

	ctx := context.Background()
	emptyStr := ""
	config, err := SetupConfig(ctx, emptyStr)
	assert.NoError(t, err)
	assert.False(t, config.AuthEnable)
}

func mustSetAppConfig(t *testing.T, k, v string) {
	t.Helper()
	err := beego.AppConfig.Set(k, v)
	if err != nil {
		t.Fatal(err)
	}
}

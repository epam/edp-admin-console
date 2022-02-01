package webapi

import (
	"testing"

	"github.com/astaxie/beego"
)

func TestSetupRouter(t *testing.T) {
	mustSetAppConfig(t, "runmode", "local")
	mustSetAppConfig(t, "keycloakAuthEnabled", "false")
	mustSetAppConfig(t, "dbEnabled", "false")
	mustSetAppConfig(t, "integrationStrategies", "Create,Clone")
	mustSetAppConfig(t, "buildTools", "maven")
	mustSetAppConfig(t, "versioningTypes", "default,edp")
	mustSetAppConfig(t, "testReportTools", "allure")
	mustSetAppConfig(t, "deploymentScript", "openshift-template")
	mustSetAppConfig(t, "ciTools", "Jenkins,GitLab CI")
	mustSetAppConfig(t, "perfDataSources", "Sonar,Jenkins,GitLab")
	SetupRouter() // rewrite this setup ASAP
}

func mustSetAppConfig(t *testing.T, k, v string) {
	t.Helper()
	err := beego.AppConfig.Set(k, v)
	if err != nil {
		t.Fatal(err)
	}
}

package webapi

import (
	"os"
	"testing"

	"github.com/astaxie/beego"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
)

func TestSetupRouter(t *testing.T) {
	err := os.Setenv(k8s.NamespaceEnv, "Test")
	if err != nil {
		t.Fatal(err)
	}
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
	fakeClient := fake.NewClientBuilder().Build()
	namespaceClient, err := k8s.NewRuntimeNamespacedClient(fakeClient, "test")
	if err != nil {
		t.Fatal(err)
	}
	conf := &config.AppConfig{AuthEnable: false}
	clusterConf := &rest.Config{}
	SetupRouter(namespaceClient, "", conf, clusterConf) // rewrite this setup ASAP
	err = os.Unsetenv(k8s.NamespaceEnv)
	if err != nil {
		t.Fatal(err)
	}
}

func mustSetAppConfig(t *testing.T, k, v string) {
	t.Helper()
	err := beego.AppConfig.Set(k, v)
	if err != nil {
		t.Fatal(err)
	}
}

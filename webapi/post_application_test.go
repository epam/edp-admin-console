package webapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

type CreateApplicationSuite struct {
	suite.Suite
	TestServer *httptest.Server
	Handler    *HandlerEnv
}

func TestCreateApplicationSuite(t *testing.T) {
	scheme := runtime.NewScheme()
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	workingDir, _ := path.Split(currentDir)
	conf := &config.AppConfig{
		BasePath:   "/",
		AuthEnable: false,
	}
	h := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMap(CreateCommonFuncMap()), WithConfig(conf))
	logger := applog.GetLogger()
	authHandler := HandlerAuthWithOption()
	router := V2APIRouter(h, authHandler, logger)

	s := &CreateApplicationSuite{
		Handler: h,
	}
	s.TestServer = httptest.NewServer(router)
	suite.Run(t, s)
}

func (s *CreateApplicationSuite) TestCreateApplication() {
	type JiraIssueMetadataForm struct {
		JiraFieldName []string `form:"jiraFieldName"`
		JiraPattern   []string `form:"jiraPattern"`
	}
	jiraServer := "edp-jira"
	appLang := "go"
	appName := "testName"
	buildTool := "Go"
	defaultBranchName := "master"
	deploymentScript := "helm-chart"
	ciTool := "Jenkins"
	framework := "operator-sdk"
	jenkinsSlave := "go"
	jobProvisioning := "default"

	fields := []string{"field1", "field2"}
	patterns := []string{"pattern1", "pattern2"}

	t := s.T()
	httpExpect := httpexpect.New(t, s.TestServer.URL)
	req := httpExpect.POST("/v2/admin/edp/application").
		WithFormField("jiraServer", jiraServer).
		WithForm(JiraIssueMetadataForm{
			JiraFieldName: fields,
			JiraPattern:   patterns}).
		WithFormField("appLang", appLang).
		WithFormField("appName", appName).
		WithFormField("buildTool", buildTool).
		WithFormField("defaultBranchName", defaultBranchName).
		WithFormField("deploymentScript", deploymentScript).
		WithFormField("ciTool", ciTool).
		WithFormField("framework", framework).
		WithFormField("jenkinsSlave", jenkinsSlave).
		WithFormField("jobProvisioning", jobProvisioning).
		WithRedirectPolicy(httpexpect.DontFollowRedirects)

	req.Expect().
		Status(http.StatusFound).
		Header("location").
		Equal(fmt.Sprintf("%s/v2/admin/edp/application/overview?%s=%s#codebaseSuccessModal", s.Handler.Config.BasePath, "waitingforcodebase", appName))
	codebase, err := s.Handler.NamespacedClient.GetCodebase(context.Background(), appName)
	assert.NoError(t, err)

	expectedJiraIssueMetadata := fmt.Sprintf(`{"%s":"%s","%s":"%s"}`, fields[0], patterns[0], fields[1], patterns[1])

	assert.NoError(t, err)
	assert.Equal(t, appName, codebase.Name)
	assert.Equal(t, jiraServer, *codebase.Spec.JiraServer)
	assert.Equal(t, appLang, codebase.Spec.Lang)
	assert.Equal(t, buildTool, codebase.Spec.BuildTool)
	assert.Equal(t, defaultBranchName, codebase.Spec.DefaultBranch)
	assert.Equal(t, deploymentScript, codebase.Spec.DeploymentScript)
	assert.Equal(t, ciTool, codebase.Spec.CiTool)
	assert.Equal(t, framework, *codebase.Spec.Framework)
	assert.Equal(t, jenkinsSlave, *codebase.Spec.JenkinsSlave)
	assert.Equal(t, jobProvisioning, *codebase.Spec.JobProvisioning)
	assert.Equal(t, expectedJiraIssueMetadata, *codebase.Spec.JiraIssueMetadataPayload)
}

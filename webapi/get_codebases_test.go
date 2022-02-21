package webapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/gavv/httpexpect/v2"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

const (
	testNamespace = "cd_namespace"
)

type GetCodebasesSuite struct {
	suite.Suite
	Router     *chi.Mux
	TestServer *httptest.Server
}

func TestGetCodebasesSuite(t *testing.T) {
	scheme := runtime.NewScheme()
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}

	builder := fake.NewClientBuilder().WithScheme(scheme)
	client := builder.Build()
	k8sClient, err := k8s.NewRuntimeNamespacedClient(client, testNamespace)
	if err != nil {
		t.Fatal(err)
	}
	conf := &config.AppConfig{
		BasePath:   "/",
		AuthEnable: false,
		EDPVersion: "v1",
	}
	h := NewHandlerEnv(WithClient(k8sClient), WithConfig(conf))
	authHandler := HandlerAuthWithOption()
	logger := applog.GetLogger()

	router := V2APIRouter(h, authHandler, logger)

	s := &GetCodebasesSuite{
		Router: router,
	}
	suite.Run(t, s)
}

func (s *GetCodebasesSuite) SetupSuite() {
	s.TestServer = httptest.NewServer(s.Router)
}

func (s *GetCodebasesSuite) TearDownSuite() {
	s.TestServer.Close()
}

func (s *GetCodebasesSuite) RedefineK8SClientWithCodebaseCR(crCodebases []*codeBaseApi.Codebase) {
	runtimeScheme := runtime.NewScheme()
	runtimeScheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.Codebase{})
	namespaceName := testNamespace

	builder := fake.NewClientBuilder().WithScheme(runtimeScheme)
	if len(crCodebases) > 0 {
		fakeObjects := make([]runtime.Object, 0)
		for _, crCodebase := range crCodebases {
			fakeObjects = append(fakeObjects, crCodebase)
		}
		builder = builder.WithRuntimeObjects(fakeObjects...)
	}

	fakeClient := builder.Build()
	k8sClient, err := k8s.NewRuntimeNamespacedClient(fakeClient, namespaceName)
	if err != nil {
		s.T().Fatal(err)
	}
	conf := &config.AppConfig{
		BasePath:   "/",
		AuthEnable: false,
	}
	testHandler := NewHandlerEnv(WithClient(k8sClient), WithConfig(conf))
	authHandler := HandlerAuthWithOption()

	newRouter := V2APIRouter(testHandler, authHandler, applog.GetLogger())
	s.TestServer.Config.Handler = newRouter
}

func WithGITUrl(gitURL string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.GitUrlPath = &gitURL
	}
}

func (s *GetCodebasesSuite) TestGetCodebases_OK() {
	t := s.T()

	namespaceName := testNamespace
	crCodebaseName_1 := "fake_spring-petclinic"
	crCodebaseGitServer_1 := "gerrit"
	crVersioningType_1 := "edp"
	crStrategy_1 := "clone"
	crGitProjectPath_1 := "git.com/foo/bar"
	deploymentScript_1 := "bash"
	crCodebaseName_2 := "fake_2"
	crCodebaseGitServer_2 := "gitlab"
	crVersioningType_2 := "default"
	crStrategy_2 := "create"
	crGitProjectPath_2 := "hg.com/foo/bar"
	deploymentScript_2 := "hands"

	stubCodebase_1 := createCodebaseCRWithOptions(
		WithCrNamespace(namespaceName),
		WithCrName(crCodebaseName_1),
		WithGitServerName(crCodebaseGitServer_1),
		WithVersioningType(crVersioningType_1),
		WithStrategy(crStrategy_1),
		WithGITUrl(crGitProjectPath_1),
		WithDeploymentScript(deploymentScript_1),
	)
	stubCodebase_2 := createCodebaseCRWithOptions(
		WithCrNamespace(namespaceName),
		WithCrName(crCodebaseName_2),
		WithGitServerName(crCodebaseGitServer_2),
		WithVersioningType(crVersioningType_2),
		WithStrategy(crStrategy_2),
		WithGITUrl(crGitProjectPath_2),
		WithDeploymentScript(deploymentScript_2),
	)
	stubCodebases := []*codeBaseApi.Codebase{
		stubCodebase_1, stubCodebase_2,
	}
	s.RedefineK8SClientWithCodebaseCR(stubCodebases)
	httpExpect := httpexpect.New(t, s.TestServer.URL)
	response := httpExpect.
		GET("/api/v2/edp/codebase").
		WithQuery("codebases", "fake_spring-petclinic,fake_2").
		Expect().
		Status(http.StatusOK).
		ContentType("application/json")

	expectedJSONBody := fmt.Sprintf(`[
{
    "deploymentScript": "%s",
    "emptyProject": false,
    "gitServer": "%s",
    "gitProjectPath": "%s",
    "jenkinsSlave": "",
    "name": "%s",
    "strategy": "%s",
    "type": "",
    "versioningType": "%s"
},
{
    "deploymentScript": "%s",
    "emptyProject": false,
    "gitServer": "%s",
    "gitProjectPath": "%s",
    "jenkinsSlave": "",
    "name": "%s",
    "strategy": "%s",
    "type": "",
    "versioningType": "%s"
}
]`,
		deploymentScript_1, crCodebaseGitServer_1, crGitProjectPath_1, crCodebaseName_1, crStrategy_1, crVersioningType_1,
		deploymentScript_2, crCodebaseGitServer_2, crGitProjectPath_2, crCodebaseName_2, crStrategy_2, crVersioningType_2,
	)

	gotPlainBody := response.Body().Raw()
	assert.JSONEq(t, expectedJSONBody, gotPlainBody, "unexpected body")
}

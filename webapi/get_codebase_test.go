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

	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

type GetCodebaseSuite struct {
	suite.Suite
	Router     *chi.Mux
	TestServer *httptest.Server
}

func TestGetCodebaseSuite(t *testing.T) {
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
	h := NewHandlerEnv(k8sClient)
	logger := applog.GetLogger()
	router := V2APIRouter(h, logger)

	s := &GetCodebaseSuite{
		Router: router,
	}
	suite.Run(t, s)
}

func (s *GetCodebaseSuite) SetupSuite() {
	s.TestServer = httptest.NewServer(s.Router)
}

func (s *GetCodebaseSuite) TearDownSuite() {
	s.TestServer.Close()
}

func (s *GetCodebaseSuite) RedefineK8SClientWithCodebaseCR(crCodebases []*codeBaseApi.Codebase) {
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
	testHandler := NewHandlerEnv(k8sClient)
	newRouter := V2APIRouter(testHandler, applog.GetLogger())
	s.TestServer.Config.Handler = newRouter
}

func createCodebaseCRWithOptions(opts ...CodebaseCROption) *codeBaseApi.Codebase {
	codebaseCR := new(codeBaseApi.Codebase)
	for i := range opts {
		opts[i](codebaseCR)
	}
	return codebaseCR
}

type CodebaseCROption func(codebase *codeBaseApi.Codebase)

func WithCrName(crName string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.ObjectMeta.Name = crName
	}
}

func WithCrNamespace(crNamespace string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.ObjectMeta.Namespace = crNamespace
	}
}

func WithGitServerName(gitServerName string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.GitServer = gitServerName
	}
}

func WithBuildTool(buildTool string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.BuildTool = buildTool
	}
}

func WithSpecType(specType string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.Type = specType
	}
}

func WithJenkinsSlave(jenkinsSlave *string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.JenkinsSlave = jenkinsSlave
	}
}

func (s *GetCodebaseSuite) TestGetCodebase_OK() {
	t := s.T()

	namespaceName := testNamespace
	crCodebaseName_1 := "fake_spring-petclinic"
	crCodebaseGitServer_1 := "gerrit"
	npmBuildTool := "npm"
	specType := "application"
	jenkinsSlave := "npm"

	stubCodebase_1 := createCodebaseCRWithOptions(
		WithCrNamespace(namespaceName),
		WithCrName(crCodebaseName_1),
		WithGitServerName(crCodebaseGitServer_1),
		WithBuildTool(npmBuildTool),
		WithSpecType(specType),
		WithJenkinsSlave(&jenkinsSlave),
	)
	stubCodebases := []*codeBaseApi.Codebase{
		stubCodebase_1,
	}
	s.RedefineK8SClientWithCodebaseCR(stubCodebases)
	httpExpect := httpexpect.New(t, s.TestServer.URL)
	response := httpExpect.
		GET(fmt.Sprintf("/api/v2/edp/codebase/%s", "fake_spring-petclinic")).
		Expect().
		Status(http.StatusOK).
		ContentType("application/json")

	expectedJSONBody := fmt.Sprintf(`{
	"build_tool" : "%s",
	"ciTool": "",
	"codebase_branch": null,
	"commitMessagePattern": "",
	"defaultBranch": "",
	"deploymentScript": "",
	"description": "",
	"emptyProject": false,
    "framework": "",
    "gitProjectPath": null,
    "gitServer": "%s",
    "git_url": "",
    "id": 0,
    "jenkinsSlave": "%s",
    "jiraIssueFields": null,
    "jiraServer": null,
    "jobProvisioning": "",
    "language": "",
    "name": "%s",
    "perf": null,
    "startFrom": null,
    "status": "",
    "strategy": "",
    "testReportFramework": "",
    "ticketNamePattern": "",
    "type": "%s",
    "versioningType": ""
}
`,
		npmBuildTool, crCodebaseGitServer_1, jenkinsSlave,
		crCodebaseName_1, specType,
	)

	gotPlainBody := response.Body().Raw()
	assert.JSONEq(t, expectedJSONBody, gotPlainBody, "unexpected body")
}

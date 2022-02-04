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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

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

func createCodebaseCR(namespace, crName, gitServerName string) *codeBaseApi.Codebase {
	return &codeBaseApi.Codebase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName,
			Namespace: namespace,
		},
		Spec: codeBaseApi.CodebaseSpec{
			GitServer: gitServerName,
		},
	}
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
	h := NewHandlerEnv(k8sClient)
	logger := applog.GetLogger()
	router := V2APIRouter(h, logger)

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
	testHandler := NewHandlerEnv(k8sClient)
	newRouter := V2APIRouter(testHandler, applog.GetLogger())
	s.TestServer.Config.Handler = newRouter
}

func (s *GetCodebasesSuite) TestGetCodebases_OK() {
	t := s.T()

	namespaceName := testNamespace
	crCodebaseName_1 := "fake_spring-petclinic"
	crCodebaseGitServer_1 := "gerrit"
	crCodebaseName_2 := "fake_2"
	crCodebaseGitServer_2 := "gitlab"

	stubCodebase_1 := createCodebaseCR(namespaceName, crCodebaseName_1, crCodebaseGitServer_1)
	stubCodebase_2 := createCodebaseCR(namespaceName, crCodebaseName_2, crCodebaseGitServer_2)
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
	"build_tool" : "",
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
    "jenkinsSlave": "",
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
    "type": "",
    "versioningType": ""
},
{
	"build_tool" : "",
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
    "jenkinsSlave": "",
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
    "type": "",
    "versioningType": ""
}
]`,
		crCodebaseGitServer_1, crCodebaseName_1,
		crCodebaseGitServer_2, crCodebaseName_2,
	)

	gotPlainBody := response.Body().Raw()
	assert.JSONEq(t, expectedJSONBody, gotPlainBody, "unexpected body")
}

package webapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/epam/edp-codebase-operator/v2/pkg/codebasebranch"
	"github.com/gavv/httpexpect/v2"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
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
	h := NewHandlerEnv(WithClient(k8sClient))
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

type crObjects struct {
	crCodebase       []*codeBaseApi.Codebase
	crCodebaseBranch []*codeBaseApi.CodebaseBranch
}

func (s *GetCodebaseSuite) RedefineK8SClientWithCodebaseCR(stubObjects crObjects) {
	runtimeScheme := runtime.NewScheme()
	runtimeScheme.AddKnownTypes(appsv1.SchemeGroupVersion,
		&codeBaseApi.Codebase{}, &codeBaseApi.CodebaseBranch{},
		&codeBaseApi.CodebaseBranchList{},
	)
	namespaceName := testNamespace

	builder := fake.NewClientBuilder().WithScheme(runtimeScheme)
	if len(stubObjects.crCodebase) > 0 {
		fakeObjects := make([]runtime.Object, 0)
		for _, crCodebase := range stubObjects.crCodebase {
			fakeObjects = append(fakeObjects, crCodebase)
		}
		builder = builder.WithRuntimeObjects(fakeObjects...)
	}

	if len(stubObjects.crCodebaseBranch) > 0 {
		fakeObjects := make([]runtime.Object, 0)
		for _, crCodebaseBranch := range stubObjects.crCodebaseBranch {
			fakeObjects = append(fakeObjects, crCodebaseBranch)
		}
		builder = builder.WithRuntimeObjects(fakeObjects...)
	}

	fakeClient := builder.Build()
	k8sClient, err := k8s.NewRuntimeNamespacedClient(fakeClient, namespaceName)
	if err != nil {
		s.T().Fatal(err)
	}
	testHandler := NewHandlerEnv(WithClient(k8sClient))
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

func WithLanguage(language string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.Lang = language
	}
}

func WithVersioningType(versioningType string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.Versioning.Type = codeBaseApi.VersioningType(versioningType)
	}
}

func WithDeploymentScript(deploymentScript string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.DeploymentScript = deploymentScript
	}
}

func WithFramework(framework *string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.Framework = framework
	}
}

func WithJobProvisioning(jobProvisioning *string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.JobProvisioning = jobProvisioning
	}
}

func WithStrategy(strategy string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.Strategy = codeBaseApi.Strategy(strategy)
	}
}

func WithCommitMessagePattern(pattern string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.Spec.CommitMessagePattern = &pattern
	}
}

func createCodebaseBranchCRWithOptions(opts ...CodebaseBranchCROption) *codeBaseApi.CodebaseBranch {
	codebaseBranchCR := new(codeBaseApi.CodebaseBranch)
	for i := range opts {
		opts[i](codebaseBranchCR)
	}
	return codebaseBranchCR
}

func cbBranchWithName(name string) CodebaseBranchCROption {
	return func(cbBranch *codeBaseApi.CodebaseBranch) {
		cbBranch.Name = name
	}
}

func cbBranchWithNamespace(namespace string) CodebaseBranchCROption {
	return func(cbBranch *codeBaseApi.CodebaseBranch) {
		cbBranch.Namespace = namespace
	}
}

func cbBranchWithLabels(labels map[string]string) CodebaseBranchCROption {
	return func(cbBranch *codeBaseApi.CodebaseBranch) {
		cbBranch.Labels = labels
	}
}

func cbBranchWithSpecBranchName(branchName string) CodebaseBranchCROption {
	return func(cbBranch *codeBaseApi.CodebaseBranch) {
		cbBranch.Spec.BranchName = branchName
	}
}

func cbBranchWithBuildNumber(buildNumber *string) CodebaseBranchCROption {
	return func(cbBranch *codeBaseApi.CodebaseBranch) {
		cbBranch.Status.Build = buildNumber
	}
}

func cbBranchWithIsRelease(isRelease bool) CodebaseBranchCROption {
	return func(cbBranch *codeBaseApi.CodebaseBranch) {
		cbBranch.Spec.Release = isRelease
	}
}

func cbBranchWithVersion(version *string) CodebaseBranchCROption {
	return func(cbBranch *codeBaseApi.CodebaseBranch) {
		cbBranch.Spec.Version = version
	}
}

type CodebaseBranchCROption func(codebase *codeBaseApi.CodebaseBranch)

func (s *GetCodebaseSuite) TestGetCodebase_OK() {
	t := s.T()

	namespaceName := testNamespace
	crCodebaseName_1 := "fake_spring-petclinic"
	crCodebaseGitServer_1 := "gerrit"
	npmBuildTool := "npm"
	specType := "application"
	jenkinsSlave := "npm"
	language := "javascript"
	versioningType := "default"
	deploymentScript := "helm-chart"
	framework := "java8"
	jobProvisioning := "default"
	strategy := "create"
	commitMessagePattern := `[JIRA\-\d{4}]\s+test|feat`

	stubCodebase_1 := createCodebaseCRWithOptions(
		WithCrNamespace(namespaceName),
		WithCrName(crCodebaseName_1),
		WithGitServerName(crCodebaseGitServer_1),
		WithBuildTool(npmBuildTool),
		WithSpecType(specType),
		WithJenkinsSlave(&jenkinsSlave),
		WithLanguage(language),
		WithVersioningType(versioningType),
		WithDeploymentScript(deploymentScript),
		WithFramework(&framework),
		WithJobProvisioning(&jobProvisioning),
		WithStrategy(strategy),
		WithCommitMessagePattern(commitMessagePattern),
	)
	stubCodebases := []*codeBaseApi.Codebase{
		stubCodebase_1,
	}

	cbBranchName_1 := "develop"
	crCbBranchName := fmt.Sprintf("%s-%s", crCodebaseName_1, cbBranchName_1)
	cbLabels_1 := map[string]string{
		codebasebranch.LabelCodebaseName: crCodebaseName_1,
	}
	buildNumber_1 := "1"
	isRelease_1 := true
	cbVersion_1 := "release-v1.2.5"
	stubCodebaseBranch_1 := createCodebaseBranchCRWithOptions(
		cbBranchWithName(crCbBranchName),
		cbBranchWithNamespace(namespaceName),
		cbBranchWithLabels(cbLabels_1),
		cbBranchWithSpecBranchName(cbBranchName_1),
		cbBranchWithBuildNumber(&buildNumber_1),
		cbBranchWithIsRelease(isRelease_1),
		cbBranchWithVersion(&cbVersion_1),
	)
	stubCbBranches := []*codeBaseApi.CodebaseBranch{
		stubCodebaseBranch_1,
	}

	stubObjects := crObjects{
		crCodebase:       stubCodebases,
		crCodebaseBranch: stubCbBranches,
	}
	s.RedefineK8SClientWithCodebaseCR(stubObjects)
	httpExpect := httpexpect.New(t, s.TestServer.URL)
	response := httpExpect.
		GET(fmt.Sprintf("/api/v2/edp/codebase/%s", "fake_spring-petclinic")).
		Expect().
		Status(http.StatusOK).
		ContentType("application/json")

	escapedJSONCommitMessagePattern, err := json.Marshal(commitMessagePattern)
	if err != nil {
		t.Fatal(err)
	}
	expectedJSONBody := fmt.Sprintf(`{
	"build_tool" : "%s",
    "codebase_branch": [
        {
             "branchName": "%s",
             "build_number": "%s",
             "release": %t,
             "version": "%s"
        }
    ],
    "commitMessagePattern": %s,
	"deploymentScript": "%s",
	"emptyProject": false,
    "framework": "%s",
    "gitServer": "%s",
    "jenkinsSlave": "%s",
    "jobProvisioning": "%s",
    "language": "%s",
    "name": "%s",
    "strategy": "%s",
    "type": "%s",
    "versioningType": "%s"
}
`,
		npmBuildTool, cbBranchName_1, buildNumber_1, isRelease_1, cbVersion_1, string(escapedJSONCommitMessagePattern),
		deploymentScript, framework, crCodebaseGitServer_1, jenkinsSlave, jobProvisioning,
		language, crCodebaseName_1, strategy, specType, versioningType,
	)

	gotPlainBody := response.Body().Raw()
	assert.JSONEq(t, expectedJSONBody, gotPlainBody, "unexpected body")
}

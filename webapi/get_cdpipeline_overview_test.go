package webapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	"github.com/epam/edp-codebase-operator/v2/pkg/codebasebranch"
	jenkinsApi "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1alpha1"
	"github.com/gavv/httpexpect/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

func TestCDPipelineOverview_OK(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion,
		&jenkinsApi.Jenkins{}, &cdPipeApi.CDPipeline{},
		&codeBaseApi.Codebase{}, &codeBaseApi.CodebaseList{},
		&codeBaseApi.CodebaseBranch{}, &codeBaseApi.CodebaseBranchList{},
		&cdPipeApi.CDPipelineList{},
	)

	crCodebaseName_1 := "test_codebase"
	codebaseCR_1 := createCodebaseCRWithOptions(
		WithCrNamespace(namespace),
		WithCrName(crCodebaseName_1),
	)

	codebaseBranchCRName_1 := "test_codebase_branch"
	codebaseBranchCR_1 := createCodebaseBranchCRWithOptions(
		cbBranchWithName(codebaseBranchCRName_1),
		cbBranchWithNamespace(namespace),
	)
	cdPipelineCRName := "test_cd"
	cdPipelineCR := createCDPipelineCRWithOptions(
		WithCDPipelineNamespace(namespace),
		WithCDPipelineName(cdPipelineCRName),
	)

	jenkinsCRName := "jenkins"
	jenkinsUrl := "example.com"
	jenkinsCR := createJenkinsCRWithOptions(
		WithJenkinsNamespace(namespace),
		WithJenkinsName(jenkinsCRName),
		WithJenkinsDNSWildcard(jenkinsUrl),
	)

	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(codebaseCR_1, codebaseBranchCR_1, cdPipelineCR, jenkinsCR).Build()

	err := codebasebranch.AddCodebaseLabel(ctx, client, codebaseBranchCR_1, crCodebaseName_1)
	if err != nil {
		t.Fatal(err)
	}

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
		BasePath:    "/",
		AuthEnable:  false,
		EDPVersion:  "v1",
		XSRFEnabled: false,
	}

	h := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMap(CreateCommonFuncMap()), WithConfig(conf))
	logger := applog.GetLogger()
	authHandler := HandlerAuthWithOption()
	router := V2APIRouter(h, authHandler, logger)
	testServer := httptest.NewServer(router)

	httpExpect := httpexpect.New(t, testServer.URL)
	httpExpect.
		GET("/v2/admin/edp/cd-pipeline/overview").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html")
}

type JenkinsCROption func(jenkins *jenkinsApi.Jenkins)

func createJenkinsCRWithOptions(opts ...JenkinsCROption) *jenkinsApi.Jenkins {
	jenkinsCR := new(jenkinsApi.Jenkins)
	for i := range opts {
		opts[i](jenkinsCR)
	}
	return jenkinsCR
}

func WithJenkinsNamespace(namespace string) JenkinsCROption {
	return func(jenkins *jenkinsApi.Jenkins) {
		jenkins.Namespace = namespace
	}
}

func WithJenkinsName(name string) JenkinsCROption {
	return func(jenkins *jenkinsApi.Jenkins) {
		jenkins.Name = name
	}
}

func WithJenkinsDNSWildcard(url string) JenkinsCROption {
	return func(jenkins *jenkinsApi.Jenkins) {
		jenkins.Spec.EdpSpec.DnsWildcard = url
	}
}

type CDPipelineCROption func(cdPipe *cdPipeApi.CDPipeline)

func createCDPipelineCRWithOptions(opts ...CDPipelineCROption) *cdPipeApi.CDPipeline {
	cdPipeCR := new(cdPipeApi.CDPipeline)
	for i := range opts {
		opts[i](cdPipeCR)
	}
	return cdPipeCR
}

func WithCDPipelineNamespace(namespace string) CDPipelineCROption {
	return func(cdPipe *cdPipeApi.CDPipeline) {
		cdPipe.Namespace = namespace
	}
}

func WithCDPipelineName(name string) CDPipelineCROption {
	return func(cdPipe *cdPipeApi.CDPipeline) {
		cdPipe.Name = name
	}
}

func WithCDPipelineInputDockerStreams(inputDockerStreams []string) CDPipelineCROption {
	return func(cdPipe *cdPipeApi.CDPipeline) {
		cdPipe.Spec.InputDockerStreams = inputDockerStreams
	}
}

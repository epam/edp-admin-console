package webapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	cdPipelineAPI "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
	"edp-admin-console/util/consts"
)

func TestPostCodebase_OK(t *testing.T) {
	scheme := runtime.NewScheme()
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	err = cdPipelineAPI.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}

	crCodebaseName_1 := "test_goapp"
	crCodebaseType_1 := "application"
	codebaseCR_1 := createCodebaseCRWithOptions(
		WithCrNamespace(namespace),
		WithCrName(crCodebaseName_1),
		WithSpecType(crCodebaseType_1),
	)

	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(codebaseCR_1).Build()
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
		XSRFKey:     []byte("secret"),
		XSRFEnabled: false,
	}
	h := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMap(CreateCommonFuncMap()), WithConfig(conf))
	logger := applog.GetLogger()
	authHandler := HandlerAuthWithOption()
	router := V2APIRouter(h, authHandler, logger)
	testServer := httptest.NewServer(router)

	codebaseName := crCodebaseName_1
	codebaseType := crCodebaseType_1

	httpExpect := httpexpect.New(t, testServer.URL)
	result := httpExpect.POST("/v2/admin/edp/codebase").
		WithFormField("name", codebaseName).
		WithFormField("codebase-type", codebaseType).
		WithRedirectPolicy(httpexpect.DontFollowRedirects)
	result.Expect().
		Status(http.StatusFound).
		Header("Location").
		Equal(fmt.Sprintf("/v2/admin/edp/application/overview?codebase=%s#codebaseIsDeleted", codebaseName))
}

type CodebaseImageStreanCROption func(codebase *codeBaseApi.CodebaseImageStream)

func createCodebaseImageStreamCRWithOptions(opts ...CodebaseImageStreanCROption) *codeBaseApi.CodebaseImageStream {
	codebaseImageStream := new(codeBaseApi.CodebaseImageStream)
	for i := range opts {
		opts[i](codebaseImageStream)
	}
	return codebaseImageStream
}

func WithCBISWithName(name string) CodebaseImageStreanCROption {
	return func(cbBranch *codeBaseApi.CodebaseImageStream) {
		cbBranch.Name = name
	}
}

func WithCBISNamespace(namespace string) CodebaseImageStreanCROption {
	return func(cbBranch *codeBaseApi.CodebaseImageStream) {
		cbBranch.Namespace = namespace
	}
}

func WithCBISSpecCodebaseName(codebaseName string) CodebaseImageStreanCROption {
	return func(cbBranch *codeBaseApi.CodebaseImageStream) {
		cbBranch.Spec.Codebase = codebaseName
	}
}

func TestPostCodebase_IsInUse(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion,
		&codeBaseApi.Codebase{}, &codeBaseApi.CodebaseList{},
		&cdPipelineAPI.CDPipeline{}, &cdPipelineAPI.CDPipelineList{},
		&codeBaseApi.CodebaseBranch{}, &codeBaseApi.CodebaseBranchList{},
		&codeBaseApi.CodebaseImageStream{}, &codeBaseApi.CodebaseImageStreamList{},
	)

	crCodebaseName_1 := "test_goapp"
	crCodebaseType_1 := "application"
	codebaseCR_1 := createCodebaseCRWithOptions(
		WithCrNamespace(namespace),
		WithCrName(crCodebaseName_1),
		WithSpecType(crCodebaseType_1),
	)
	crCBIS_1 := "test_goapp-develop"
	cbBranchCR_1 := createCodebaseImageStreamCRWithOptions(
		WithCBISWithName(crCBIS_1),
		WithCBISNamespace(namespace),
		WithCBISSpecCodebaseName(crCodebaseName_1),
	)

	cdPipelineName_1 := "cd_pipeline_1"
	inputDockerStreams := []string{crCBIS_1}
	cdPipelineCR := createCDPipelineCRWithOptions(
		WithCDPipelineNamespace(namespace),
		WithCDPipelineName(cdPipelineName_1),
		WithCDPipelineInputDockerStreams(inputDockerStreams),
	)

	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(codebaseCR_1, cdPipelineCR, cbBranchCR_1).Build()
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
		XSRFKey:     []byte("secret"),
		XSRFEnabled: false,
	}
	h := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMap(CreateCommonFuncMap()), WithConfig(conf))
	logger := applog.GetLogger()
	authHandler := HandlerAuthWithOption()
	router := V2APIRouter(h, authHandler, logger)
	testServer := httptest.NewServer(router)

	codebaseName := crCodebaseName_1
	codebaseType := crCodebaseType_1

	httpExpect := httpexpect.New(t, testServer.URL)
	result := httpExpect.POST("/v2/admin/edp/codebase").
		WithFormField("name", codebaseName).
		WithFormField("codebase-type", codebaseType).
		WithRedirectPolicy(httpexpect.DontFollowRedirects)
	errmsg := fmt.Sprintf(
		"application '%s' and its image stream '%s' are used by '%s' pipeline",
		codebaseName, crCBIS_1, cdPipelineName_1)
	result.Expect().
		Status(http.StatusFound).
		Header("Location").
		Equal(fmt.Sprintf("/v2/admin/edp/application/overview?codebase=%s&errmsg=%s#codebaseIsUsed", codebaseName, errmsg))

	k8sCodebase, err := namespacedClient.GetCodebase(ctx, codebaseName)
	assert.NoError(t, err)

	expectedCodebaseCR := &codeBaseApi.Codebase{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Codebase",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            crCodebaseName_1,
			Namespace:       namespace,
			ResourceVersion: "999",
		},
		Spec: codeBaseApi.CodebaseSpec{
			Type: consts.Application,
		},
		Status: codeBaseApi.CodebaseStatus{},
	}
	assert.Equal(t, expectedCodebaseCR, k8sCodebase, "codebase CR must exist")
}

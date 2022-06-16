package webapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"

	cdPipelineAPI "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type StageCROption func(codebase *cdPipelineAPI.Stage)

func createStageCRWithOptions(opts ...StageCROption) *cdPipelineAPI.Stage {
	stageCR := new(cdPipelineAPI.Stage)
	for i := range opts {
		opts[i](stageCR)
	}
	return stageCR
}

func stageCRWithName(crName string) StageCROption {
	return func(stageCR *cdPipelineAPI.Stage) {
		stageCR.Name = crName
	}
}

func stageCRWithNamespace(crNamespace string) StageCROption {
	return func(stageCR *cdPipelineAPI.Stage) {
		stageCR.Namespace = crNamespace
	}
}

func stageCRWithSpecCDPipeline(cdPipelineName string) StageCROption {
	return func(stageCR *cdPipelineAPI.Stage) {
		stageCR.Spec.CdPipeline = cdPipelineName
	}
}

func TestHandlerEnv_DeleteCD_OK(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion,
		&cdPipelineAPI.CDPipeline{},
		&cdPipelineAPI.Stage{}, &cdPipelineAPI.StageList{},
	)
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}

	crCDPipelineName_1 := "test_cd_pipeline_1"
	cdPipelineCR_1 := createCDPipelineCRWithOptions(
		WithCDPipelineNamespace(namespace),
		WithCDPipelineName(crCDPipelineName_1),
	)

	crStageName_1 := "stage_0"
	crStageName_2 := "stage_1"
	crStageName_3 := "stage_2"
	crStage_1 := createStageCRWithOptions(
		stageCRWithName(crStageName_1),
		stageCRWithNamespace(namespace),
		stageCRWithSpecCDPipeline(crCDPipelineName_1),
	)
	crStage_2 := createStageCRWithOptions(
		stageCRWithName(crStageName_2),
		stageCRWithNamespace(namespace),
		stageCRWithSpecCDPipeline(crCDPipelineName_1),
	)

	foreignNamespace := "foreign " + namespace
	crStage_3 := createStageCRWithOptions( // foreign custom resource
		stageCRWithName(crStageName_3),
		stageCRWithNamespace(foreignNamespace),
		stageCRWithSpecCDPipeline(crCDPipelineName_1),
	)

	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(cdPipelineCR_1, crStage_1, crStage_2, crStage_3).Build()
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
		BasePath:    "",
		AuthEnable:  false,
		XSRFKey:     []byte("secret"),
		XSRFEnabled: false,
	}
	h := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMap(CreateCommonFuncMap()), WithConfig(conf))
	logger := applog.GetLogger()
	authHandler := HandlerAuthWithOption()
	router := V2APIRouter(h, authHandler, logger)
	testServer := httptest.NewServer(router)

	httpExpect := httpexpect.New(t, testServer.URL)
	result := httpExpect.POST("/v2/admin/edp/cd-pipeline/delete").
		WithFormField("name", crCDPipelineName_1).
		WithRedirectPolicy(httpexpect.DontFollowRedirects)
	result.Expect().
		Status(http.StatusFound).
		Header("Location").Equal(fmt.Sprintf("/v2/admin/edp/cd-pipeline/overview?name=%s#cdPipelineDeletedSuccessModal", crCDPipelineName_1))

	ctx := context.Background()
	nsn := types.NamespacedName{
		Name:      crStageName_1,
		Namespace: namespace,
	}
	k8sStage_1 := new(cdPipelineAPI.Stage)
	err = client.Get(ctx, nsn, k8sStage_1)
	assert.Error(t, err, metav1.StatusReasonNotFound)

	nsn = types.NamespacedName{
		Name:      crStageName_2,
		Namespace: namespace,
	}
	k8sStage_2 := new(cdPipelineAPI.Stage)
	err = client.Get(ctx, nsn, k8sStage_2)
	assert.Error(t, err, metav1.StatusReasonNotFound)

	nsn = types.NamespacedName{
		Name:      crStageName_3,
		Namespace: foreignNamespace,
	}
	k8sStage_3 := new(cdPipelineAPI.Stage)
	err = client.Get(ctx, nsn, k8sStage_3)
	assert.NoError(t, err)
	crStage_3.TypeMeta.Kind = "Stage" // these fields are out of our control, but we need to assert them
	crStage_3.TypeMeta.APIVersion = "apps/v1"
	assert.Equal(t, crStage_3, k8sStage_3)

	nsn = types.NamespacedName{
		Name:      crCDPipelineName_1,
		Namespace: namespace,
	}
	k8sCPPipeline := new(cdPipelineAPI.CDPipeline)
	err = client.Get(ctx, nsn, k8sCPPipeline)
	assert.Error(t, err, metav1.StatusReasonNotFound)
}

func TestHandlerEnv_DeleteCD_CDPipelineNotFound(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion,
		&cdPipelineAPI.CDPipeline{},
	)
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}

	crCDPipelineName_1 := "test_cd_pipeline_1"

	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects().Build()
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
		BasePath:    "",
		AuthEnable:  false,
		XSRFKey:     []byte("secret"),
		XSRFEnabled: false,
	}
	h := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMap(CreateCommonFuncMap()), WithConfig(conf))
	logger := applog.GetLogger()
	authHandler := HandlerAuthWithOption()
	router := V2APIRouter(h, authHandler, logger)
	testServer := httptest.NewServer(router)

	httpExpect := httpexpect.New(t, testServer.URL)
	result := httpExpect.POST("/v2/admin/edp/cd-pipeline/delete").
		WithFormField("name", crCDPipelineName_1).
		WithRedirectPolicy(httpexpect.DontFollowRedirects)
	result.Expect().
		Status(http.StatusFound).
		Header("Location").Equal("/v2/admin/edp/cd-pipeline/overview#cdPipelineDeletedErrorModal")
}

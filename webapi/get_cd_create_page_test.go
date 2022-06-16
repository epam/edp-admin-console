package webapi

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	edpComponentAPI "github.com/epam/edp-component-operator/pkg/apis/v1/v1alpha1"
	jenkinsAPI "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1alpha1"
	"github.com/gavv/httpexpect/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGetCDCreatePage_OK(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion,
		&edpComponentAPI.EDPComponent{},
		&codeBaseApi.Codebase{}, &codeBaseApi.CodebaseList{},
		&codeBaseApi.CodebaseBranch{}, &codeBaseApi.CodebaseBranchList{},
		&jenkinsAPI.JenkinsList{},
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

	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(codebaseCR_1, codebaseBranchCR_1).Build()
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
		EDPVersion: "v1",
		XSRFKey:    []byte("secret"),
	}

	h := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMap(CreateCommonFuncMap()), WithConfig(conf))
	logger := applog.GetLogger()
	authHandler := HandlerAuthWithOption()
	router := V2APIRouter(h, authHandler, logger)
	testServer := httptest.NewServer(router)

	httpExpect := httpexpect.New(t, testServer.URL)
	httpExpect.
		GET("/v2/admin/edp/cd-pipeline/create").
		WithCookie("_edp_csrf", "csrf_token").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html")
}

package webapi

import (
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
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestHandlerEnv_DeleteCodebaseBranch(t *testing.T) {
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
	codebaseBranchName_1 := "develop"
	crCodebaseBranchName_1 := buildCobaseBranchCRName(crCodebaseName_1, codebaseBranchName_1)
	codebaseCR_1 := createCodebaseBranchCRWithOptions(
		cbBranchWithName(crCodebaseBranchName_1),
		cbBranchWithNamespace(namespace),
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

	cbBranchName := codebaseBranchName_1
	codebaseName := crCodebaseName_1

	httpExpect := httpexpect.New(t, testServer.URL)
	result := httpExpect.POST("/v2/admin/edp/codebase/branch/delete").
		WithFormField("name", cbBranchName).
		WithFormField("codebase-name", codebaseName).
		WithRedirectPolicy(httpexpect.DontFollowRedirects)
	result.Expect().
		Status(http.StatusFound).
		Header("Location").
		Equal(fmt.Sprintf("%s/v2/admin/edp/codebase/%s/overview?name=%s#branchDeletedSuccessModal", conf.BasePath, codebaseName, cbBranchName))
}

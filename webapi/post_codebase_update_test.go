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

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	"github.com/gavv/httpexpect/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestEditCodebase_OK(t *testing.T) {
	scheme := runtime.NewScheme()
	err := codeBaseApi.AddToScheme(scheme)
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
	jiraServerToggle := "on"
	jiraFieldNames := []string{"issueKey", "prokectKey"}
	jiraPatterns := []string{`ISSUE-\d`, `\w{2}.\d`}
	commitMessagePattern := "commit pattern"
	ticketNamePattern := "ticket name pattern"
	jiraServer := "test-jira"

	httpExpect := httpexpect.New(t, testServer.URL)
	result := httpExpect.POST(fmt.Sprintf("/v2/admin/edp/codebase/%s/update", codebaseName)).
		WithFormField("jiraServerToggle", jiraServerToggle).
		WithFormField("jiraFieldName", jiraFieldNames).
		WithFormField("name", codebaseName).
		WithFormField("jiraPattern", jiraPatterns).
		WithFormField("commitMessagePattern", commitMessagePattern).
		WithFormField("ticketNamePattern", ticketNamePattern).
		WithFormField("jiraServer", jiraServer).
		WithRedirectPolicy(httpexpect.DontFollowRedirects)
	result.Expect().
		Status(http.StatusFound).
		Header("Location").Equal("/v2/admin/edp/application/overview#codebaseUpdateSuccessModal")
}

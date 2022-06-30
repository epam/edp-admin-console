package webapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

type CreateBranchSuite struct {
	suite.Suite
	TestServer *httptest.Server
	Handler    *HandlerEnv
}

func TestCreateBranchSuite(t *testing.T) {
	scheme := runtime.NewScheme()
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	client := fake.NewClientBuilder().WithScheme(scheme).Build()
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
		XSRFEnabled: false,
	}
	h := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMap(CreateCommonFuncMap()), WithConfig(conf))
	logger := applog.GetLogger()
	authHandler := HandlerAuthWithOption()
	router := V2APIRouter(h, authHandler, logger)

	s := &CreateBranchSuite{
		Handler: h,
	}
	s.TestServer = httptest.NewServer(router)
	suite.Run(t, s)
}

func (s *CreateBranchSuite) TestCreateBranch() {
	t := s.T()
	ctx := context.Background()
	defaultBranch := "master"
	appName := "appName"
	codebaseSpec := codeBaseApi.CodebaseSpec{
		DefaultBranch: defaultBranch,
	}
	err := s.Handler.NamespacedClient.CreateCodebaseByCustomFields(ctx, appName, codebaseSpec)
	if err != nil {
		t.Fatal(err)
	}

	branchName := "new"
	version := "1.1"
	versionPostfix := "RC"
	isRelease := "true"

	httpExpect := httpexpect.New(t, s.TestServer.URL)
	req := httpExpect.POST(fmt.Sprintf("/v2/admin/edp/codebase/%s/branch", appName)).
		WithFormField("name", branchName).
		WithFormField("version", version).
		WithFormField("versioningPostfix", versionPostfix).
		WithFormField("releaseBranch", isRelease).
		WithFormField("masterVersion", version).
		WithFormField("snapshotStaticField", versionPostfix).
		WithRedirectPolicy(httpexpect.DontFollowRedirects)

	expectedURL := fmt.Sprintf("%s/v2/admin/edp/codebase/%s/overview?%s=%s#branchSuccessModal", s.Handler.Config.BasePath,
		appName, paramWaitingForBranch, url.PathEscape(branchName))
	req.Expect().
		Status(http.StatusFound).
		Header("location").
		Equal(expectedURL)

	codebase, err := s.Handler.NamespacedClient.GetCodebase(ctx, appName)
	assert.NoError(t, err)
	assert.Equal(t, *GetVersionOrNil(version, versionPostfix), codebase.Spec.DefaultBranch)

	codebaseBranch, err := s.Handler.NamespacedClient.GetCBBranch(ctx, fmt.Sprintf("%s-%s", appName, branchName))
	assert.NoError(t, err)
	assert.Equal(t, branchName, codebaseBranch.Spec.BranchName)
	assert.True(t, codebaseBranch.Spec.Release)
	assert.Equal(t, GetVersionOrNil(version, versionPostfix), codebaseBranch.Spec.Version)
	assert.Equal(t, appName, codebaseBranch.Spec.CodebaseName)

}

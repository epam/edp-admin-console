package webapi

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

type ApplicationOverviewSuite struct {
	suite.Suite
	Server  *httptest.Server
	Handler *HandlerEnv
}

func TestApplicationOverviewSuite(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion, &codeBaseApi.CodebaseList{}, &codeBaseApi.Codebase{})
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
		BasePath:   "/",
		AuthEnable: false,
		EDPVersion: "v1",
	}

	h := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMap(CreateCommonFuncMap()), WithConfig(conf))
	logger := applog.GetLogger()
	authHandler := HandlerAuthWithOption()
	router := V2APIRouter(h, authHandler, logger)

	s := &ApplicationOverviewSuite{
		Handler: h,
	}
	s.Server = httptest.NewServer(router)
	suite.Run(t, s)
}

func (s *ApplicationOverviewSuite) TestApplicationOverview() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.Server.URL)
	httpExpect.
		GET("/v2/admin/edp/application/overview").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html")
}

package webapi

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

type IndexSuite struct {
	suite.Suite
	TestServer *httptest.Server
	Handler    *HandlerEnv
}

func TestIndexSuite(t *testing.T) {
	client := fake.NewClientBuilder().Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	workingDir, _ := path.Split(currentDir)
	h := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMapTemplate(CreateCommonFuncMap()))
	logger := applog.GetLogger()
	router := V2APIRouter(h, logger)

	s := &IndexSuite{
		Handler: h,
	}
	s.TestServer = httptest.NewServer(router)
	suite.Run(t, s)
}

func (s *IndexSuite) TestIndex() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.TestServer.URL)
	httpExpect.
		GET("/v2/admin/edp/index").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html")
}

package webapi

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	jenkinsApi "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

type PipelineUpdateSuite struct {
	suite.Suite
	Server  *httptest.Server
	Handler *HandlerEnv
}

func TestPipelineUpdateSuite(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion, &codeBaseApi.CodebaseList{}, &codeBaseApi.Codebase{}, &codeBaseApi.CodebaseImageStream{})
	scheme.AddKnownTypes(jenkinsApi.SchemeGroupVersion, &jenkinsApi.JenkinsList{}, &jenkinsApi.Jenkins{})
	scheme.AddKnownTypes(cdPipeApi.SchemeGroupVersion, &cdPipeApi.CDPipeline{})

	cdPipeline := cdPipeApi.CDPipeline{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "testCDPipeline",
			Namespace: namespace,
		},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&cdPipeline).Build()
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

	s := &PipelineUpdateSuite{
		Handler: h,
	}
	s.Server = httptest.NewServer(router)
	suite.Run(t, s)
}

func (s *PipelineUpdateSuite) TestPipelineUpdatePage() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.Server.URL)
	httpExpect.
		GET("/v2/admin/edp/cd-pipeline/testCDPipeline/update").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html")
}

package webapi

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	jenkinsApi "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1alpha1"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	perfApi "github.com/epam/edp-perf-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

type CreateApplicationOverviewSuite struct {
	suite.Suite
	Server  *httptest.Server
	Handler *HandlerEnv
}

func TestCreateApplicationOverviewSuite(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion, &codeBaseApi.GitServerList{}, &codeBaseApi.GitServer{})
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion, &codeBaseApi.JiraServerList{}, &codeBaseApi.JiraServer{})
	scheme.AddKnownTypes(perfApi.SchemeGroupVersion, &perfApi.PerfServerList{}, &perfApi.PerfServer{})
	scheme.AddKnownTypes(jenkinsApi.SchemeGroupVersion, &jenkinsApi.JenkinsList{}, &jenkinsApi.Jenkins{})
	scheme.AddKnownTypes(v1.SchemeGroupVersion, &v1.ConfigMap{})

	data := make(map[string]string)
	data[perfIntegrationEnabledKey] = "false"

	configMap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      edpConfigMapName,
			Namespace: namespace,
		},
		Data: data,
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&configMap).Build()
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

	s := &CreateApplicationOverviewSuite{
		Handler: h,
	}
	s.Server = httptest.NewServer(router)
	suite.Run(t, s)
}

func (s *CreateApplicationOverviewSuite) TestCreateApplicationOverview() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.Server.URL)
	httpExpect.
		GET("/v2/admin/edp/application/create").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html")
}

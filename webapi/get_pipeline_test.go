package webapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

const (
	firstApplication  = "DOCKER IMAGE"
	secondApplication = "ANOTHER IMAGE"
)

type GetPipelineSuite struct {
	suite.Suite
	Router     *chi.Mux
	TestServer *httptest.Server
}

func TestGetPipelineSuite(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(cdPipeApi.SchemeGroupVersion, &cdPipeApi.CDPipeline{})
	err := cdPipeApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}

	cdPipeline := &cdPipeApi.CDPipeline{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      validCDPipelineName,
			Namespace: testNamespace,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			ApplicationsToPromote: []string{firstApplication, secondApplication},
		},
		Status: cdPipeApi.CDPipelineStatus{},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(cdPipeline).Build()
	k8sClient, err := k8s.NewRuntimeNamespacedClient(client, testNamespace)
	if err != nil {
		t.Fatal(err)
	}
	h := NewHandlerEnv(k8sClient)
	logger := applog.GetLogger()
	router := V2APIRouter(h, logger)
	s := &GetPipelineSuite{
		Router: router,
	}
	suite.Run(t, s)
}

func (s *GetPipelineSuite) SetupSuite() {
	s.TestServer = httptest.NewServer(s.Router)
}
func (s *GetPipelineSuite) TearDownSuite() {
	s.TestServer.Close()
}

func (s *GetPipelineSuite) TestGetStagePipeline_ValidFirst() {
	t := s.T()
	httpExpect := httpexpect.New(t, s.TestServer.URL)
	response := httpExpect.
		GET("/api/v2/edp/cd-pipeline/" + validCDPipelineName).
		Expect().
		Status(http.StatusOK).
		ContentType("application/json")

	expectedJSONBody := fmt.Sprintf(`{
   "codebaseBranches":[
      {
         "appName":"%s"
      },
      {
         "appName":"%s"
      }
   ],
   "applicationsToPromote":[
      "%s",
      "%s"
   ]
}`,
		firstApplication, secondApplication,
		firstApplication, secondApplication)

	gotPlainBody := response.Body().Raw()
	assert.JSONEq(t, expectedJSONBody, gotPlainBody, "unexpected body")
}

func (s *GetPipelineSuite) TestGetStagePipeline_CDPipelineNotFound() {
	t := s.T()
	httpExpect := httpexpect.New(t, s.TestServer.URL)
	response := httpExpect.
		GET("/api/v2/edp/cd-pipeline/" + "unexpectedCDPipeline").
		Expect().
		Status(http.StatusInternalServerError)

	expectedBody := "get CDPipeline by name failed"
	assert.Equal(t, expectedBody, response.Body().Raw(), unexpectedBody)
}

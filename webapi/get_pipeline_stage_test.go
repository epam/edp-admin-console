package webapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/internal/imagestream"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

const (
	validCDPipelineName = "cd1"
	namespace           = "ns"
	firstStage          = "first"
	unexpectedBody      = "unexpected body"
)

type StagePipelineSuite struct {
	suite.Suite
	TestServer *httptest.Server
	Handler    *HandlerEnv
}

func TestStagePipelineSuite(t *testing.T) {
	scheme := runtime.NewScheme()
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal("failed to add codebaseApi")
	}
	err = cdPipeApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal("failed to add cdPipeApi")
	}

	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	h := NewHandlerEnv(namespacedClient)
	logger := applog.GetLogger()
	router := V2APIRouter(h, logger)

	s := &StagePipelineSuite{
		Handler: h,
	}
	s.TestServer = httptest.NewServer(router)
	suite.Run(t, s)
}

func (s *StagePipelineSuite) TearDownSuite() {
	s.TestServer.Close()
}

func (s *StagePipelineSuite) TestGetStagePipeline_StageNotFound() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.TestServer.URL)
	response := httpExpect.
		GET("/api/v2/edp/cd-pipeline/build_pipeline/stage/tests").
		Expect().
		Status(http.StatusNotFound)

	expectedBody := "stage not found"
	assert.Equal(t, expectedBody, response.Body().Raw(), unexpectedBody)
}

func (s *StagePipelineSuite) TestGetStagePipeline_ValidFirst() {
	validFirstStageName := validCDPipelineName + "-" + firstStage
	t := s.T()
	ctx := context.Background()

	validCDPipeline := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      validCDPipelineName,
			Namespace: namespace,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			Name:                  validCDPipelineName,
			InputDockerStreams:    []string{"java11-cd1-master"},
			ApplicationsToPromote: []string{"java11-cd1"},
		},
	}
	validFirstStage := cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      validFirstStageName,
			Namespace: namespace,
		},
		Spec: cdPipeApi.StageSpec{
			CdPipeline:      validCDPipelineName,
			Description:     "description",
			JobProvisioning: "manual",
			Order:           0,
			TriggerType:     "Manual",
		},
	}
	fakeK8SClient := s.Handler.NamespacedClient
	err := fakeK8SClient.Create(ctx, &validCDPipeline)
	if err != nil {
		t.Fatal(err)
	}
	err = fakeK8SClient.Create(ctx, &validFirstStage)
	if err != nil {
		t.Fatal(err)
	}

	httpExpect := httpexpect.New(t, s.TestServer.URL)
	response := httpExpect.
		GET("/api/v2/edp/cd-pipeline/" + validCDPipelineName + "/stage/" + firstStage).
		Expect().
		Status(http.StatusOK).
		ContentType("application/json")

	expectedBody := `
{"name":"first",
"cdPipeline":"cd1",
"description":"description",
"triggerType":"Manual",
"order":"0",
"applications":
	[{"name":"java11-cd1",
	"branchName":"",
	"inputIs":"java11-cd1-master",
	"outputIs":"cd1-first-java11-cd1-verified"}],
"qualityGates":null,
"jobProvisioning":"manual"}`

	assert.JSONEq(t, expectedBody, response.Body().Raw(), unexpectedBody)

	err = fakeK8SClient.Delete(ctx, &validCDPipeline)
	if err != nil {
		t.Fatal(err)
	}
	err = fakeK8SClient.Delete(ctx, &validFirstStage)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *StagePipelineSuite) TestGetStagePipeline_ValidSecond() {
	secondStage := "second"
	validSecondStageName := validCDPipelineName + "-" + secondStage
	CISNameForSecondStage := "cd1-first-java11-cd1-verified"
	t := s.T()
	ctx := context.Background()
	validCDPipeline := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      validCDPipelineName,
			Namespace: namespace,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			Name:                  validCDPipelineName,
			InputDockerStreams:    []string{"java11-cd1-master"},
			ApplicationsToPromote: []string{"java11-cd1"},
		},
	}
	validSecondStage := cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:        validSecondStageName,
			Namespace:   namespace,
			Annotations: map[string]string{imagestream.PreviousStageNameAnnotationKey: firstStage},
		},
		Spec: cdPipeApi.StageSpec{
			CdPipeline:      validCDPipelineName,
			Description:     "description",
			JobProvisioning: "manual",
			Order:           1,
			TriggerType:     "Manual",
		},
	}
	cisForSecondStage := codeBaseApi.CodebaseImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CISNameForSecondStage,
			Namespace: namespace,
		},
	}
	fakeK8SClient := s.Handler.NamespacedClient
	err := fakeK8SClient.Create(ctx, &validCDPipeline)
	if err != nil {
		t.Fatal(err)
	}
	err = fakeK8SClient.Create(ctx, &validSecondStage)
	if err != nil {
		t.Fatal(err)
	}
	err = fakeK8SClient.Create(ctx, &cisForSecondStage)
	if err != nil {
		t.Fatal(err)
	}

	httpExpect := httpexpect.New(t, s.TestServer.URL)
	response := httpExpect.
		GET("/api/v2/edp/cd-pipeline/" + validCDPipelineName + "/stage/" + secondStage).
		Expect().
		Status(http.StatusOK).
		ContentType("application/json")

	expectedBody := `
{"name":"second",
"cdPipeline":"cd1",
"description":"description",
"triggerType":"Manual",
"order":"1",
"applications":
	[{"name":"java11-cd1",
	"branchName":"",
	"inputIs":"cd1-first-java11-cd1-verified",
	"outputIs":"cd1-second-java11-cd1-verified"}],
"qualityGates":null,
"jobProvisioning":"manual"}`
	assert.JSONEq(t, expectedBody, response.Body().Raw(), unexpectedBody)

	err = fakeK8SClient.Delete(ctx, &validCDPipeline)
	if err != nil {
		t.Fatal(err)
	}
	err = fakeK8SClient.Delete(ctx, &validSecondStage)
	if err != nil {
		t.Fatal(err)
	}
	err = fakeK8SClient.Delete(ctx, &cisForSecondStage)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *StagePipelineSuite) TestGetStagePipeline_CDPipeNotFound() {
	notValidCDPipelineName := "cd2"
	nameCollision := notValidCDPipelineName + "-" + firstStage
	t := s.T()
	ctx := context.Background()
	nameCollisionStage := cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nameCollision,
			Namespace: namespace,
		},
	}
	fakeK8SClient := s.Handler.NamespacedClient
	err := fakeK8SClient.Create(ctx, &nameCollisionStage)
	if err != nil {
		t.Fatal(err)
	}

	httpExpect := httpexpect.New(t, s.TestServer.URL)
	response := httpExpect.
		GET("/api/v2/edp/cd-pipeline/" + notValidCDPipelineName + "/stage/" + firstStage).
		Expect().
		Status(http.StatusNotFound)

	expectedBody := `cdPipeline not found`
	assert.Equal(t, expectedBody, response.Body().Raw(), unexpectedBody)

	err = fakeK8SClient.Delete(ctx, &nameCollisionStage)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *StagePipelineSuite) TestGetStagePipeline_inputISNotFound() {
	emptyInputISCDPipeName := "emptyInputIS"
	emptyInputISStageName := emptyInputISCDPipeName + "-" + firstStage
	t := s.T()
	ctx := context.Background()
	emptyInputISCDPipe := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      emptyInputISCDPipeName,
			Namespace: namespace,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			Name: emptyInputISCDPipeName,
		},
	}
	emptyInputISStage := cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      emptyInputISStageName,
			Namespace: namespace,
		},
	}
	fakeK8SClient := s.Handler.NamespacedClient
	err := fakeK8SClient.Create(ctx, &emptyInputISStage)
	if err != nil {
		t.Fatal(err)
	}
	err = fakeK8SClient.Create(ctx, &emptyInputISCDPipe)
	if err != nil {
		t.Fatal(err)
	}

	httpExpect := httpexpect.New(t, s.TestServer.URL)
	response := httpExpect.
		GET("/api/v2/edp/cd-pipeline/" + emptyInputISCDPipeName + "/stage/" + firstStage).
		Expect().
		Status(http.StatusNotFound)

	expectedBody := `input IS not found`
	assert.Equal(t, expectedBody, response.Body().Raw(), unexpectedBody)

	err = fakeK8SClient.Delete(ctx, &emptyInputISCDPipe)
	if err != nil {
		t.Fatal(err)
	}
	err = fakeK8SClient.Delete(ctx, &emptyInputISStage)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *StagePipelineSuite) TestGetStagePipeline_outputISNotFound() {
	emptyOutputISCDPipeName := "emptyOutputIS"
	emptyOutputISStageName := emptyOutputISCDPipeName + "-" + firstStage
	t := s.T()
	ctx := context.Background()
	emptyOutputISCDPipe := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      emptyOutputISCDPipeName,
			Namespace: namespace,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			Name:               validCDPipelineName,
			InputDockerStreams: []string{"java11-cd1-master"},
		},
	}

	emptyOutputISStage := cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      emptyOutputISStageName,
			Namespace: namespace,
		},
		Spec: cdPipeApi.StageSpec{
			CdPipeline: emptyOutputISCDPipeName,
		},
	}
	fakeK8SClient := s.Handler.NamespacedClient
	err := fakeK8SClient.Create(ctx, &emptyOutputISStage)
	if err != nil {
		t.Fatal(err)
	}
	err = fakeK8SClient.Create(ctx, &emptyOutputISCDPipe)
	if err != nil {
		t.Fatal(err)
	}

	httpExpect := httpexpect.New(t, s.TestServer.URL)
	response := httpExpect.
		GET("/api/v2/edp/cd-pipeline/" + emptyOutputISCDPipeName + "/stage/" + firstStage).
		Expect().
		Status(http.StatusNotFound)

	expectedBody := `output IS not found`
	assert.Equal(t, expectedBody, response.Body().Raw(), unexpectedBody)

	err = fakeK8SClient.Delete(ctx, &emptyOutputISStage)
	if err != nil {
		t.Fatal(err)
	}
	err = fakeK8SClient.Delete(ctx, &emptyOutputISCDPipe)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *StagePipelineSuite) TestGetStagePipeline_ISDifferentSize() {
	diffISSizeCDPipeName := "diffSize"
	diffISSizeStageName := diffISSizeCDPipeName + "-" + firstStage
	t := s.T()
	ctx := context.Background()
	diffSizeISCDPipe := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      diffISSizeCDPipeName,
			Namespace: namespace,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			Name:                  validCDPipelineName,
			InputDockerStreams:    []string{"java11-cd1-master", "new"},
			ApplicationsToPromote: []string{"java11-cd1"},
		},
	}
	diffISSizeStage := cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      diffISSizeStageName,
			Namespace: namespace,
		},
		Spec: cdPipeApi.StageSpec{
			CdPipeline: diffISSizeCDPipeName,
		},
	}
	fakeK8SClient := s.Handler.NamespacedClient
	err := fakeK8SClient.Create(ctx, &diffSizeISCDPipe)
	if err != nil {
		t.Fatal(err)
	}
	err = fakeK8SClient.Create(ctx, &diffISSizeStage)
	if err != nil {
		t.Fatal(err)
	}

	httpExpect := httpexpect.New(t, s.TestServer.URL)
	response := httpExpect.
		GET("/api/v2/edp/cd-pipeline/" + diffISSizeCDPipeName + "/stage/" + firstStage).
		Expect().
		Status(http.StatusNotFound)

	expectedBody := `inputIS and outputIS not the same size`
	assert.Equal(t, expectedBody, response.Body().Raw(), unexpectedBody)

	err = fakeK8SClient.Delete(ctx, &diffISSizeStage)
	if err != nil {
		t.Fatal(err)
	}
	err = fakeK8SClient.Delete(ctx, &diffSizeISCDPipe)
	if err != nil {
		t.Fatal(err)
	}
}

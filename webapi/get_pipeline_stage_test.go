package webapi

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	applog "edp-admin-console/service/logger"
)

type StagePipelineSuite struct {
	suite.Suite
	Router *chi.Mux
}

func TestStagePipelineSuite(t *testing.T) {
	h := &HandlerEnv{}
	logger := applog.GetLogger()
	router := V2APIRouter(h, logger)

	s := &StagePipelineSuite{
		Router: router,
	}
	suite.Run(t, s)
}

func (s *StagePipelineSuite) TestGetStagePipeline_OK() {
	t := s.T()
	router := s.Router

	gotResponse := makeHTTPRequest(t, router, http.MethodGet, "/api/v2/edp/cd-pipeline/build_pipeline/stage/tests", nil)

	expectedHTTPCode := http.StatusOK
	expectedJSONBody := `{
	"name": "tests",
	"cdPipeline": "build_pipeline"
}`
	assert.Equal(t, expectedHTTPCode, gotResponse.statusCode, "unexpected http code")
	assert.JSONEq(t, expectedJSONBody, gotResponse.bodyBuffer.String(), "unexpected body")
}

type testResponse struct {
	statusCode int
	bodyBuffer *bytes.Buffer
}

func makeHTTPRequest(t *testing.T, router http.Handler, httpMethod, URI string, body io.Reader) *testResponse {
	t.Helper()

	responseRecorder := httptest.NewRecorder()
	testRequest := httptest.NewRequest(httpMethod, URI, body)

	router.ServeHTTP(responseRecorder, testRequest)

	statusCode := responseRecorder.Result().StatusCode
	httpBodyReader := responseRecorder.Result().Body
	defer func() {
		if closeErr := httpBodyReader.Close(); closeErr != nil {
			t.Fatal(closeErr)
		}
	}()

	bodyData := make([]byte, 0)
	httpBodyBuffer := bytes.NewBuffer(bodyData)
	_, err := io.Copy(httpBodyBuffer, httpBodyReader)
	if err != nil {
		t.Fatal(err)
	}
	return &testResponse{
		statusCode: statusCode,
		bodyBuffer: httpBodyBuffer,
	}
}

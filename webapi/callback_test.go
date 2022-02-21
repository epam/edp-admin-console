package webapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2"

	"edp-admin-console/internal/config"
	applog "edp-admin-console/service/logger"
)

const TokenURL = "http://domain"
const authID = "id"

type CallbackSuite struct {
	suite.Suite
	TestServer  *httptest.Server
	Handler     *HandlerEnv
	HandlerAuth *HandlerAuth
}

func TestCallbackSuite(t *testing.T) {

	conf := &config.AppConfig{
		BasePath:   "/",
		AuthEnable: true,
	}
	authController := &config.AuthController{
		Oauth2Service: oauth2.Config{
			Endpoint: oauth2.Endpoint{
				TokenURL:  TokenURL,
				AuthStyle: 0,
			},
			Scopes: []string{},
		},
	}

	h := NewHandlerEnv(WithConfig(conf))
	authHandler := HandlerAuthWithOption(WithAuthController(authController))
	logger := applog.GetLogger()

	router := V2APIRouter(h, authHandler, logger)

	s := &CallbackSuite{
		Handler:     h,
		HandlerAuth: authHandler,
	}
	s.TestServer = httptest.NewServer(router)
	suite.Run(t, s)
}

func (s *CallbackSuite) TestCallbackNoCookie() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.TestServer.URL)
	httpExpect.
		GET("/v2/auth/callback").
		Expect().
		Status(http.StatusInternalServerError).
		ContentType("text/plain").
		Body().
		Contains("cant find auth session id")
}

func (s *CallbackSuite) TestCallbackNoSessionState() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.TestServer.URL)

	httpExpect.
		Request(http.MethodGet, "/v2/auth/callback").
		WithCookie(AuthSessionIDName, authID).
		Expect().
		Status(http.StatusInternalServerError).
		ContentType("text/plain").
		Body().
		Contains("cant find session state")
}

func (s *CallbackSuite) TestCallbackOK() {
	accessToken := "123"
	type tokenJSON struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int32  `json:"expires_in"`
	}
	token := tokenJSON{
		AccessToken: accessToken,
		ExpiresIn:   -1,
	}

	t := s.T()
	authState := "state"
	code := "code"
	s.HandlerAuth.StateMap[authID] = authState
	s.HandlerAuth.UrlMap[authID] = s.TestServer.URL
	httpExpect := httpexpect.New(t, s.TestServer.URL)

	httpmock.Activate()
	jsonResponder, err := httpmock.NewJsonResponder(http.StatusOK, token)
	if err != nil {
		t.Fatal(err)
	}
	httpmock.RegisterResponder(http.MethodPost, TokenURL, jsonResponder)
	client := resty.New()
	r := httpExpect.
		Request(http.MethodGet, "/v2/auth/callback").
		WithCookie(AuthSessionIDName, authID).
		WithQuery("state", authState).
		WithQuery("code", code).
		WithClient(client.GetClient()). //the use of a non-default client is due to the fact that the default client is mocked and used for auth mock
		WithRedirectPolicy(httpexpect.DontFollowRedirects)

	r.Expect().
		Status(http.StatusFound).Header("location").Equal(s.TestServer.URL)
	httpmock.DeactivateAndReset()
}

func TestHandlerEnv_getRedirectPath(t *testing.T) {
	sessionID := "id"
	expectedPath := "path"
	urlMap := map[string]string{sessionID: expectedPath}
	handler := HandlerAuth{
		UrlMap: urlMap,
	}
	redirectPath := handler.getRedirectPath(sessionID)
	assert.Equal(t, expectedPath, redirectPath)
}

func TestHandlerEnv_getRedirectPathDefaultPath(t *testing.T) {
	sessionID := "id"
	handler := HandlerAuth{
		BasePath: "",
	}
	redirectPath := handler.getRedirectPath(sessionID)
	assert.Equal(t, "/admin/edp/overview", redirectPath)
}

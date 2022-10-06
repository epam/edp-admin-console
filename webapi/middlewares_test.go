package webapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"edp-admin-console/internal/applog"
	"edp-admin-console/internal/config"
)

const TokenURL = "http://domain"

func testContextWithLogger(t *testing.T) context.Context {
	t.Helper()
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	return applog.ContextWithLogger(context.Background(), logger)
}

func TestGetCookieByName(t *testing.T) {
	authSessionId := "id"
	r := &http.Request{
		Header: make(map[string][]string),
	}
	r = r.WithContext(testContextWithLogger(t))
	cookie := &http.Cookie{
		Name:    AuthSessionIDName,
		Value:   authSessionId,
		Expires: time.Now().Add(5 * time.Minute),
		Path:    "/",
	}
	r.AddCookie(cookie)
	requestId, ok := GetCookieByName(r, AuthSessionIDName)
	assert.Equal(t, authSessionId, requestId)
	assert.True(t, ok)
}

func TestGetCookieByNameNoCookie(t *testing.T) {
	r := &http.Request{}
	r = r.WithContext(testContextWithLogger(t))
	requestId, ok := GetCookieByName(r, AuthSessionIDName)
	assert.Empty(t, requestId)
	assert.False(t, ok)
}

type testHandler struct {
	ExpectedUser AuthorisedUser
	T            *testing.T
}

func (h testHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	user, ok := ctx.Value(AuthorisedUserKey{}).(AuthorisedUser)

	assert.True(h.T, ok)
	assert.Equal(h.T, h.ExpectedUser.UserName(), user.UserName())
	assert.Equal(h.T, h.ExpectedUser.IsAdmin(), user.IsAdmin())
	assert.Equal(h.T, h.ExpectedUser.IsDeveloper(), user.IsDeveloper())

	writer.WriteHeader(http.StatusOK)
}

type Verifier struct {
	Payload []byte
}

func (v Verifier) VerifySignature(_ context.Context, _ string) (payload []byte, err error) {
	return v.Payload, nil
}

func TestMiddleware(t *testing.T) {
	ctx := context.Background()

	testName := "name"
	authCode := "code"
	sessionId := "sessionId"
	adminRole := "admin"
	devRole := "dev"
	ttl := 1000
	roles := []string{adminRole}
	//generating jwt token with AuthorisedUser Claims
	expectedTokenClaim := TokenClaim{
		Name: testName,
		RealmAccess: RealmAccess{
			Roles: roles,
		},
	}
	claims := jwt.MapClaims{
		"name":         expectedTokenClaim.Name,
		"realm_access": expectedTokenClaim.RealmAccess,
	}
	payload, err := json.Marshal(claims)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signingString, err := jwtToken.SignedString([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}

	type tokenJSON struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int32  `json:"expires_in"`
	}

	token := tokenJSON{
		AccessToken: signingString,
		ExpiresIn:   int32(ttl),
	}

	oidcConfig := oidc.Config{
		SupportedSigningAlgs: []string{"HS256"},
		SkipClientIDCheck:    true,
		SkipExpiryCheck:      true,
		SkipIssuerCheck:      true,
	}

	verifier := Verifier{Payload: payload}
	authController := &config.AuthController{
		AdminRoleName: adminRole,
		DevRoleName:   devRole,
		Oauth2Service: oauth2.Config{
			Endpoint: oauth2.Endpoint{
				TokenURL:  TokenURL,
				AuthStyle: 0,
			},
			Scopes: []string{},
		},
		Verifier: oidc.NewVerifier(TokenURL, verifier, &oidcConfig),
	}

	httpmock.Activate()
	jsonResponder, err := httpmock.NewJsonResponder(http.StatusOK, token)
	if err != nil {
		t.Fatal(err)
	}
	httpmock.RegisterResponder(http.MethodPost, TokenURL, jsonResponder)
	tokenAuth, err := authController.Oauth2Service.Exchange(ctx, authCode)
	if err != nil {
		t.Fatal(err)
	}

	ts := authController.Oauth2Service.TokenSource(ctx, tokenAuth)
	tokenMap := map[string]oauth2.TokenSource{
		sessionId: ts,
	}

	mw := WithAuthZ(tokenMap, authController)
	expectedUser := NewAuthorisedUser(expectedTokenClaim, ConfigRoles{DevRole: devRole, AdminRole: adminRole})
	initHandler := testHandler{ExpectedUser: expectedUser, T: t}
	handler := mw(initHandler)
	w := httptest.NewRecorder()
	var data []byte
	reader := bytes.NewReader(data)
	req := httptest.NewRequest(http.MethodGet, "/", reader)
	cookie := &http.Cookie{
		Name:    SessionIDName,
		Value:   sessionId,
		Expires: time.Now().Add(5 * time.Minute),
		Path:    "/",
	}

	req.AddCookie(cookie)
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	httpmock.DeactivateAndReset()

}

type stubOKHandler struct {
}

func (s *stubOKHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
}

func TestWithLogRequestBoundaries(t *testing.T) {
	mw := WithLogRequestBoundaries()
	next := new(stubOKHandler)

	w := httptest.NewRecorder()
	var data []byte
	reader := bytes.NewReader(data)
	r := httptest.NewRequest(http.MethodGet, "/end/point", reader)

	mw(next).ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

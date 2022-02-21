package webapi

import (
	"context"
	"net/http"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func testContextWithLogger(t *testing.T) context.Context {
	t.Helper()
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	return ContextWithLogger(context.Background(), logger)
}

func TestSessionIdByRequest(t *testing.T) {
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
	requestId, ok := SessionIdByRequest(r, AuthSessionIDName)
	assert.Equal(t, authSessionId, requestId)
	assert.True(t, ok)
}

func TestSessionIdByRequestNoCookie(t *testing.T) {
	r := &http.Request{}
	r = r.WithContext(testContextWithLogger(t))
	requestId, ok := SessionIdByRequest(r, AuthSessionIDName)
	assert.Empty(t, requestId)
	assert.False(t, ok)
}

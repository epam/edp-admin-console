package filters

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
	"github.com/stretchr/testify/assert"
)

func TestAccessControlEnv_RoleAccessControlRestFilter(t *testing.T) {
	beegoCookieStorage := &session.CookieSessionStore{}
	err := beegoCookieStorage.Flush()
	if err != nil {
		t.Fatal(err)
	}

	err = beegoCookieStorage.Set("realm_roles", false)
	if err != nil {
		t.Fatal(err)
	}
	beegoInput := context.NewInput()
	beegoInput.CruSession = beegoCookieStorage

	beegoCtx := context.NewContext()
	beegoCtx.Input = beegoInput

	httpRersponseWriter := httptest.NewRecorder()
	beegoCtx.ResponseWriter = &context.Response{
		ResponseWriter: httpRersponseWriter,
		Started:        false,
		Status:         0,
		Elapsed:        0,
	}

	accessControlEnv := &AccessControlEnv{Permissions: map[string][]string{}}
	accessControlEnv.RoleAccessControlRestFilter(beegoCtx)

	expectedCode := http.StatusForbidden
	assert.Equal(t, expectedCode, httpRersponseWriter.Code)
}

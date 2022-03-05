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

	"edp-admin-console/internal/config"
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
	conf := &config.AppConfig{
		BasePath:   "/",
		AuthEnable: false,
		EDPVersion: "v1",
	}
	h := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMap(CreateCommonFuncMap()), WithConfig(conf))
	logger := applog.GetLogger()
	authHandler := HandlerAuthWithOption()

	router := V2APIRouter(h, authHandler, logger)

	s := &IndexSuite{
		Handler: h,
	}
	s.TestServer = httptest.NewServer(router)
	suite.Run(t, s)
}

func CustomUser(userName string, isAdmin, isDeveloper bool) AuthorisedUser {
	adminRole := "admin"
	devRole := "dev"
	var roles []string
	if isAdmin {
		roles = append(roles, adminRole)
	}
	if isDeveloper {
		roles = append(roles, devRole)
	}
	return AuthorisedUser{
		TokenClaim: TokenClaim{
			Name: userName,
			RealmAccess: RealmAccess{
				Roles: roles,
			},
		},
		ConfigRoles: ConfigRoles{
			DevRole:   adminRole,
			AdminRole: devRole,
		},
	}
}

func (s *IndexSuite) TestIndex() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.TestServer.URL)
	httpExpect.
		GET("/v2").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html")
}

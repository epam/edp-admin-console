package webapi

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/gavv/httpexpect/v2"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	edpComponentAPI "github.com/epam/edp-component-operator/pkg/apis/v1/v1"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

func TestHandlerEnv_GetOverviewPage(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(edpComponentAPI.SchemeGroupVersion, &edpComponentAPI.EDPComponent{}, &edpComponentAPI.EDPComponentList{})
	componentName := "EDP component"
	componentURL := "https://example.com"
	icon := "icon"
	componentType := "edp"
	component := edpComponentAPI.EDPComponent{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      componentName,
			Namespace: namespace,
		},
		Spec: edpComponentAPI.EDPComponentSpec{
			Url:     componentURL,
			Visible: true,
			Icon:    icon,
			Type:    componentType,
		},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&component).Build()
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
	server := httptest.NewServer(router)

	httpExpect := httpexpect.New(t, server.URL)
	httpExpect.
		GET("/v2/admin/edp/overview").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html")
}

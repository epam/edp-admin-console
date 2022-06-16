package webapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	cdPipelineAPI "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	edpComponentAPI "github.com/epam/edp-component-operator/pkg/apis/v1/v1alpha1"
	"github.com/gavv/httpexpect/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

func TestHandlerEnv_GetPipelineOverviewPage(t *testing.T) {
	scheme := runtime.NewScheme()
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	err = cdPipelineAPI.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	err = edpComponentAPI.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}

	pipeName := "cd1"
	deploymentType := "container"
	initialISName := []string{"app1master"}

	initialCDPipelineCR := &cdPipelineAPI.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pipeName,
			Namespace: namespace,
		},
		Spec: cdPipelineAPI.CDPipelineSpec{
			DeploymentType:     deploymentType,
			InputDockerStreams: initialISName,
		},
	}

	imageSteamCR := &codeBaseApi.CodebaseImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      initialISName[0],
			Namespace: namespace,
		},
		Spec: codeBaseApi.CodebaseImageStreamSpec{},
	}

	jenkinsUrl := "https://domain"
	jenkinsComponentCR := &edpComponentAPI.EDPComponent{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jenkins",
			Namespace: namespace,
		},
		Spec: edpComponentAPI.EDPComponentSpec{
			Url: jenkinsUrl,
		},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initialCDPipelineCR, jenkinsComponentCR, imageSteamCR).Build()
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
		BasePath:    "/",
		AuthEnable:  false,
		XSRFEnabled: false,
	}
	h := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMap(CreateCommonFuncMap()), WithConfig(conf))
	logger := applog.GetLogger()
	authHandler := HandlerAuthWithOption()
	router := V2APIRouter(h, authHandler, logger)
	testServer := httptest.NewServer(router)
	httpExpect := httpexpect.New(t, testServer.URL)
	httpExpect.
		GET(fmt.Sprintf("/v2/admin/edp/cd-pipeline/%s/overview", pipeName)).
		Expect().
		Status(http.StatusOK).
		ContentType("text/html")

}

package webapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	cdPipelineAPI "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

func TestHandlerEnv_UpdateCDPipeline(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	err = cdPipelineAPI.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}

	pipeName := "cd1"
	deploymentType := "container"
	initialISName := []string{"app1master", "app2master"}

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

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initialCDPipelineCR).Build()
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

	newApps := []string{"application1"}
	newAppsToPromote := []string{newApps[0]}
	newISName := []string{"app1dev"}
	newStages := []string{"newStage"}
	stepsForNewStage := []string{"step-new-stage1"}
	stepAutotest := "null"
	stepBranch := "null"
	manualQualityGate := "manual"
	stageDescriptions := []string{"stage1 desc"}
	triggerType := "Auto"
	defaultPipelineLibraryName := "default"
	defaultJobProvisioning := "default"

	testServer := httptest.NewServer(router)
	httpExpect := httpexpect.New(t, testServer.URL)
	req := httpExpect.POST(fmt.Sprintf("/v2/admin/edp/cd-pipeline/%s/update", pipeName)).
		WithForm(map[string][]string{"app": newApps}).
		WithFormField("app", newApps[0]).
		WithFormField("pipelineName", pipeName).
		WithFormField("deploymentType", deploymentType).
		WithFormField(strings.Join([]string{newApps[0], "promote"}, "-"), "true").
		WithFormField(newApps[0], newISName[0]).
		WithFormField("stageName", newStages[0]).
		WithFormField(strings.Join([]string{newStages[0], "stageDesc"}, "-"), stageDescriptions[0]).
		WithFormField(strings.Join([]string{newStages[0], "triggerType"}, "-"), triggerType).
		WithFormField(strings.Join([]string{newStages[0], "pipelineLibraryName"}, "-"), defaultPipelineLibraryName).
		WithFormField(strings.Join([]string{newStages[0], "jobProvisioning"}, "-"), defaultJobProvisioning).
		WithFormField(strings.Join([]string{newStages[0], "stageStepName"}, "-"), stepsForNewStage[0]).
		WithFormField(strings.Join([]string{newStages[0], stepsForNewStage[0], "stageAutotests"}, "-"), stepAutotest).
		WithFormField(strings.Join([]string{newStages[0], stepsForNewStage[0], "stageBranch"}, "-"), stepBranch).
		WithFormField(strings.Join([]string{newStages[0], stepsForNewStage[0], "stageQualityGateType"}, "-"), manualQualityGate).
		WithRedirectPolicy(httpexpect.DontFollowRedirects)
	expectedURL := fmt.Sprintf("%s/v2/admin/edp/cd-pipeline/overview#cdPipelineEditSuccessModal", h.Config.BasePath)

	req.Expect().
		Status(http.StatusFound).
		Header("location").
		Equal(expectedURL)

	pipeline, err := h.NamespacedClient.GetCDPipeline(ctx, pipeName)
	assert.NoError(t, err)
	assert.Equal(t, newApps, pipeline.Spec.Applications)
	assert.Equal(t, newAppsToPromote, pipeline.Spec.ApplicationsToPromote)
	assert.Equal(t, newISName, pipeline.Spec.InputDockerStreams)

	stage1, err := h.NamespacedClient.GetCDStage(ctx, strings.Join([]string{pipeName, newStages[0]}, "-"))
	assert.NoError(t, err)
	assert.Equal(t, newStages[0], stage1.Spec.Name)
	assert.Equal(t, triggerType, stage1.Spec.TriggerType)
	assert.Equal(t, stageDescriptions[0], stage1.Spec.Description)
	assert.Equal(t, defaultJobProvisioning, stage1.Spec.JobProvisioning)

	assert.Equal(t, manualQualityGate, stage1.Spec.QualityGates[0].QualityGateType)
	assert.Equal(t, stepsForNewStage[0], stage1.Spec.QualityGates[0].StepName)
	assert.Nil(t, stage1.Spec.QualityGates[0].AutotestName)
	assert.Nil(t, stage1.Spec.QualityGates[0].BranchName)

	assert.Equal(t, defaultPipelineLibraryName, stage1.Spec.Source.Type)

}

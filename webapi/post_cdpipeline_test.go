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
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	applog "edp-admin-console/service/logger"
)

func TestHandlerEnv_CreateCDPipeline(t *testing.T) {
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
	client := fake.NewClientBuilder().WithScheme(scheme).Build()
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

	apps := []string{"application1", "application2"}
	appsToPromote := []string{apps[0]}
	pipeName := "cd1"
	deploymentType := "container"
	ISName := []string{"app1master", "app2master"}
	stages := []string{"stage1", "stage2"}
	stepsForFirstStage := []string{"step1stage1", "step2stage1"}
	stepsForSecondStage := "step1stage2"
	stepAutotest := "null"
	stepBranch := "null"
	manualQualityGate := "manual"
	autoQualityGate := "auto"
	stageDescriptions := []string{"stage1 desc", "stage2 desc"}
	triggerType := "Auto"
	defaultPipelineLibraryName := "default"
	defaultJobProvisioning := "default"

	testServer := httptest.NewServer(router)
	httpExpect := httpexpect.New(t, testServer.URL)
	req := httpExpect.POST(fmt.Sprintf("/v2/admin/edp/cd-pipeline")).
		WithForm(map[string][]string{"app": apps}).
		WithFormField("app", apps[0]).WithFormField("app", apps[1]).
		WithFormField("pipelineName", pipeName).
		WithFormField("deploymentType", deploymentType).
		WithFormField(strings.Join([]string{apps[0], "promote"}, "-"), "true").
		WithFormField(strings.Join([]string{apps[1], "promote"}, "-"), "false").
		WithFormField(apps[0], ISName[0]).
		WithFormField(apps[1], ISName[1]).
		WithFormField("stageName", stages[0]).
		WithFormField("stageName", stages[1]).
		WithFormField(strings.Join([]string{stages[0], "stageDesc"}, "-"), stageDescriptions[0]).
		WithFormField(strings.Join([]string{stages[1], "stageDesc"}, "-"), stageDescriptions[1]).
		WithFormField(strings.Join([]string{stages[0], "triggerType"}, "-"), triggerType).
		WithFormField(strings.Join([]string{stages[1], "triggerType"}, "-"), triggerType).
		WithFormField(strings.Join([]string{stages[0], "pipelineLibraryName"}, "-"), defaultPipelineLibraryName).
		WithFormField(strings.Join([]string{stages[1], "pipelineLibraryName"}, "-"), defaultPipelineLibraryName).
		WithFormField(strings.Join([]string{stages[0], "jobProvisioning"}, "-"), defaultJobProvisioning).
		WithFormField(strings.Join([]string{stages[1], "jobProvisioning"}, "-"), defaultJobProvisioning).
		WithFormField(strings.Join([]string{stages[0], "stageStepName"}, "-"), stepsForFirstStage[0]).
		WithFormField(strings.Join([]string{stages[0], "stageStepName"}, "-"), stepsForFirstStage[1]).
		WithFormField(strings.Join([]string{stages[0], stepsForFirstStage[0], "stageAutotests"}, "-"), stepAutotest).
		WithFormField(strings.Join([]string{stages[0], stepsForFirstStage[0], "stageBranch"}, "-"), stepBranch).
		WithFormField(strings.Join([]string{stages[0], stepsForFirstStage[0], "stageQualityGateType"}, "-"), manualQualityGate).
		WithFormField(strings.Join([]string{stages[0], stepsForFirstStage[1], "stageAutotests"}, "-"), stepAutotest).
		WithFormField(strings.Join([]string{stages[0], stepsForFirstStage[1], "stageBranch"}, "-"), stepBranch).
		WithFormField(strings.Join([]string{stages[0], stepsForFirstStage[1], "stageQualityGateType"}, "-"), manualQualityGate).
		WithFormField(strings.Join([]string{stages[1], "stageStepName"}, "-"), stepsForSecondStage).
		WithFormField(strings.Join([]string{stages[1], stepsForSecondStage, "stageAutotests"}, "-"), stepAutotest).
		WithFormField(strings.Join([]string{stages[1], stepsForSecondStage, "stageBranch"}, "-"), stepBranch).
		WithFormField(strings.Join([]string{stages[1], stepsForSecondStage, "stageQualityGateType"}, "-"), autoQualityGate).
		WithRedirectPolicy(httpexpect.DontFollowRedirects)
	expectedURL := fmt.Sprintf("%s/v2/admin/edp/cd-pipeline/overview?%s=%s#cdPipelineSuccessModal", h.Config.BasePath, paramWaitingForCdPipeline, pipeName)

	req.Expect().
		Status(http.StatusFound).
		Header("location").
		Equal(expectedURL)

	pipeline, err := h.NamespacedClient.GetCDPipeline(ctx, pipeName)
	assert.NoError(t, err)
	assert.Equal(t, pipeName, pipeline.Spec.Name)
	assert.Equal(t, apps, pipeline.Spec.Applications)
	assert.Equal(t, appsToPromote, pipeline.Spec.ApplicationsToPromote)
	assert.Equal(t, ISName, pipeline.Spec.InputDockerStreams)
	assert.Equal(t, deploymentType, pipeline.Spec.DeploymentType)

	stage1, err := h.NamespacedClient.GetCDStage(ctx, strings.Join([]string{pipeName, stages[0]}, "-"))
	assert.NoError(t, err)
	assert.Equal(t, stages[0], stage1.Spec.Name)
	assert.Equal(t, triggerType, stage1.Spec.TriggerType)
	assert.Equal(t, stageDescriptions[0], stage1.Spec.Description)
	assert.Equal(t, defaultJobProvisioning, stage1.Spec.JobProvisioning)

	assert.Equal(t, manualQualityGate, stage1.Spec.QualityGates[0].QualityGateType)
	assert.Equal(t, stepsForFirstStage[0], stage1.Spec.QualityGates[0].StepName)
	assert.Nil(t, stage1.Spec.QualityGates[0].AutotestName)
	assert.Nil(t, stage1.Spec.QualityGates[0].BranchName)

	assert.Equal(t, manualQualityGate, stage1.Spec.QualityGates[1].QualityGateType)
	assert.Equal(t, stepsForFirstStage[1], stage1.Spec.QualityGates[1].StepName)
	assert.Nil(t, stage1.Spec.QualityGates[1].AutotestName)
	assert.Nil(t, stage1.Spec.QualityGates[1].BranchName)

	assert.Equal(t, defaultPipelineLibraryName, stage1.Spec.Source.Type)

	stage2, err := h.NamespacedClient.GetCDStage(ctx, strings.Join([]string{pipeName, stages[1]}, "-"))
	assert.NoError(t, err)
	assert.Equal(t, stages[1], stage2.Spec.Name)
	assert.Equal(t, triggerType, stage2.Spec.TriggerType)
	assert.Equal(t, stageDescriptions[1], stage2.Spec.Description)
	assert.Equal(t, defaultJobProvisioning, stage2.Spec.JobProvisioning)

	assert.Equal(t, autoQualityGate, stage2.Spec.QualityGates[0].QualityGateType)
	assert.Equal(t, stepsForSecondStage, stage2.Spec.QualityGates[0].StepName)
	assert.Nil(t, stage2.Spec.QualityGates[0].AutotestName)
	assert.Nil(t, stage2.Spec.QualityGates[0].BranchName)

	assert.Equal(t, defaultPipelineLibraryName, stage2.Spec.Source.Type)
}

package pipelinestage

import (
	"context"
	"testing"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/k8s"
)

const (
	namespace      = "ns"
	cdPipelineName = "cdPipelineName"
	branchName     = "branch"
	stageName      = "stage"
)

func TestCreateApplicationStage_Err(t *testing.T) {
	inputIS := []string{"input1", "input2"}
	outputIS := []string{"output1"}
	appNames := []string{"app1"}
	client := fake.NewClientBuilder().Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	stages, err := BuildApplicationStages(ctx, namespacedClient, inputIS, outputIS, appNames)
	assert.Error(t, err)
	assert.Nil(t, stages)
}

func TestCreateApplicationStage(t *testing.T) {
	inputIS := []string{"input1"}
	outputIS := []string{"output1"}
	appNames := []string{"app1"}

	scheme := runtime.NewScheme()
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	cbBranch := codeBaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      inputIS[0],
			Namespace: namespace,
		},
		Spec: codeBaseApi.CodebaseBranchSpec{
			BranchName: branchName,
		},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&cbBranch).Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	expectedStages := []ApplicationStage{
		{
			Name:       appNames[0],
			InputIs:    inputIS[0],
			OutputIs:   outputIS[0],
			BranchName: branchName,
		},
	}

	stages, err := BuildApplicationStages(ctx, namespacedClient, inputIS, outputIS, appNames)
	assert.NoError(t, err)
	assert.Equal(t, expectedStages, stages)
}

func TestCdPipelineNameByCRName_Err(t *testing.T) {
	client := fake.NewClientBuilder().Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	name, err := CdPipelineAppNamesByCRName(ctx, namespacedClient, cdPipelineName)
	assert.Error(t, err)
	assert.True(t, runtime.IsNotRegisteredError(err))
	assert.Empty(t, name)
}

func TestCdPipelineNameByCRName(t *testing.T) {
	appNames := []string{"app1"}
	scheme := runtime.NewScheme()
	err := cdPipeApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	cdPipe := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cdPipelineName,
			Namespace: namespace,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			ApplicationsToPromote: appNames,
		},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&cdPipe).Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	name, err := CdPipelineAppNamesByCRName(ctx, namespacedClient, cdPipelineName)
	assert.NoError(t, err)
	assert.Equal(t, appNames, name)
}

func TestStageViewByCRName(t *testing.T) {
	description := "test description"
	triggerType := "manual"
	order := 0
	jobProvisioning := "default"
	qgType := "manual"
	stepName := "dev"
	autotestName := "autotest"
	branchName := branchName

	scheme := runtime.NewScheme()
	err := cdPipeApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	err = codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}

	stage := cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      stageName,
			Namespace: namespace,
		},
		Spec: cdPipeApi.StageSpec{
			Name:        stageName,
			CdPipeline:  cdPipelineName,
			Description: description,
			TriggerType: triggerType,
			Order:       order,
			QualityGates: []cdPipeApi.QualityGate{
				{
					QualityGateType: qgType,
					StepName:        stepName,
					AutotestName:    &autotestName,
					BranchName:      &branchName,
				},
			},
			JobProvisioning: jobProvisioning,
		},
	}

	branch := codeBaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      branchName,
			Namespace: namespace,
		},
		Spec: codeBaseApi.CodebaseBranchSpec{
			BranchName: branchName,
		},
	}
	codebase := codeBaseApi.Codebase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      autotestName,
			Namespace: namespace,
		},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&stage, &branch, &codebase).Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	emptyString := ""
	cdStageID := 0
	expectedStageView := StageMainData{
		Description:     description,
		TriggerType:     triggerType,
		Order:           "0",
		JobProvisioning: jobProvisioning,
		QualityGates: []QualityGate{
			{
				Id:              0,
				QualityGateType: qgType,
				StepName:        stepName,
				Autotest: &Codebase{
					Name:      autotestName,
					GitServer: &emptyString,
				},
				Branch:    &CodebaseBranch{Name: branchName},
				CdStageId: &cdStageID,
			},
		},
	}

	name, err := StageViewByCRName(ctx, namespacedClient, stageName)
	assert.NoError(t, err)
	assert.Equal(t, &expectedStageView, name)

}

func TestStageViewByCRName_Err(t *testing.T) {
	client := fake.NewClientBuilder().Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	name, err := StageViewByCRName(ctx, namespacedClient, stageName)
	assert.Error(t, err)
	assert.Nil(t, name)
}

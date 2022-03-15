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
	client := fake.NewClientBuilder().Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	stages, err := BuildApplicationStages(ctx, namespacedClient, inputIS, outputIS)
	assert.Error(t, err)
	assert.Nil(t, stages)
}

func TestCreateApplicationStage(t *testing.T) {
	inputIS := []string{"input1"}
	outputIS := []string{"output1"}
	appName := "app1"

	scheme := runtime.NewScheme()
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	inputISCR := codeBaseApi.CodebaseImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      inputIS[0],
			Namespace: namespace,
		},
		Spec: codeBaseApi.CodebaseImageStreamSpec{
			Codebase: appName,
		},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&inputISCR).Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	expectedStages := []ApplicationStage{
		{
			Name:     appName,
			InputIs:  inputIS[0],
			OutputIs: outputIS[0],
		},
	}

	stages, err := BuildApplicationStages(ctx, namespacedClient, inputIS, outputIS)
	assert.NoError(t, err)
	assert.Equal(t, expectedStages, stages)
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

func TestStageListByPipelineName_OK(t *testing.T) {
	cdPipelineName_1 := "test_cd_pipeline_1"
	cdPipelineName_2 := "test_cd_pipeline_2"
	stageName_1 := "test_stage_1"
	stageName_2 := "test_stage_2"
	stageName_3 := "test_stage_3"
	order_1 := 0
	order_2 := 1
	order_3 := 2

	scheme := runtime.NewScheme()
	err := cdPipeApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	err = codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}

	stage_1 := &cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      stageName_1,
			Namespace: namespace,
		},
		Spec: cdPipeApi.StageSpec{
			Name:       stageName_1,
			CdPipeline: cdPipelineName_1,
			Order:      order_1,
		},
	}
	stage_2 := &cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      stageName_2,
			Namespace: namespace,
		},
		Spec: cdPipeApi.StageSpec{
			Name:       stageName_2,
			CdPipeline: cdPipelineName_1,
			Order:      order_2,
		},
	}
	stage_3 := &cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      stageName_3,
			Namespace: namespace,
		},
		Spec: cdPipeApi.StageSpec{
			Name:       stageName_3,
			CdPipeline: cdPipelineName_2,
			Order:      order_3,
		},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(stage_1, stage_2, stage_3).Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	expectedStageList := []cdPipeApi.Stage{
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       "",
				APIVersion: "",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            stageName_1,
				Namespace:       namespace,
				ResourceVersion: "999",
			},
			Spec: cdPipeApi.StageSpec{
				Name:       stageName_1,
				CdPipeline: cdPipelineName_1,
				Order:      order_1,
			},
		},
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       "",
				APIVersion: "",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            stageName_2,
				Namespace:       namespace,
				ResourceVersion: "999",
			},
			Spec: cdPipeApi.StageSpec{
				Name:       stageName_2,
				CdPipeline: cdPipelineName_1,
				Order:      order_2,
			},
		},
	}

	gotStageList, err := StageListByPipelineName(ctx, namespacedClient, cdPipelineName_1)
	assert.NoError(t, err)
	assert.Equal(t, expectedStageList, gotStageList)
}

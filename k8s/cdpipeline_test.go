package k8s

import (
	"context"
	"testing"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestRuntimeNamespacedClient_CreateCDPipelineBySpec_AlreadyExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.CDPipeline{})

	initCR := createCDPipelineCR(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initCR).Build()
	nameCD := "name"
	spec := cdPipeApi.CDPipelineSpec{Name: nameCD}

	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	err = k8sClient.CreateCDPipelineBySpec(ctx, name, spec)
	assert.Error(t, err)
	assert.True(t, k8serrors.IsAlreadyExists(err))

	cdPipe, err := k8sClient.GetCDPipeline(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, name, cdPipe.Name)
}

func TestRuntimeNamespacedClient_CreateCDPipelineBySpec(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.CDPipeline{})

	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	nameCD := "name"
	spec := cdPipeApi.CDPipelineSpec{Name: nameCD}

	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	err = k8sClient.CreateCDPipelineBySpec(ctx, name, spec)
	assert.NoError(t, err)

	cdPipe, err := k8sClient.GetCDPipeline(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, nameCD, cdPipe.Spec.Name)
}

func TestRuntimeNamespacedClient_CreateCDStageBySpec_AlreadyExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{})

	initCR := createCDStageCR(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initCR).Build()
	cdName := "name"
	spec := cdPipeApi.StageSpec{
		CdPipeline: cdName,
	}

	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	err = k8sClient.CreateCDStageBySpec(ctx, name, spec)
	assert.Error(t, err)
	assert.True(t, k8serrors.IsAlreadyExists(err))

	stage, err := k8sClient.GetCDStage(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, name, stage.Name)
}

func TestRuntimeNamespacedClient_CreateCDStageBySpec(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{})

	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	cdName := "name"
	spec := cdPipeApi.StageSpec{
		CdPipeline: cdName,
	}

	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}
	err = k8sClient.CreateCDStageBySpec(ctx, name, spec)
	assert.NoError(t, err)

	stage, err := k8sClient.GetCDStage(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, cdName, stage.Spec.CdPipeline)
}

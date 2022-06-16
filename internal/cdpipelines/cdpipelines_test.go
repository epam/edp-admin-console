package cdpipelines

import (
	"context"
	"testing"

	"edp-admin-console/k8s"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type CDPipelineCROption func(cdPipe *cdPipeApi.CDPipeline)

func createCDPipelineCRWithOptions(opts ...CDPipelineCROption) *cdPipeApi.CDPipeline {
	cdPipeCR := new(cdPipeApi.CDPipeline)
	for i := range opts {
		opts[i](cdPipeCR)
	}
	return cdPipeCR
}

func WithCDPipelineNamespace(namespace string) CDPipelineCROption {
	return func(cdPipe *cdPipeApi.CDPipeline) {
		cdPipe.Namespace = namespace
	}
}

func WithCDPipelineName(name string) CDPipelineCROption {
	return func(cdPipe *cdPipeApi.CDPipeline) {
		cdPipe.Name = name
	}
}

func TestByNameIFExists_OK(t *testing.T) {
	namespace := "test_namespace"
	crCDPipelineName_1 := "test_cd_pipeline_1"
	cdPipelineCR_1 := createCDPipelineCRWithOptions(
		WithCDPipelineNamespace(namespace),
		WithCDPipelineName(crCDPipelineName_1),
	)

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion,
		&cdPipeApi.CDPipeline{},
	)
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cdPipelineCR_1).Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	ctx := context.Background()
	k8sCDPipeline, err := ByNameIFExists(ctx, namespacedClient, crCDPipelineName_1)
	assert.NoError(t, err)

	expectedCDPipeline := cdPipelineCR_1.DeepCopy()
	expectedCDPipeline.TypeMeta.Kind = "CDPipeline"
	expectedCDPipeline.TypeMeta.APIVersion = "apps/v1"
	assert.Equal(t, expectedCDPipeline, k8sCDPipeline)
}

func TestByNameIFExists_NotFound(t *testing.T) {
	namespace := "test_namespace"
	crCDPipelineName_1 := "test_not_found_cd_pipeline_1"

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion,
		&cdPipeApi.CDPipeline{},
	)
	err := codeBaseApi.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}

	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	namespacedClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	ctx := context.Background()
	k8sCDPipeline, err := ByNameIFExists(ctx, namespacedClient, crCDPipelineName_1)
	assert.NoError(t, err)
	var expectedCDPipeline *cdPipeApi.CDPipeline
	assert.Equal(t, expectedCDPipeline, k8sCDPipeline)
}

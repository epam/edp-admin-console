package imagestream

import (
	"context"
	"strings"
	"testing"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"edp-admin-console/k8s"
)

const (
	name              = "name"
	ns                = "ns"
	firstImage        = "firstImage"
	secondImage       = "secondImage"
	cdPipelineName    = "CDPipeline"
	zeroOrder         = 0
	nonZeroOrder      = 1
	previousStageName = "previous-stage-name"
	appName           = "app-name"
)

func createStageCR(order int, cdPipeName string) cdPipeApi.Stage {
	return cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createStageCrName(cdPipelineName, name),
			Namespace: ns,
		},
		Spec: cdPipeApi.StageSpec{
			Order:      order,
			CdPipeline: cdPipeName,
		},
	}
}

func TestGetImageStreamFromStage_BadClient(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{})

	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	k8sClient := k8s.RuntimeNamespacedClient{Client: client}

	stage, err := GetInputISForStage(ctx, &k8sClient, name, cdPipelineName)
	assert.Error(t, err)
	assert.Nil(t, stage)
	assert.True(t, k8s.AsEmptyNamespaceErr(err))
}

func TestGetImageStreamFromStage_GetStageErr(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{})

	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)

	stage, err := GetInputISForStage(ctx, k8sClient, name, cdPipelineName)
	assert.Error(t, err)
	assert.Nil(t, stage)
	assert.True(t, k8serrors.IsNotFound(err))
}

func TestGetImageStreamFromStage_NonZeroOrder(t *testing.T) {
	ctx := context.Background()
	stageCR := createStageCR(nonZeroOrder, "")
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{})

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&stageCR).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)

	stage, err := GetInputISForStage(ctx, k8sClient, name, cdPipelineName)
	assert.Error(t, err)
	assert.Nil(t, stage)
	assert.True(t, strings.Contains(err.Error(), "Spec.CdPipeline is empty in Stage CR named"))
}

func TestGetImageStreamFromStage_EmptyCDPipeName(t *testing.T) {
	ctx := context.Background()
	stageCR := createStageCR(zeroOrder, "")
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{})

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&stageCR).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)

	stage, err := GetInputISForStage(ctx, k8sClient, name, cdPipelineName)
	assert.Error(t, err)
	assert.Nil(t, stage)
	assert.True(t, strings.Contains(err.Error(), "Spec.CdPipeline is empty in Stage CR named"))
}

func TestGetImageStreamFromStage_GetCdPipeErr(t *testing.T) {
	ctx := context.Background()
	stageCR := createStageCR(zeroOrder, name)
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{})

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&stageCR).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)

	stage, err := GetInputISForStage(ctx, k8sClient, name, cdPipelineName)
	assert.Error(t, err)
	assert.Nil(t, stage)
	assert.True(t, runtime.IsNotRegisteredError(err))
}

func TestGetImageStreamFromStage_EmptyImageStreamErr(t *testing.T) {
	ctx := context.Background()
	stageCR := createStageCR(zeroOrder, name)
	cdPipelineCR := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{}, &cdPipeApi.CDPipeline{})

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&stageCR, &cdPipelineCR).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)

	stage, err := GetInputISForStage(ctx, k8sClient, name, cdPipelineName)
	var emptyImageStreamEr *EmptyImageStreamErr
	assert.ErrorAs(t, err, &emptyImageStreamEr)
	assert.Nil(t, stage)
}

func TestGetImageStreamFromStage(t *testing.T) {
	ctx := context.Background()
	stageCR := createStageCR(zeroOrder, name)
	inputStreams := []string{"is1", "is2"}

	cdPipelineCR := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			InputDockerStreams: inputStreams,
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{}, &cdPipeApi.CDPipeline{})

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&stageCR, &cdPipelineCR).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)

	stage, err := GetInputISForStage(ctx, k8sClient, name, cdPipelineName)
	assert.NoError(t, err)
	assert.Equal(t, inputStreams, stage)
}

func TestGetImageStreamFromStage_NonZeroStageOrder(t *testing.T) {
	ctx := context.Background()
	stageCR := createStageCR(nonZeroOrder, cdPipelineName)
	stageCR.Annotations = make(map[string]string)
	stageCR.Annotations[PreviousStageNameAnnotationKey] = previousStageName

	cdPipelineCR := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cdPipelineName,
			Namespace: ns,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			ApplicationsToPromote: []string{firstImage, secondImage, "nonExistingImage"},
		},
	}

	firstImageStream := codeBaseApi.CodebaseImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createCISName(cdPipelineCR.Name, previousStageName, firstImage),
			Namespace: ns,
		},
	}

	secondImageStream := codeBaseApi.CodebaseImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createCISName(cdPipelineCR.Name, previousStageName, secondImage),
			Namespace: ns,
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{}, &cdPipeApi.CDPipeline{}, &codeBaseApi.CodebaseImageStream{})

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&stageCR, &cdPipelineCR, &firstImageStream, &secondImageStream).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)

	expectedResult := []string{createCISName(cdPipelineCR.Name, previousStageName, firstImage), createCISName(cdPipelineCR.Name, previousStageName, secondImage)}

	applicationsToPromote, err := GetInputISForStage(ctx, k8sClient, name, cdPipelineName)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, applicationsToPromote)
}

func TestGetImageStreamFromStage_BadApplicationsToPromoteValues(t *testing.T) {
	ctx := context.Background()
	stageCR := createStageCR(nonZeroOrder, cdPipelineName)
	stageCR.Annotations = make(map[string]string)
	stageCR.Annotations[PreviousStageNameAnnotationKey] = previousStageName

	cdPipelineCR := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cdPipelineName,
			Namespace: ns,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			ApplicationsToPromote: []string{"nonExistingImage"},
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{}, &cdPipeApi.CDPipeline{}, &codeBaseApi.CodebaseImageStream{})

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&stageCR, &cdPipelineCR).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)

	applicationsToPromote, err := GetInputISForStage(ctx, k8sClient, name, cdPipelineName)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "InputIS verification failed"))
	assert.Nil(t, applicationsToPromote)
}

func TestGetImageStreamFromStage_EmptyAnnotationToPromote(t *testing.T) {
	ctx := context.Background()
	stageCR := createStageCR(nonZeroOrder, cdPipelineName)
	stageCR.Annotations = make(map[string]string)
	stageCR.Annotations[PreviousStageNameAnnotationKey] = previousStageName

	cdPipelineCR := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cdPipelineName,
			Namespace: ns,
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{}, &cdPipeApi.CDPipeline{})

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&stageCR, &cdPipelineCR).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)

	applicationsToPromote, err := GetInputISForStage(ctx, k8sClient, name, cdPipelineName)
	var emptyImageStreamEr *EmptyImageStreamErr
	assert.ErrorAs(t, err, &emptyImageStreamEr)
	assert.Nil(t, applicationsToPromote)
}

func TestGetImageStreamFromStage_NoPreviousStageName(t *testing.T) {
	ctx := context.Background()
	stageCR := createStageCR(nonZeroOrder, cdPipelineName)

	cdPipelineCR := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cdPipelineName,
			Namespace: ns,
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{}, &cdPipeApi.CDPipeline{}, &codeBaseApi.CodebaseImageStream{})

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&stageCR, &cdPipelineCR).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)

	applicationsToPromote, err := GetInputISForStage(ctx, k8sClient, name, cdPipelineName)
	assert.Nil(t, applicationsToPromote)
	assert.True(t, strings.Contains(err.Error(), "there is no annotation"))
}

func TestGetOutputISForStage_GetErr(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{}, &cdPipeApi.CDPipeline{})
	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)
	imageStreams, err := GetOutputISForStage(ctx, k8sClient, cdPipelineName, name)
	assert.Error(t, err)
	assert.Nil(t, imageStreams)
}
func TestGetOutputISForStage_EmptyIsErr(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{}, &cdPipeApi.CDPipeline{})
	cdPipelineCR := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cdPipelineName,
			Namespace: ns,
		},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&cdPipelineCR).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)
	imageStreams, err := GetOutputISForStage(ctx, k8sClient, cdPipelineName, name)
	var emptyImageStreamEr *EmptyImageStreamErr
	assert.ErrorAs(t, err, &emptyImageStreamEr)
	assert.Nil(t, imageStreams)
}
func TestGetOutputISForStage(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	expectedIS := []string{createCISName(cdPipelineName, name, appName)}
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{}, &cdPipeApi.CDPipeline{})
	cdPipelineCR := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cdPipelineName,
			Namespace: ns,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			ApplicationsToPromote: []string{appName},
		},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&cdPipelineCR).Build()
	k8sClient := k8s.NewRuntimeNamespacedClient(client, ns)
	imageStreams, err := GetOutputISForStage(ctx, k8sClient, cdPipelineName, name)
	assert.NoError(t, err)
	assert.Equal(t, expectedIS, imageStreams)
}

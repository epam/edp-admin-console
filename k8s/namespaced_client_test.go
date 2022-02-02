package k8s

import (
	"context"
	"testing"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	name = "name"
	ns   = "ns"
	ns2  = "ns2"
)

func createcbBranchCRByNamespace(ns string) *codeBaseApi.CodebaseBranch {
	return &codeBaseApi.CodebaseBranch{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CodebaseBranch",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: codeBaseApi.CodebaseBranchSpec{
			BranchName: "master",
			Release:    true,
		},
	}
}

func createCDStageCR(ns string) *cdPipeApi.Stage {
	return &cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}

func createCDPipelineCR(ns string) *cdPipeApi.CDPipeline {
	return &cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}

func TestK8SClient_GetCBBranch(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	cbBranchCR := createcbBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(cbBranchCR).Build()
	k8sClient := NewRuntimeNamespacedClient(client, "ns")
	branch, err := k8sClient.GetCBBranch(ctx, "name")
	assert.NoError(t, err)
	assert.Equal(t, name, branch.Name)
	assert.Equal(t, ns, branch.Namespace)
}

func TestK8SClient_GetCBBranch_NotExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	cbBranchCR := createcbBranchCRByNamespace(ns2)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cbBranchCR).Build()
	k8sClient := NewRuntimeNamespacedClient(client, ns)
	branch, err := k8sClient.GetCBBranch(ctx, name)
	assert.True(t, k8serrors.IsNotFound(err))
	assert.Nil(t, branch)
}

func TestNamespacedClient_DeleteCBBranch_NotFound(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	cbBranchCR := createcbBranchCRByNamespace(ns2)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cbBranchCR).Build()
	k8sClient := NewRuntimeNamespacedClient(client, ns)
	err := k8sClient.DeleteCBBranch(ctx, name)
	assert.True(t, k8serrors.IsNotFound(err))
}

func TestNamespacedClient_DeleteCBBranch(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	cbBranchCR := createcbBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cbBranchCR).Build()
	k8sClient := NewRuntimeNamespacedClient(client, ns)
	err := k8sClient.DeleteCBBranch(ctx, name)
	assert.NoError(t, err)

	branch, err := k8sClient.GetCBBranch(ctx, name)
	assert.True(t, k8serrors.IsNotFound(err))
	assert.Nil(t, branch)
}

func TestNamespacedClient_UpdateCBBranchByCustomFields(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	initCR := createcbBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initCR).Build()
	spec := codeBaseApi.CodebaseBranchSpec{
		BranchName: "master",
		Release:    false,
	}

	k8sClient := NewRuntimeNamespacedClient(client, ns)
	err := k8sClient.UpdateCBBranchByCustomFields(ctx, name, spec, codeBaseApi.CodebaseBranchStatus{})
	assert.NoError(t, err)
	branch, err := k8sClient.GetCBBranch(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, spec.Release, branch.Spec.Release)
}

func TestNamespacedClient_UpdateCBBranchByCustomFields_NotExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects().Build()
	spec := codeBaseApi.CodebaseBranchSpec{
		BranchName: "master",
		Release:    false,
	}

	k8sClient := NewRuntimeNamespacedClient(client, ns)
	err := k8sClient.UpdateCBBranchByCustomFields(ctx, name, spec, codeBaseApi.CodebaseBranchStatus{})
	assert.Error(t, err)
	assert.True(t, k8serrors.IsNotFound(err))
}

func TestNamespacedClient_CreateCBBranchByCustomFields(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	initCR := createcbBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initCR).Build()
	spec := codeBaseApi.CodebaseBranchSpec{BranchName: "test"}

	k8sClient := NewRuntimeNamespacedClient(client, ns)
	err := k8sClient.CreateCBBranchByCustomFields(context.TODO(), name, spec, codeBaseApi.CodebaseBranchStatus{})
	assert.Error(t, err)
	assert.True(t, k8serrors.IsAlreadyExists(err))

	branch, err := k8sClient.GetCBBranch(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, "master", branch.Spec.BranchName)
}

func TestNamespacedClient_GetCBBranch_BadClient(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	CR := createcbBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(CR).Build()
	k8sClient := RuntimeNamespacedClient{
		Client: client,
	}
	branch, err := k8sClient.GetCBBranch(ctx, "name")
	assert.Nil(t, branch)
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
}

func TestNamespacedClient_UpdateCBBranchByCustomFields_BadClient(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	initCR := createcbBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initCR).Build()
	spec := codeBaseApi.CodebaseBranchSpec{
		BranchName: "master",
		Release:    false,
	}
	k8sClient := RuntimeNamespacedClient{
		Client: client,
	}

	err := k8sClient.UpdateCBBranchByCustomFields(ctx, name, spec, codeBaseApi.CodebaseBranchStatus{})
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
}

func TestNamespacedClient_DeleteCBBranch_BadClient(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	cbBranchCR := createcbBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(cbBranchCR).Build()
	k8sClient := RuntimeNamespacedClient{
		Client: client,
	}
	err := k8sClient.DeleteCBBranch(ctx, "name")
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
}

func TestNamespacedClient_CreateCBBranchByCustomFields_BadClient(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	initCR := createcbBranchCRByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initCR).Build()
	spec := codeBaseApi.CodebaseBranchSpec{BranchName: "test"}

	k8sClient := RuntimeNamespacedClient{
		Client: client,
	}
	err := k8sClient.CreateCBBranchByCustomFields(context.TODO(), name, spec, codeBaseApi.CodebaseBranchStatus{})
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
}

func TestK8SClient_GetCDStage_BadClient(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects().Build()
	k8sClient := RuntimeNamespacedClient{Client: client}
	stage, err := k8sClient.GetCDStage(ctx, name)
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
	assert.Nil(t, stage)
}

func TestK8SClient_GetCDStage_NotExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{})
	cdStageCR := createCDStageCR(ns2)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cdStageCR).Build()
	k8sClient := NewRuntimeNamespacedClient(client, ns)
	stage, err := k8sClient.GetCDStage(ctx, name)
	assert.True(t, k8serrors.IsNotFound(err))
	assert.Nil(t, stage)
}

func TestK8SClient_GetCDStage(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.Stage{})
	cdStageCR := createCDStageCR(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cdStageCR).Build()
	k8sClient := NewRuntimeNamespacedClient(client, ns)
	stage, err := k8sClient.GetCDStage(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, name, stage.Name)
	assert.Equal(t, ns, stage.Namespace)
}

func TestK8SClient_GetCDPipeline_BadClient(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects().Build()
	k8sClient := RuntimeNamespacedClient{Client: client}
	cdPipeline, err := k8sClient.GetCDPipeline(ctx, name)
	var emptyNamespaceErr *EmptyNamespaceErr
	assert.ErrorAs(t, err, &emptyNamespaceErr)
	assert.Nil(t, cdPipeline)
}

func TestK8SClient_GetCDPipeline_NotExist(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.CDPipeline{})
	cdPipelineCR := createCDPipelineCR(ns2)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cdPipelineCR).Build()
	k8sClient := NewRuntimeNamespacedClient(client, ns)
	cdPipeline, err := k8sClient.GetCDPipeline(ctx, name)
	assert.True(t, k8serrors.IsNotFound(err))
	assert.Nil(t, cdPipeline)
}

func TestK8SClient_GetCDPipeline(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &cdPipeApi.CDPipeline{})
	cdPipelineCR := createCDPipelineCR(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cdPipelineCR).Build()
	k8sClient := NewRuntimeNamespacedClient(client, ns)
	cdPipeline, err := k8sClient.GetCDPipeline(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, name, cdPipeline.Name)
	assert.Equal(t, ns, cdPipeline.Namespace)
}

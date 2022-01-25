package k8s

import (
	ctx "context"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

const (
	name = "name"
	ns   = "ns"
	ns2  = "ns2"
)

func createCBBranchInstanceByNamespace(ns string) *codeBaseApi.CodebaseBranch {
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

func TestK8SClient_GetCBBranch(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	instance := createCBBranchInstanceByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(instance).Build()
	k8sClient := NewNamespacedClient(client, "ns")
	branch, err := k8sClient.GetCBBranch(ctx.TODO(), "name")
	assert.NoError(t, err)
	assert.Equal(t, name, branch.Name)
	assert.Equal(t, ns, branch.Namespace)
}

func TestK8SClient_GetCBBranch_NotExist(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	instance := createCBBranchInstanceByNamespace(ns2)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(instance).Build()
	k8sClient := NewNamespacedClient(client, ns)
	branch, err := k8sClient.GetCBBranch(ctx.TODO(), name)
	assert.True(t, k8serrors.IsNotFound(err))
	assert.Empty(t, branch)
}

func TestNamespacedClient_DeleteCBBranch_NotFound(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	instance := createCBBranchInstanceByNamespace(ns2)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(instance).Build()
	k8sClient := NewNamespacedClient(client, ns)
	err := k8sClient.DeleteCBBranch(ctx.TODO(), name)
	assert.True(t, k8serrors.IsNotFound(err))
}

func TestNamespacedClient_DeleteCBBranch(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	instance := createCBBranchInstanceByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(instance).Build()
	k8sClient := NewNamespacedClient(client, ns)
	err := k8sClient.DeleteCBBranch(ctx.TODO(), name)
	assert.NoError(t, err)
	branch, err := k8sClient.GetCBBranch(ctx.TODO(), name)
	assert.True(t, k8serrors.IsNotFound(err))
	assert.Empty(t, branch)
}

func TestNamespacedClient_UpdateCBBranchByCustomFields(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	initInstance := createCBBranchInstanceByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initInstance).Build()
	spec := codeBaseApi.CodebaseBranchSpec{
		BranchName: "master",
		Release:    false,
	}

	k8sClient := NewNamespacedClient(client, ns)
	err := k8sClient.UpdateCBBranchByCustomFields(ctx.TODO(), name, spec, codeBaseApi.CodebaseBranchStatus{})
	assert.NoError(t, err)
	branch, err := k8sClient.GetCBBranch(ctx.TODO(), name)
	assert.NoError(t, err)
	assert.Equal(t, spec.Release, branch.Spec.Release)
}

func TestNamespacedClient_UpdateCBBranchByCustomFields_NotExist(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects().Build()
	spec := codeBaseApi.CodebaseBranchSpec{
		BranchName: "master",
		Release:    false,
	}

	k8sClient := NewNamespacedClient(client, ns)
	err := k8sClient.UpdateCBBranchByCustomFields(ctx.TODO(), name, spec, codeBaseApi.CodebaseBranchStatus{})
	assert.Error(t, err)
	assert.True(t, k8serrors.IsNotFound(err))
}

func TestNamespacedClient_CreateCBBranchByCustomFields(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	initInstance := createCBBranchInstanceByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initInstance).Build()
	spec := codeBaseApi.CodebaseBranchSpec{BranchName: "test"}

	k8sClient := NewNamespacedClient(client, ns)
	err := k8sClient.CreateCBBranchByCustomFields(context.TODO(), name, spec, codeBaseApi.CodebaseBranchStatus{})
	assert.Error(t, err)
	assert.True(t, k8serrors.IsAlreadyExists(err))

	branch, err := k8sClient.GetCBBranch(ctx.TODO(), name)
	assert.NoError(t, err)
	assert.Equal(t, "master", branch.Spec.BranchName)
}

func TestNamespacedClient_GetCBBranch_BadClient(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	instance := createCBBranchInstanceByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(instance).Build()
	k8sClient := NamespacedClient{
		Client: client,
	}
	branch, err := k8sClient.GetCBBranch(ctx.TODO(), "name")
	assert.Error(t, err)
	assert.Empty(t, branch)
	assert.True(t, AsEmptyNamespaceErr(err))
}

func TestNamespacedClient_UpdateCBBranchByCustomFields_BadClient(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	initInstance := createCBBranchInstanceByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initInstance).Build()
	spec := codeBaseApi.CodebaseBranchSpec{
		BranchName: "master",
		Release:    false,
	}
	k8sClient := NamespacedClient{
		Client: client,
	}

	err := k8sClient.UpdateCBBranchByCustomFields(ctx.TODO(), name, spec, codeBaseApi.CodebaseBranchStatus{})
	assert.Error(t, err)
	assert.True(t, AsEmptyNamespaceErr(err))
}

func TestNamespacedClient_DeleteCBBranch_BadClient(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})
	instance := createCBBranchInstanceByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(instance).Build()
	k8sClient := NamespacedClient{
		Client: client,
	}
	err := k8sClient.DeleteCBBranch(ctx.TODO(), "name")
	assert.Error(t, err)
	assert.True(t, AsEmptyNamespaceErr(err))
}

func TestNamespacedClient_CreateCBBranchByCustomFields_BadClient(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &codeBaseApi.CodebaseBranch{})

	initInstance := createCBBranchInstanceByNamespace(ns)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initInstance).Build()
	spec := codeBaseApi.CodebaseBranchSpec{BranchName: "test"}

	k8sClient := NamespacedClient{
		Client: client,
	}
	err := k8sClient.CreateCBBranchByCustomFields(context.TODO(), name, spec, codeBaseApi.CodebaseBranchStatus{})
	assert.Error(t, err)
	assert.True(t, AsEmptyNamespaceErr(err))
}

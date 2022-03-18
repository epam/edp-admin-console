package k8s

import (
	"context"
	"testing"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func createCodebaseCRWithOptions(opts ...CodebaseCROption) *codeBaseApi.Codebase {
	codebaseCR := new(codeBaseApi.Codebase)
	for i := range opts {
		opts[i](codebaseCR)
	}
	return codebaseCR
}

type CodebaseCROption func(codebase *codeBaseApi.Codebase)

func WithCBCrName(crName string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.ObjectMeta.Name = crName
	}
}

func WithCBCrNamespace(crNamespace string) CodebaseCROption {
	return func(codebase *codeBaseApi.Codebase) {
		codebase.ObjectMeta.Namespace = crNamespace
	}
}

func TestRuntimeNamespacedClient_DeleteCodebase(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion,
		&codeBaseApi.Codebase{},
	)
	testNamespace_1 := "namespace_1"
	testCodebaseName_1 := "codebase_1"
	codebaseCR_1 := createCodebaseCRWithOptions(
		WithCBCrName(testCodebaseName_1),
		WithCBCrNamespace(testNamespace_1),
	)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(codebaseCR_1).Build()
	k8sClient, err := NewRuntimeNamespacedClient(fakeClient, testNamespace_1)
	if err != nil {
		t.Fatal(err)
	}

	err = k8sClient.DeleteCodebase(ctx, codebaseCR_1)
	assert.NoError(t, err)
	k8sCodebase, err := k8sClient.GetCodebase(ctx, testCodebaseName_1)
	assert.True(t, k8sErrors.IsNotFound(err))

	var expectedK8SCodebase *codeBaseApi.Codebase
	assert.Equal(t, expectedK8SCodebase, k8sCodebase)
}

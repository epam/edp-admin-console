package k8s

import (
	"context"
	"testing"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type CodebaseImageStreamCROption func(codebaseImageStreamCR *codeBaseApi.CodebaseImageStream)

func WithCISCRName(crName string) CodebaseImageStreamCROption {
	return func(codebaseImageStreamCR *codeBaseApi.CodebaseImageStream) {
		codebaseImageStreamCR.Name = crName
	}
}

func WithCISCRNamespace(crNamespace string) CodebaseImageStreamCROption {
	return func(codebaseImageStreamCR *codeBaseApi.CodebaseImageStream) {
		codebaseImageStreamCR.Namespace = crNamespace
	}
}

func createCISCRWithOptions(opts ...CodebaseImageStreamCROption) *codeBaseApi.CodebaseImageStream {
	cisCR := new(codeBaseApi.CodebaseImageStream)
	for _, opt := range opts {
		opt(cisCR)
	}
	return cisCR
}

func TestRuntimeNamespacedClient_CodebaseImageStreamList_OK(t *testing.T) {
	ctx := context.Background()
	namespace := "test_ns_1"
	cisCRName_1 := "codebaseImageStream_1"

	cisCR_1 := createCISCRWithOptions(
		WithCISCRName(cisCRName_1),
		WithCISCRNamespace(namespace),
	)

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion,
		&codeBaseApi.CodebaseImageStream{}, &codeBaseApi.CodebaseImageStreamList{},
	)
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cisCR_1).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	expectedList := []codeBaseApi.CodebaseImageStream{*cisCR_1}
	gotList, err := k8sClient.CodebaseImageStreamList(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedList, gotList)
}

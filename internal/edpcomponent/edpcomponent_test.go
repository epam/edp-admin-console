package edpcomponent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	edpComponentAPI "github.com/epam/edp-component-operator/pkg/apis/v1/v1"

	"edp-admin-console/k8s"
)

type EDPComponentCROption func(codebase *edpComponentAPI.EDPComponent)

func createEDPComponentCRWithOptions(opts ...EDPComponentCROption) *edpComponentAPI.EDPComponent {
	jiraServerCR := new(edpComponentAPI.EDPComponent)
	for i := range opts {
		opts[i](jiraServerCR)
	}
	return jiraServerCR
}

func WithEDPComponentCrName(crName string) EDPComponentCROption {
	return func(codebase *edpComponentAPI.EDPComponent) {
		codebase.ObjectMeta.Name = crName
	}
}

func WithEDPComponentCrNamespace(crNamespace string) EDPComponentCROption {
	return func(codebase *edpComponentAPI.EDPComponent) {
		codebase.ObjectMeta.Namespace = crNamespace
	}
}

func TestByNameIFExists_Exists(t *testing.T) {
	ctx := context.Background()
	namespace := "test_ns_1"
	crName := "gerrit_1"
	edpComponentCR := createEDPComponentCRWithOptions(
		WithEDPComponentCrName(crName),
		WithEDPComponentCrNamespace(namespace),
	)

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(edpComponentAPI.SchemeGroupVersion, &edpComponentAPI.EDPComponent{})
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(edpComponentCR).Build()
	k8sClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	expectedEDPComponentCR := &edpComponentAPI.EDPComponent{
		ObjectMeta: metaV1.ObjectMeta{
			Name:            crName,
			Namespace:       namespace,
			ResourceVersion: "999",
		},
		TypeMeta: metaV1.TypeMeta{
			Kind:       "EDPComponent",
			APIVersion: "v1.edp.epam.com/v1",
		},
	}
	gotEDPComponent, err := ByNameIFExists(ctx, k8sClient, crName)

	assert.NoError(t, err)
	assert.Equal(t, expectedEDPComponentCR, gotEDPComponent)
}

func TestByNameIFExists_NotFound(t *testing.T) {
	ctx := context.Background()
	namespace := "test_ns_1"
	crName := "gerrit_1"

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(edpComponentAPI.SchemeGroupVersion, &edpComponentAPI.EDPComponent{})
	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	k8sClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	gotEDPComponent, err := ByNameIFExists(ctx, k8sClient, crName)

	assert.NoError(t, err)
	var expectedEDPComponent *edpComponentAPI.EDPComponent
	assert.Equal(t, expectedEDPComponent, gotEDPComponent)
}

func TestByNameIFExists_UnknownType(t *testing.T) {
	ctx := context.Background()
	namespace := "test_ns_1"
	crName := "gerrit_1"

	scheme := runtime.NewScheme()
	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	k8sClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	gotEDPComponent, err := ByNameIFExists(ctx, k8sClient, crName)

	assert.Error(t, err)
	var expectedEDPComponent *edpComponentAPI.EDPComponent
	assert.Equal(t, expectedEDPComponent, gotEDPComponent)
}

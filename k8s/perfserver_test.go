package k8s

import (
	"context"
	"testing"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	perfApi "github.com/epam/edp-perf-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestRuntimeNamespacedClient_GetPerfServerList(t *testing.T) {
	ctx := context.Background()

	perfServer := perfApi.PerfServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion, &perfApi.PerfServerList{}, &perfApi.PerfServer{})
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&perfServer).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}

	expectedPerfServer := perfApi.PerfServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       ns,
			ResourceVersion: "999",
		},
	}
	expectedList := []perfApi.PerfServer{expectedPerfServer}

	gotPerfServer, err := k8sClient.GetPerfServerList(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedList, gotPerfServer.Items)
}

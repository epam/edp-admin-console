package k8s

import (
	"context"
	"testing"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestRuntimeNamespacedClient_GetGitServerList(t *testing.T) {
	ctx := context.Background()

	gitServer := codeBaseApi.GitServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion, &codeBaseApi.GitServerList{}, &codeBaseApi.GitServer{})
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&gitServer).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}

	expectedGitServer := codeBaseApi.GitServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       ns,
			ResourceVersion: "999",
		},
	}
	expectedList := []codeBaseApi.GitServer{expectedGitServer}

	gotGitServer, err := k8sClient.GetGitServerList(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedList, gotGitServer.Items)
}

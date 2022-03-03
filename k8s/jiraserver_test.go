package k8s

import (
	"context"
	"testing"

	"github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type JiraServerCROption func(codebase *v1alpha1.JiraServer)

func createJiraServerCRWithOptions(opts ...JiraServerCROption) *v1alpha1.JiraServer {
	jiraServerCR := new(v1alpha1.JiraServer)
	for i := range opts {
		opts[i](jiraServerCR)
	}
	return jiraServerCR
}

func WithCrName(crName string) JiraServerCROption {
	return func(codebase *v1alpha1.JiraServer) {
		codebase.ObjectMeta.Name = crName
	}
}

func WithCrNamespace(crNamespace string) JiraServerCROption {
	return func(codebase *v1alpha1.JiraServer) {
		codebase.ObjectMeta.Namespace = crNamespace
	}
}

func TestJiraServersList_Success(t *testing.T) {
	ctx := context.Background()
	namespace := "test_ns_1"
	crName := "jira_server_1"
	jiraServerCR := createJiraServerCRWithOptions(
		WithCrName(crName),
		WithCrNamespace(namespace),
	)

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(v1alpha1.SchemeGroupVersion, &v1alpha1.JiraServer{}, &v1alpha1.JiraServerList{})
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(jiraServerCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	expectedJiraServerCR := v1alpha1.JiraServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:            crName,
			Namespace:       namespace,
			ResourceVersion: "999",
		},
	}
	expectedList := []v1alpha1.JiraServer{expectedJiraServerCR}
	gotJiraServers, err := k8sClient.JiraServersList(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedList, gotJiraServers)
}

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

type JiraServerCROption func(codebase *codeBaseApi.JiraServer)

func createJiraServerCRWithOptions(opts ...JiraServerCROption) *codeBaseApi.JiraServer {
	jiraServerCR := new(codeBaseApi.JiraServer)
	for i := range opts {
		opts[i](jiraServerCR)
	}
	return jiraServerCR
}

func WithCrName(crName string) JiraServerCROption {
	return func(codebase *codeBaseApi.JiraServer) {
		codebase.ObjectMeta.Name = crName
	}
}

func WithCrNamespace(crNamespace string) JiraServerCROption {
	return func(codebase *codeBaseApi.JiraServer) {
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
	scheme.AddKnownTypes(codeBaseApi.SchemeGroupVersion, &codeBaseApi.JiraServer{}, &codeBaseApi.JiraServerList{})
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(jiraServerCR).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	expectedJiraServerCR := codeBaseApi.JiraServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:            crName,
			Namespace:       namespace,
			ResourceVersion: "999",
		},
	}
	expectedList := []codeBaseApi.JiraServer{expectedJiraServerCR}
	gotJiraServers, err := k8sClient.JiraServersList(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedList, gotJiraServers)
}

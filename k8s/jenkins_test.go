package k8s

import (
	"context"
	"testing"

	jenkinsApi "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1alpha1"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestRuntimeNamespacedClient_GetJenkinsList(t *testing.T) {
	ctx := context.Background()

	jenkins := jenkinsApi.Jenkins{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(jenkinsApi.SchemeGroupVersion, &jenkinsApi.JenkinsList{}, &jenkinsApi.Jenkins{})
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&jenkins).Build()
	k8sClient, err := NewRuntimeNamespacedClient(client, ns)
	if err != nil {
		t.Fatal(err)
	}

	expectedJenkins := jenkinsApi.Jenkins{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       ns,
			ResourceVersion: "999",
		},
	}
	expectedList := []jenkinsApi.Jenkins{expectedJenkins}

	gotJenkins, err := k8sClient.GetJenkinsList(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedList, gotJenkins.Items)
}

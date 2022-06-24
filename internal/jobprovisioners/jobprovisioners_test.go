package jobprovisioners

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	jenkinsApi "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1"

	"edp-admin-console/k8s"
)

const testNamespace = "test_namespace_1"

func TestList_OK(t *testing.T) {
	namespace := testNamespace
	crName_1 := "jenkins_cr_1"
	ctx := context.Background()

	jobProvisioner_1_Name := "jenkins_job_name_1"
	jobProvisioner_1_Scope := "jenkins_job_scope_1"

	jobProvisioners := []jenkinsApi.JobProvision{
		{
			Name:  jobProvisioner_1_Name,
			Scope: jobProvisioner_1_Scope,
		},
	}
	jenkins := jenkinsApi.Jenkins{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      crName_1,
			Namespace: namespace,
		},
		Status: jenkinsApi.JenkinsStatus{
			JobProvisions: jobProvisioners,
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(jenkinsApi.SchemeGroupVersion, &jenkinsApi.JenkinsList{}, &jenkinsApi.Jenkins{})
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&jenkins).Build()
	k8sClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	expectedList := jobProvisioners

	gotJobProvisioners, err := List(ctx, k8sClient)
	assert.NoError(t, err)
	assert.Equal(t, expectedList, gotJobProvisioners)
}

func TestListNames_OK(t *testing.T) {
	namespace := testNamespace
	crName_1 := "jenkins_cr_1"
	ctx := context.Background()

	jobProvisioner_1_Name := "jenkins_job_name_1"
	jobProvisioner_1_Scope := "jenkins_job_scope_1"
	jobProvisioner_2_Name := jobProvisioner_1_Name
	jobProvisioner_2_Scope := "jenkins_job_scope_2"

	jobProvisioners := []jenkinsApi.JobProvision{
		{
			Name:  jobProvisioner_1_Name,
			Scope: jobProvisioner_1_Scope,
		},
		{
			Name:  jobProvisioner_2_Name,
			Scope: jobProvisioner_2_Scope,
		},
	}
	jenkins := jenkinsApi.Jenkins{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      crName_1,
			Namespace: namespace,
		},
		Status: jenkinsApi.JenkinsStatus{
			JobProvisions: jobProvisioners,
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(jenkinsApi.SchemeGroupVersion, &jenkinsApi.JenkinsList{}, &jenkinsApi.Jenkins{})
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(&jenkins).Build()
	k8sClient, err := k8s.NewRuntimeNamespacedClient(client, namespace)
	if err != nil {
		t.Fatal(err)
	}

	expectedList := []string{jobProvisioner_1_Name}

	gotJobProvisioners, err := ListNames(ctx, k8sClient)
	assert.NoError(t, err)
	assert.Equal(t, expectedList, gotJobProvisioners)
}

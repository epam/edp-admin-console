package k8s

import (
	"context"

	jenkinsAPI "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1alpha1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// GetJenkinsList retrieves all Jenkins structure ptr for the given custom resource name from the Kubernetes Cluster CR.
func (c *RuntimeNamespacedClient) GetJenkinsList(ctx context.Context) (*jenkinsAPI.JenkinsList, error) {
	if c.Namespace == "" {
		return nil, NemEmptyNamespaceErr("client namespace is not set")
	}

	jenkinsList := &jenkinsAPI.JenkinsList{}
	err := c.List(ctx, jenkinsList, &runtimeClient.ListOptions{
		Namespace: c.Namespace,
	})
	if err != nil {
		return nil, err
	}
	return jenkinsList, err
}

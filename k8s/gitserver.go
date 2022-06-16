package k8s

import (
	"context"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// GetGitServerList retrieves all GitServer structure ptr for the given custom resource name from the Kubernetes Cluster CR.
func (c *RuntimeNamespacedClient) GetGitServerList(ctx context.Context) (*codeBaseApi.GitServerList, error) {
	if c.Namespace == "" {
		return nil, NemEmptyNamespaceErr("client namespace is not set")
	}

	gitServerList := &codeBaseApi.GitServerList{}
	err := c.List(ctx, gitServerList, &runtimeClient.ListOptions{
		Namespace: c.Namespace,
	})
	if err != nil {
		return nil, err
	}
	return gitServerList, err
}

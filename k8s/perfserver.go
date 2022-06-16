package k8s

import (
	"context"

	perfApi "github.com/epam/edp-perf-operator/v2/pkg/apis/edp/v1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// GetPerfServerList retrieves all PerfServer structure ptr for the given custom resource name from the Kubernetes Cluster CR.
func (c *RuntimeNamespacedClient) GetPerfServerList(ctx context.Context) (*perfApi.PerfServerList, error) {
	if c.Namespace == "" {
		return nil, NemEmptyNamespaceErr("client namespace is not set")
	}

	perfServerList := &perfApi.PerfServerList{}
	err := c.List(ctx, perfServerList, &runtimeClient.ListOptions{
		Namespace: c.Namespace,
	})
	if err != nil {
		return nil, err
	}
	return perfServerList, err
}

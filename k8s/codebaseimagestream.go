package k8s

import (
	"context"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// CodebaseImageStreamList returns full list of codebase image streams
func (c *RuntimeNamespacedClient) CodebaseImageStreamList(ctx context.Context) ([]codeBaseApi.CodebaseImageStream, error) {
	if c.Namespace == "" {
		return nil, NemEmptyNamespaceErr("client namespace is not set")
	}
	codebaseImageStream := new(codeBaseApi.CodebaseImageStreamList)
	err := c.List(ctx, codebaseImageStream, &runtimeClient.ListOptions{
		Namespace: c.Namespace,
	})
	if err != nil {
		return nil, err
	}
	return codebaseImageStream.Items, err
}

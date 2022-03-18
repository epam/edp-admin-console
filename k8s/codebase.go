package k8s

import (
	"context"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
)

func (c *RuntimeNamespacedClient) DeleteCodebase(ctx context.Context, codebaseCR *codeBaseApi.Codebase) error {
	return c.Delete(ctx, codebaseCR)
}

package k8s

import (
	"context"

	cdPipelineAPI "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
)

func (c *RuntimeNamespacedClient) DeleteCDPipeline(ctx context.Context, cdPipelineCR *cdPipelineAPI.CDPipeline) error {
	return c.Delete(ctx, cdPipelineCR)
}

package k8s

import (
	"context"

	cdPipelineAPI "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func (c *RuntimeNamespacedClient) StageList(ctx context.Context) ([]cdPipelineAPI.Stage, error) {
	stageList := new(cdPipelineAPI.StageList)
	err := c.List(ctx, stageList, &runtimeClient.ListOptions{Namespace: c.Namespace})
	if err != nil {
		return nil, err
	}
	return stageList.Items, nil
}

func (c *RuntimeNamespacedClient) DeleteStage(ctx context.Context, stageCR *cdPipelineAPI.Stage) error {
	return c.Delete(ctx, stageCR)
}

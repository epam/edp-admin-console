package k8s

import (
	"context"

	cdPipelineAPI "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func (c *RuntimeNamespacedClient) DeleteCDPipeline(ctx context.Context, cdPipelineCR *cdPipelineAPI.CDPipeline) error {
	return c.Delete(ctx, cdPipelineCR)
}

// CreateCDPipelineBySpec creates CDPipeline CR by custom field and name
func (c *RuntimeNamespacedClient) CreateCDPipelineBySpec(ctx context.Context, crName string, spec cdPipelineAPI.CDPipelineSpec) error {
	if c.Namespace == "" {
		return NemEmptyNamespaceErr("client namespace is not set")
	}
	cdPipeline := &cdPipelineAPI.CDPipeline{
		ObjectMeta: v1.ObjectMeta{
			Name:      crName,
			Namespace: c.Namespace,
		},
		Spec: spec,
	}
	err := c.Create(ctx, cdPipeline)
	return err
}

// CreateCDStageBySpec creates CDPipeline CR by custom field and name
func (c *RuntimeNamespacedClient) CreateCDStageBySpec(ctx context.Context, crName string, spec cdPipelineAPI.StageSpec) error {
	if c.Namespace == "" {
		return NemEmptyNamespaceErr("client namespace is not set")
	}
	Stage := &cdPipelineAPI.Stage{
		ObjectMeta: v1.ObjectMeta{
			Name:      crName,
			Namespace: c.Namespace,
		},
		Spec: spec,
	}
	err := c.Create(ctx, Stage)
	return err
}

func (c *RuntimeNamespacedClient) CDPipelineList(ctx context.Context) ([]cdPipelineAPI.CDPipeline, error) {
	list := new(cdPipelineAPI.CDPipelineList)
	err := c.List(ctx, list, &runtimeClient.ListOptions{Namespace: c.Namespace})
	return list.Items, err
}

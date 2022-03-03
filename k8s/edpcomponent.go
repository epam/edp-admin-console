package k8s

import (
	"context"

	edpComponentAPI "github.com/epam/edp-component-operator/pkg/apis/v1/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

func (c *RuntimeNamespacedClient) EDPComponentByCRName(ctx context.Context, componentCRName string) (*edpComponentAPI.EDPComponent, error) {
	edpComponent := new(edpComponentAPI.EDPComponent)
	namespacedName := types.NamespacedName{
		Name:      componentCRName,
		Namespace: c.Namespace,
	}
	err := c.Get(ctx, namespacedName, edpComponent)
	if err != nil {
		return nil, err
	}
	return edpComponent, nil
}

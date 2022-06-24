package k8s

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"

	edpComponentAPI "github.com/epam/edp-component-operator/pkg/apis/v1/v1"
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

func (c *RuntimeNamespacedClient) EDPComponentList(ctx context.Context) ([]edpComponentAPI.EDPComponent, error) {
	edpComponentList := &edpComponentAPI.EDPComponentList{}
	err := c.List(ctx, edpComponentList, &runtimeClient.ListOptions{Namespace: c.Namespace})
	if err != nil {
		return nil, err
	}
	return edpComponentList.Items, nil
}

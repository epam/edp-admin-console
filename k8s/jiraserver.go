package k8s

import (
	"context"

	"github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func (c *RuntimeNamespacedClient) JiraServersList(ctx context.Context) ([]v1alpha1.JiraServer, error) {
	jiraServersList := new(v1alpha1.JiraServerList)
	err := c.List(ctx, jiraServersList, &runtimeClient.ListOptions{Namespace: c.Namespace})
	if err != nil {
		return nil, err
	}
	return jiraServersList.Items, nil
}

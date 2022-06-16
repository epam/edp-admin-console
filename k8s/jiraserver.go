package k8s

import (
	"context"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func (c *RuntimeNamespacedClient) JiraServersList(ctx context.Context) ([]codeBaseApi.JiraServer, error) {
	jiraServersList := new(codeBaseApi.JiraServerList)
	err := c.List(ctx, jiraServersList, &runtimeClient.ListOptions{Namespace: c.Namespace})
	if err != nil {
		return nil, err
	}
	return jiraServersList.Items, nil
}

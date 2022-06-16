package cdpipelines

import (
	"context"
	"errors"
	"net/http"

	"edp-admin-console/k8s"

	cdPipelineAPI "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
)

func ByNameIFExists(ctx context.Context, k8sClient *k8s.RuntimeNamespacedClient, cdName string) (*cdPipelineAPI.CDPipeline, error) {
	edpComponent, err := k8sClient.GetCDPipeline(ctx, cdName)
	if err != nil {
		var statusErr *k8sErrors.StatusError
		if errors.As(err, &statusErr) {
			if statusErr.ErrStatus.Code == http.StatusNotFound {
				return nil, nil // not found
			}
		}
		return nil, err
	}

	return edpComponent, nil
}

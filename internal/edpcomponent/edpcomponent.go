package edpcomponent

import (
	"context"
	"errors"
	"net/http"

	edpComponentAPI "github.com/epam/edp-component-operator/pkg/apis/v1/v1alpha1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"

	"edp-admin-console/k8s"
)

func ByNameIFExists(ctx context.Context, client *k8s.RuntimeNamespacedClient, componentName string) (*edpComponentAPI.EDPComponent, error) {
	edpComponent, err := client.EDPComponentByCRName(ctx, componentName)
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

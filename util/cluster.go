package util

import (
	"edp-admin-console/context"
	"edp-admin-console/util/consts"
	edpv1alpha1 "github.com/epmd-edp/codebase-operator/v2/pkg/apis/edp/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
)

func GetCodebaseCR(c *rest.RESTClient, name string) (*edpv1alpha1.Codebase, error) {
	r := &edpv1alpha1.Codebase{}
	err := c.Get().Namespace(context.Namespace).Resource(consts.CodebasePlural).Name(name).Do().Into(r)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r, nil
}

package util

import (
	ctx "context"
	"edp-admin-console/context"
	"edp-admin-console/util/consts"
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
)

func GetCodebaseCR(c *rest.RESTClient, name string) (*codeBaseApi.Codebase, error) {
	r := &codeBaseApi.Codebase{}
	err := c.Get().Namespace(context.Namespace).Resource(consts.CodebasePlural).Name(name).Do(ctx.TODO()).Into(r)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r, nil
}

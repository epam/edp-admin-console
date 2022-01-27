package filters

import (
	"fmt"

	bgCtx "github.com/astaxie/beego/context"
	"go.uber.org/zap"
)

type AccessControlEnv struct {
	Permissions map[string][]string
}

func (ac *AccessControlEnv) RoleAccessControlFilter(context *bgCtx.Context) {
	log.Debug("Start Role Access Control filter..")
	contextRoles := context.Input.Session("realm_roles").([]string)

	permissions := ac.Permissions
	isPageAvailable, err := IsPageAvailable(permissions, fmt.Sprintf("%s %s", context.Input.Method(), context.Input.URI()), contextRoles)
	if err != nil {
		log.Error("define is page available failed", zap.Error(err),
			zap.String("method", context.Input.Method()),
			zap.String("url", context.Input.URI()),
		)
		context.Abort(200, "403")
		return
	}

	if !isPageAvailable {
		log.Error("Access is denied", zap.String("url", context.Input.URI()))
		context.Abort(200, "403")
		return
	}
}

package filters

import (
	"fmt"
	"net/http"

	bgCtx "github.com/astaxie/beego/context"
	"go.uber.org/zap"
)

func (ac *AccessControlEnv) RoleAccessControlRestFilter(context *bgCtx.Context) {
	log.Debug("Start Role Access Control filter..")
	contextRoles, ok := context.Input.Session("realm_roles").([]string)
	if !ok {
		log.Error("invalid realm_roles, expected []string",
			zap.Any("realm_roles", context.Input.Session("realm_roles")),
			zap.Any("realm_roles_type", fmt.Sprintf("%T", context.Input.Session("realm_roles"))),
		)
		http.Error(context.ResponseWriter, "Forbidden.", http.StatusForbidden)
		return
	}

	permissions := ac.Permissions
	isPageAvailable, err := IsPageAvailable(permissions, fmt.Sprintf("%s %s", context.Input.Method(), context.Input.URI()), contextRoles)
	if err != nil {
		log.Error("define is page available failed", zap.Error(err),
			zap.String("method", context.Input.Method()),
			zap.String("url", context.Input.URI()),
		)
		http.Error(context.ResponseWriter, "Forbidden.", http.StatusForbidden)
		return
	}

	if !isPageAvailable {
		log.Error("Access is denied", zap.String("url", context.Input.URI()))
		http.Error(context.ResponseWriter, "Forbidden.", http.StatusForbidden)
		return
	}
}

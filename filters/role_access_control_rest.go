package filters

import (
	"fmt"
	bgCtx "github.com/astaxie/beego/context"
	"go.uber.org/zap"
	"net/http"
)

func RoleAccessControlRestFilter(context *bgCtx.Context) {
	log.Debug("Start Role Access Control filter..")
	contextRoles := context.Input.Session("realm_roles").([]string)

	isPageAvailable := IsPageAvailable(fmt.Sprintf("%s %s", context.Input.Method(), context.Input.URI()), contextRoles)

	if !isPageAvailable {
		log.Error("Access is denied", zap.String("url", context.Input.URI()))
		http.Error(context.ResponseWriter, "Forbidden.", http.StatusForbidden)
		return
	}
}

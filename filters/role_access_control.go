package filters

import (
	"fmt"
	bgCtx "github.com/astaxie/beego/context"
	"go.uber.org/zap"
)

func RoleAccessControlFilter(context *bgCtx.Context) {
	log.Debug("Start Role Access Control filter..")
	contextRoles := context.Input.Session("realm_roles").([]string)

	isPageAvailable := IsPageAvailable(fmt.Sprintf("%s %s", context.Input.Method(), context.Input.URI()), contextRoles)

	if !isPageAvailable {
		log.Error("Access is denied", zap.String("url", context.Input.URI()))
		context.Abort(200, "403")
		return
	}
}

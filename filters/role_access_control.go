package filters

import (
	"fmt"
	"github.com/astaxie/beego"
	bgCtx "github.com/astaxie/beego/context"
	"log"
)

func RoleAccessControlFilter(context *bgCtx.Context) {
	log.Println("Start Role Access Control filter..")
	resourceAccess := context.Input.Session("resource_access").(map[string][]string)
	edpTenantName := beego.AppConfig.String("cicdNamespace")
	contextRoles := resourceAccess[edpTenantName+"-edp"]

	if contextRoles == nil {
		contextRoles = context.Input.Session("realm_roles").([]string)
	}

	isPageAvailable := IsPageAvailable(fmt.Sprintf("%s %s", context.Input.Method(), context.Input.URI()), contextRoles)

	if !isPageAvailable {
		context.Abort(200, "403")
		return
	}
}

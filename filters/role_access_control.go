package filters

import (
	"fmt"
	bgCtx "github.com/astaxie/beego/context"
	"log"
	"net/http"
)

func RoleAccessControlFilter(context *bgCtx.Context) {
	log.Println("Start Role Access Control filter..")
	resourceAccess := context.Input.Session("resource_access").(map[string][]string)
	contextRoles := resourceAccess[context.Input.Param(":name")+"-edp"]

	if contextRoles == nil {
		nonValidTenantMsg := fmt.Sprintf("Couldn't find tenant by %s name.", context.Input.Param(":name"))
		http.Error(context.ResponseWriter, nonValidTenantMsg, http.StatusNotFound)
		return
	}

	isPageAvailable := IsPageAvailable(fmt.Sprintf("%s %s", context.Input.Method(), context.Input.URI()), contextRoles)

	if !isPageAvailable {
		context.Abort(200, "403")
		return
	}
}

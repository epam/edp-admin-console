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
	isPageAvailable := IsPageAvailable(fmt.Sprintf("%s %s", context.Input.Method(), context.Input.URI()), contextRoles)

	if !isPageAvailable {
		http.Error(context.ResponseWriter, "Forbidden.", http.StatusForbidden)
		return
	}
}

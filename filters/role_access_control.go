package filters

import (
	bgCtx "github.com/astaxie/beego/context"
	"log"
	"net/http"
	"strings"
)

func RoleAccessControlFilter(context *bgCtx.Context) {
	log.Println("Start Role Access Control filter..")
	resourceAccess := context.Input.Session("resource_access").(map[string][]string)
	edpName := context.Input.Param(":name")
	contextRoles := resourceAccess[edpName+"-edp"]
	isPageAvailable := IsPageAvailable(context.Input.Method(), strings.Replace(context.Input.URI(), edpName, "{edpName}", -1), contextRoles)

	if !isPageAvailable {
		http.Error(context.ResponseWriter, "Forbidden.", http.StatusForbidden)
		return
	}
}

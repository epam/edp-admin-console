package filters

import (
	"fmt"
	"github.com/astaxie/beego"
	bgCtx "github.com/astaxie/beego/context"
	"log"
	"net/http"
)

func RoleAccessControlRestFilter(context *bgCtx.Context) {
	log.Println("Start Role Access Control filter..")
	resourceAccess := context.Input.Session("resource_access").(map[string][]string)
	edpTenantName := beego.AppConfig.String("cicdNamespace")
	contextRoles := resourceAccess[edpTenantName+"-edp"]

	if contextRoles == nil {
		nonValidTenantMsg := fmt.Sprintf("Couldn't find tenant by %s name.", edpTenantName)
		http.Error(context.ResponseWriter, nonValidTenantMsg, http.StatusNotFound)
		return
	}

	isPageAvailable := IsPageAvailable(fmt.Sprintf("%s %s", context.Input.Method(), context.Input.URI()), contextRoles)

	if !isPageAvailable {
		http.Error(context.ResponseWriter, "Forbidden.", http.StatusForbidden)
		return
	}
}

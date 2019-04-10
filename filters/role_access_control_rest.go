package filters

import (
	"fmt"
	bgCtx "github.com/astaxie/beego/context"
	"log"
	"net/http"
)

func RoleAccessControlRestFilter(context *bgCtx.Context) {
	log.Println("Start Role Access Control filter..")
	contextRoles := context.Input.Session("realm_roles").([]string)

	isPageAvailable := IsPageAvailable(fmt.Sprintf("%s %s", context.Input.Method(), context.Input.URI()), contextRoles)

	if !isPageAvailable {
		http.Error(context.ResponseWriter, "Forbidden.", http.StatusForbidden)
		return
	}
}

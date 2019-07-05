package filters

import (
	"fmt"
	bgCtx "github.com/astaxie/beego/context"
	"log"
)

func RoleAccessControlFilter(context *bgCtx.Context) {
	log.Println("Start Role Access Control filter..")
	contextRoles := context.Input.Session("realm_roles").([]string)

	isPageAvailable := IsPageAvailable(fmt.Sprintf("%s %s", context.Input.Method(), context.Input.URI()), contextRoles)

	if !isPageAvailable {
		log.Printf("Access to %v is denied", context.Input.URI())
		context.Abort(200, "403")
		return
	}
}

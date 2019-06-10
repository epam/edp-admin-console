package filters

import (
	bgCtx "github.com/astaxie/beego/context"
	"log"
)

func StubAuthFilter(context *bgCtx.Context) {
	log.Println("Start stub auth filter..")
	context.Output.Session("realm_roles", []string{"administrator"})
	context.Output.Session("username", "admin")
}

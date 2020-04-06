package filters

import (
	bgCtx "github.com/astaxie/beego/context"
)

func StubAuthFilter(context *bgCtx.Context) {
	log.Debug("Start stub auth filter..")
	context.Output.Session("realm_roles", []string{"administrator"})
	context.Output.Session("username", "admin")
}

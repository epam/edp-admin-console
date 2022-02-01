package main

import (
	"edp-admin-console/template_function"
	"edp-admin-console/webapi"

	"github.com/astaxie/beego"
)

func main() {
	webapi.SetupRouter()
	template_function.RegisterTemplateFuncs()
	beego.Run()
}

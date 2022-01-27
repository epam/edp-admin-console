package main

import (
	"edp-admin-console/routers"
	"edp-admin-console/template_function"

	"github.com/astaxie/beego"
)

func main() {
	routers.SetupRouter()
	template_function.RegisterTemplateFuncs()
	beego.Run()
}

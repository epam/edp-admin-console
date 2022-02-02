package main

import (
	"log"

	"github.com/astaxie/beego"

	"edp-admin-console/k8s"
	"edp-admin-console/template_function"
	"edp-admin-console/webapi"
)

func main() {
	namespacedClient, err := k8s.SetupNamespacedClient()
	if err != nil {
		log.Fatal(err)
	}
	webapi.SetupRouter(namespacedClient)
	template_function.RegisterTemplateFuncs()
	beego.Run()
}

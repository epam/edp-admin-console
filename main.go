package main

import (
	"log"
	"os"

	"github.com/astaxie/beego"

	"edp-admin-console/k8s"
	"edp-admin-console/template_function"
	"edp-admin-console/webapi"
)

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	namespacedClient, err := k8s.SetupNamespacedClient()
	if err != nil {
		log.Fatal(err)
	}
	webapi.SetupRouter(namespacedClient, workingDir)
	template_function.RegisterTemplateFuncs()
	beego.Run()
}

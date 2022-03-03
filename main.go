package main

import (
	"context"
	"log"
	"os"

	"github.com/astaxie/beego"

	"edp-admin-console/internal/config"
	"edp-admin-console/k8s"
	"edp-admin-console/template_function"
	"edp-admin-console/webapi"
)

func main() {
	conf, err := config.SetupConfig(context.Background(), "conf/app.conf")
	if err != nil {
		log.Fatal(err)
	}
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	namespacedClient, err := k8s.SetupNamespacedClient()
	if err != nil {
		log.Fatal(err)
	}

	webapi.SetupRouter(namespacedClient, workingDir, conf)
	template_function.RegisterTemplateFuncs()
	beego.Run()
}

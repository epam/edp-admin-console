package main

import (
	_ "edp-admin-console/routers"
	"flag"
	"github.com/astaxie/beego"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/spf13/pflag"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func main() {
	pflag.CommandLine.AddFlagSet(zap.FlagSet())
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	logf.SetLogger(zap.Logger())

	beego.Run()
}

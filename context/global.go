package context

import (
	"github.com/astaxie/beego"
	"strings"
)

var (
	Namespace = beego.AppConfig.String("cicdNamespace")
	Tenant    = strings.TrimSuffix(Namespace, "-edp-cicd")
)

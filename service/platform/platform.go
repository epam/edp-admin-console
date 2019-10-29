package platform

import (
	"github.com/astaxie/beego"
)

const (
	platformType      = "platformType"
	platformOpenshift = "openshift"
)

func IsOpenshift() bool {
	return beego.AppConfig.String(platformType) == platformOpenshift
}

package filters

import (
	"fmt"
	"github.com/astaxie/beego"
)

var roles map[string][]string

func init() {
	administrator := beego.AppConfig.String("adminRole")
	developer := beego.AppConfig.String("developerRole")
	roles = map[string][]string{
		"GET /admin/edp/{edpName}/overview":             {administrator, developer},
		"GET /admin/edp/{edpName}/application/overview": {administrator, developer},
		"GET /admin/edp/{edpName}/application/create":   {administrator},
		"POST /admin/edp/{edpName}/application":         {administrator},

		"GET /api/v1/edp/{edpName}/vcs":          {administrator, developer},
		"GET /api/v1/edp/{edpName}":              {administrator, developer},
		"POST /api/v1/edp/{edpName}/application": {administrator},
	}
}

func IsPageAvailable(method string, uri string, contextRoles []string) bool {
	roles := roles[fmt.Sprintf("%s %s", method, uri)]

	if roles == nil || getIntersectionOfRoles(contextRoles, roles) == nil {
		return false
	}
	return true
}

func getIntersectionOfRoles(a, b []string) (c []string) {
	m := make(map[string]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}

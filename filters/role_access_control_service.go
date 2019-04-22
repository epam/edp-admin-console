package filters

import (
	"github.com/astaxie/beego"
	"regexp"
)

var roles map[string][]string

func init() {
	administrator := beego.AppConfig.String("adminRole")
	developer := beego.AppConfig.String("developerRole")
	roles = map[string][]string{
		"GET /admin/edp/overview$":                    {administrator, developer},
		"GET /admin/edp/application/overview":         {administrator, developer},
		"GET /admin/edp/application/create$":          {administrator},
		"GET /admin/edp/application/([^/]*)/overview": {administrator, developer},
		"GET /admin/edp/cd-pipeline/overview$":        {administrator, developer},
		"GET /admin/edp/cd-pipeline/create$":          {administrator, developer},
		"POST /admin/edp/cd-pipeline$":                {administrator},
		"POST /admin/edp/application$":                {administrator},
		"POST /admin/edp/application/([^/]*)/branch$": {administrator},

		"GET /api/v1/edp/vcs$":                 {administrator, developer},
		"GET /api/v1/edp/application$":         {administrator, developer},
		"GET /api/v1/edp/application/([^/]*)$": {administrator, developer},
		"GET /api/v1/edp/cd-pipeline/([^/]*)$": {administrator, developer},
		"POST /api/v1/edp/application$":        {administrator},
	}
}

func IsPageAvailable(key string, contextRoles []string) bool {
	pageRoles := getValue(key)
	if pageRoles == nil {
		return false
	}

	if roles == nil || getIntersectionOfRoles(contextRoles, pageRoles) == nil {
		return false
	}
	return true
}

func getValue(key string) []string {
	for k, v := range roles {
		match, _ := regexp.MatchString(k, key)
		if match {
			return v
		}
	}
	return nil
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

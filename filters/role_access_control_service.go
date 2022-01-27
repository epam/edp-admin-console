package filters

import (
	"regexp"

	"github.com/astaxie/beego"
)

func PermissionsMap() map[string][]string {
	administrator := beego.AppConfig.String("adminRole")
	developer := beego.AppConfig.String("developerRole")
	return map[string][]string{
		"GET /admin/edp/overview$":                          {administrator, developer},
		"GET /admin/edp/application/overview":               {administrator, developer},
		"GET /admin/edp/application/create$":                {administrator},
		"GET /admin/edp/cd-pipeline/overview":               {administrator, developer},
		"GET /admin/edp/cd-pipeline/create$":                {administrator, developer},
		"GET /admin/edp/cd-pipeline/([^/]*)/overview":       {administrator, developer},
		"POST /admin/edp/cd-pipeline$":                      {administrator},
		"POST /admin/edp/application$":                      {administrator},
		"POST /admin/edp/codebase/([^/]*)/branch$":          {administrator},
		"GET /admin/edp/autotest/overview":                  {administrator, developer},
		"GET /admin/edp/autotest/create$":                   {administrator, developer},
		"GET /admin/edp/codebase/([^/]*)/overview":          {administrator, developer},
		"POST /admin/edp/autotest$":                         {administrator},
		"GET /admin/edp/library/overview":                   {administrator, developer},
		"GET /admin/edp/library/create$":                    {administrator, developer},
		"POST /admin/edp/library$":                          {administrator},
		"GET /admin/edp/cd-pipeline/([^/]*)/update":         {administrator, developer},
		"POST /admin/edp/cd-pipeline/([^/]*)/update":        {administrator},
		"POST /admin/edp/codebase$":                         {administrator},
		"POST /admin/edp/stage$":                            {administrator},
		"POST /admin/edp/cd-pipeline/delete":                {administrator},
		"GET /admin/edp/diagram/overview":                   {administrator, developer},
		"POST /admin/edp/cd-pipeline/([^/]*)/cd-stage/edit": {administrator},

		"GET /api/v1/edp/vcs$":                               {administrator, developer},
		"GET /api/v1/edp/codebase":                           {administrator, developer},
		"GET /api/v1/edp/codebase/([^/]*)$":                  {administrator, developer},
		"GET /api/v1/edp/cd-pipeline/([^/]*)$":               {administrator, developer},
		"GET /api/v1/edp/cd-pipeline/([^/]*)/stage/([^/]*)$": {administrator, developer},
		"POST /api/v1/edp/codebase$":                         {administrator},
		"POST /api/v1/edp/cd-pipeline$":                      {administrator},
		"PUT /api/v1/edp/cd-pipeline/([^/]*)$":               {administrator},
		"DELETE /api/v1/edp/codebase$":                       {administrator},
		"DELETE /api/v1/edp/stage$":                          {administrator},
	}
}

func IsPageAvailable(permissions map[string][]string, key string, contextRoles []string) (bool, error) {
	pageRoles, err := getValue(permissions, key)
	if err != nil {
		return false, err
	}
	if pageRoles == nil {
		return true, nil
	}

	if getIntersectionOfRoles(contextRoles, pageRoles) == nil {
		return false, nil
	}
	return true, nil
}

func getValue(permissions map[string][]string, key string) ([]string, error) {
	for k, v := range permissions {
		match, err := regexp.MatchString(k, key)
		if err != nil {
			return nil, err
		}
		if match {
			return v, nil
		}
	}
	return nil, nil
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

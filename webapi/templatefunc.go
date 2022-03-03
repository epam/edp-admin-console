package webapi

import (
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/pkg/errors"

	"edp-admin-console/util"
)

func add(a, b int) int {
	return a + b
}

func getDefaultBranchVersion(cb []codebaseBranch, defaultBranch string) string {
	for _, branch := range cb {
		if branch.Name == defaultBranch {
			v := branch.Version
			return util.TrimSuffix(*v, "-SNAPSHOT")
		}
	}
	return ""
}

func incrementVersion(v string) (*string, error) {
	pv, err := semver.Make(v)
	if err != nil {
		return nil, err
	}

	pv.Minor += 1

	res := util.TrimSuffix(pv.String(), "-SNAPSHOT")
	return &res, nil
}

func compareJiraServer(codebaseJiraServer *string, jiraServer string) bool {
	if codebaseJiraServer == nil {
		return false
	}
	return *codebaseJiraServer == jiraServer
}

func params(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid params call")
	}
	p := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		k, ok := values[i].(string)
		if !ok {
			return nil, errors.New("params keys must be strings")
		}
		p[k] = values[i+1]
	}
	return p, nil
}

func getCurrentYear() int {
	return time.Now().Year()
}

func CapitalizeFirstLetter(s string) string {
	return strings.Title(s)
}

func CapitalizeAll(s string) string {
	return strings.ToUpper(s)
}

func LowercaseAll(s string) string {
	return strings.ToLower(s)
}

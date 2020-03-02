package template_function

import (
	"github.com/astaxie/beego"
	"github.com/blang/semver"
	"github.com/pkg/errors"
	"strings"
)

func init() {
	if err := beego.AddFuncMap("add", add); err != nil {
		panic("couldn't register 'add' function to go template")
	}
	if err := beego.AddFuncMap("params", params); err != nil {
		panic("couldn't register 'params' function to go template")
	}
	if err := beego.AddFuncMap("trimSuffix", trimSuffix); err != nil {
		panic("couldn't register 'trimSuffix' function to go template")
	}
	if err := beego.AddFuncMap("incrementVersion", incrementVersion); err != nil {
		panic("couldn't register 'incrementVersion' function to go template")
	}
}

func trimSuffix(v, s string) string {
	return strings.TrimSuffix(v, s)
}

func incrementVersion(v string) (*string, error) {
	pv, err := semver.Make(v)
	if err != nil {
		return nil, err
	}

	pv.Minor += 1

	res := trimSuffix(pv.String(), "-SNAPSHOT")
	return &res, nil
}

func add(a, b int) int {
	return a + b
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

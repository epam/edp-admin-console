package util

import (
	"bytes"
	"edp-admin-console/models/query"
	"edp-admin-console/service/logger"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"go.uber.org/zap"
)

var (
	log                 = logger.GetLogger()
	imageStreamReplacer = strings.NewReplacer("/", "-", ".", "-")
)

func GetValuesFromConfig(name string) []string {
	values := beego.AppConfig.String(name)
	if values == "" {
		log.Debug("variable is empty.", zap.String("name", name))
		return nil
	}

	s := strings.Split(values, ",")
	log.Info("Fetched data variable",
		zap.String("name", name), zap.Strings("values", s))
	return s
}

func GetVersionOrNil(value, postfix string) *string {
	if value == "" {
		return nil
	}

	if postfix == "" {
		v := value
		return &v
	}

	v := fmt.Sprintf("%v-%v", value, postfix)

	return &v
}

func EncodeStructToBytes(s interface{}) ([]byte, error) {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(s)
	if err != nil {
		return nil, err
	}
	return reqBodyBytes.Bytes(), nil
}

func TrimSuffix(v, s string) string {
	return strings.TrimSuffix(v, s)
}

func ProcessNameToKubernetesConvention(name string) string {
	return strings.Replace(name, "/", "-", -1)
}

func ProcessCodeBaseImageStreamNameConvention(name string) string {
	return imageStreamReplacer.Replace(name)
}

func GetStringP(val string) *string {
	return &val
}

func GetCodebaseTypeP(val query.CodebaseType) *query.CodebaseType {
	return &val
}

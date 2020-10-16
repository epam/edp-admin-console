package perfboard

import (
	"edp-admin-console/models/query"
	"github.com/astaxie/beego/orm"
)

type IPerfServer interface {
	GetPerfServers() ([]*query.PerfServer, error)
}

type PerfServer struct {
}

func (PerfServer) GetPerfServers() ([]*query.PerfServer, error) {
	o := orm.NewOrm()
	var servers []*query.PerfServer
	_, err := o.QueryTable(new(query.PerfServer)).
		OrderBy("name").
		All(&servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

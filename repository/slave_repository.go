package repository

import (
	"edp-admin-console/models/query"
	"github.com/astaxie/beego/orm"
)

type ISlaveRepository interface {
	GetAllSlaves() ([]*query.JenkinsSlave, error)
}

type SlaveRepository struct {
}

func (s SlaveRepository) GetAllSlaves() ([]*query.JenkinsSlave, error) {
	o := orm.NewOrm()
	var jenkinsSlaves []*query.JenkinsSlave

	_, err := o.QueryTable(new(query.JenkinsSlave)).
		OrderBy("name").
		All(&jenkinsSlaves)
	if err != nil {
		return nil, err
	}

	return jenkinsSlaves, nil
}

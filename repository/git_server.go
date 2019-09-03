package repository

import (
	"edp-admin-console/models/query"
	"github.com/astaxie/beego/orm"
)

type IGitServerRepository interface {
	GetGitServersByCriteria(criteria query.GitServerCriteria) ([]*query.GitServer, error)
}

type GitServerRepository struct {
	IGitServerRepository
}

func (GitServerRepository) GetGitServersByCriteria(criteria query.GitServerCriteria) ([]*query.GitServer, error) {
	o := orm.NewOrm()
	var gitServers []*query.GitServer

	_, err := o.QueryTable(new(query.GitServer)).
		Filter("available", criteria.Available).
		OrderBy("name").
		All(&gitServers)
	if err != nil {
		return nil, err
	}

	return gitServers, nil
}

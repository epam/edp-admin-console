package repository

import (
	"edp-admin-console/models/query"
	"github.com/astaxie/beego/orm"
)

type IGitServerRepository interface {
	GetGitServersByCriteria(criteria query.GitServerCriteria) ([]*query.GitServer, error)
	GetGitServerByName(name string) (*query.GitServer, error)
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

func (GitServerRepository) GetGitServerByName(name string) (*query.GitServer, error) {
	o := orm.NewOrm()
	gitServer := query.GitServer{Name: name}

	err := o.Read(&gitServer, "Name")

	if err == orm.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &gitServer, nil
}

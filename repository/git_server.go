package repository

import (
	"edp-admin-console/models/query"
)

type IGitServerRepository interface {
	GetGitServersByCriteria(criteria query.GitServerCriteria) ([]query.GitServer, error)
}

type GitServerRepository struct {
	IGitServerRepository
}

func (GitServerRepository) GetGitServersByCriteria(criteria query.GitServerCriteria) ([]query.GitServer, error) {
	return []query.GitServer{
		{
			Id:     1,
			Name:   "https://git.epam.com",
			Status: "active",
		},
		{
			Id:     2,
			Name:   "https://gerrit-edp-cicd-delivery.delivery.aws.main.edp.projects.epam.com",
			Status: "active",
		},
	}, nil
}

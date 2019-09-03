package service

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository"
	"log"
)

type GitServerService struct {
	IGitServerRepository repository.IGitServerRepository
}

func (s GitServerService) GetServers(criteria query.GitServerCriteria) ([]*query.GitServer, error) {
	log.Println("Start fetching Git Servers...")

	gitServers, err := s.IGitServerRepository.GetGitServersByCriteria(criteria)
	if err != nil {
		log.Printf("An error has occurred while fetching Git Servers from DB: %v", err)
		return nil, err
	}
	log.Printf("Fetched Git Servers. Count: %v. Values: %v", len(gitServers), gitServers)

	return gitServers, nil
}

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

func (s GitServerService) GetGitServer(name string) (*query.GitServer, error) {
	log.Printf("Start fetching Git Server %v...", name)

	g, err := s.IGitServerRepository.GetGitServerByName(name)
	if err != nil {
		log.Printf("An error has occurred while fetching Git Server %v from DB: %v", name, err)
		return nil, err
	}
	log.Printf("Fetched Git Server: %v ", g)

	return g, nil
}

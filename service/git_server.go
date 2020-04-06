package service

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository"
	"go.uber.org/zap"
)

type GitServerService struct {
	IGitServerRepository repository.IGitServerRepository
}

func (s GitServerService) GetServers(criteria query.GitServerCriteria) ([]*query.GitServer, error) {
	log.Debug("Start fetching Git Servers...")

	gitServers, err := s.IGitServerRepository.GetGitServersByCriteria(criteria)
	if err != nil {
		log.Error("An error has occurred while fetching Git Servers from DB", zap.Error(err))
		return nil, err
	}
	log.Info("Fetched Git Servers",
		zap.Int("count", len(gitServers)), zap.Any("git servers", gitServers))

	return gitServers, nil
}

func (s GitServerService) GetGitServer(name string) (*query.GitServer, error) {
	log.Debug("Start fetching Git Server...", zap.String("name", name))

	g, err := s.IGitServerRepository.GetGitServerByName(name)
	if err != nil {
		log.Error("An error has occurred while fetching Git Server from DB",
			zap.String("name", name), zap.Error(err))
		return nil, err
	}
	log.Info("Fetched Git Server", zap.Any("git server", g))
	return g, nil
}

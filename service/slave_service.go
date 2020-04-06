package service

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type SlaveService struct {
	ISlaveRepository repository.ISlaveRepository
}

//GetAllSlaves gets all slave entities from DB
func (s SlaveService) GetAllSlaves() ([]*query.JenkinsSlave, error) {
	log.Debug("Start fetching all available Slaves...")
	jenkinsSlaves, err := s.ISlaveRepository.GetAllSlaves()
	if err != nil {
		return nil, errors.Wrap(err, "an error has occurred while fetching slave entities from DB")
	}
	log.Info("Fetched Slaves", zap.Any("slaves", jenkinsSlaves))
	return jenkinsSlaves, nil
}

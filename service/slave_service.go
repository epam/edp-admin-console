package service

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository"
	"github.com/pkg/errors"
	"log"
)

type SlaveService struct {
	ISlaveRepository repository.ISlaveRepository
}

//GetAllSlaves gets all slave entities from DB
func (s SlaveService) GetAllSlaves() ([]*query.JenkinsSlave, error) {
	log.Println("Start fetching all available Slaves...")

	jenkinsSlaves, err := s.ISlaveRepository.GetAllSlaves()
	if err != nil {
		return nil, errors.Wrap(err, "an error has occurred while fetching slave entities from DB")
	}

	log.Printf("Fetched Slaves: %v", jenkinsSlaves)

	return jenkinsSlaves, nil
}

package service

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository"
	"go.uber.org/zap"
)

type ThirdPartyService struct {
	IServiceCatalogRepository repository.IServiceCatalogRepository
}

func (s ThirdPartyService) GetAllServices() ([]query.ThirdPartyService, error) {
	log.Debug("Start execution of GetAllServices method...")
	services, err := s.IServiceCatalogRepository.GetAllServices()
	if err != nil {
		log.Error("An error has occurred while getting services from database", zap.Error(err))
		return nil, err
	}
	log.Info("Fetched services",
		zap.Int("count", len(services)), zap.Any("services", services))
	return services, nil
}

package service

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository"
	"log"
)

type ThirdPartyService struct {
	IServiceCatalogRepository repository.IServiceCatalogRepository
}

func (s ThirdPartyService) GetAllServices() ([]query.ThirdPartyService, error) {
	log.Println("Start execution of GetAllServices method...")
	services, err := s.IServiceCatalogRepository.GetAllServices()
	if err != nil {
		log.Printf("An error has occurred while getting services from database: %s", err)
		return nil, err
	}
	log.Printf("Fetched services. Count: %v. Rows: %v", len(services), services)

	return services, nil
}

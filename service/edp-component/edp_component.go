package edp_component

import (
	"edp-admin-console/context"
	"edp-admin-console/models/query"
	ec "edp-admin-console/repository/edp-component"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

type EDPComponentService struct {
	IEDPComponent ec.IEDPComponentRepository
}

//GetEDPComponent gets EDP component by type from DB
func (s EDPComponentService) GetEDPComponent(componentType string) (*query.EDPComponent, error) {
	log.Printf("Start fetching EDP Component by %v type...", componentType)

	c, err := s.IEDPComponent.GetEDPComponent(componentType)
	if err != nil {
		return nil, errors.Wrapf(err, "an error has occurred while fetching EDP Component by %v type from DB",
			componentType)
	}

	log.Printf("Fetched EDP Component: %+v", c)

	return c, nil
}

//GetEDPComponents gets all EDP components from DB
func (s EDPComponentService) GetEDPComponents() ([]*query.EDPComponent, error) {
	log.Println("Start fetching EDP Components...")

	c, err := s.IEDPComponent.GetEDPComponents()
	if err != nil {
		return nil, errors.Wrap(err, "an error has occurred while fetching EDP Components from DB")
	}
	log.Printf("Fetched EDP Components: %+v", c)

	for i, v := range c {
		if v.Type == "openshift" {
			c[i].Url = fmt.Sprintf("%s/console/project/%s-edp-cicd/overview", v.Url, context.Tenant)
		}
	}

	return c, nil
}

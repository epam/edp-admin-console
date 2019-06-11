package repository

import (
	"edp-admin-console/models/query"
	"github.com/astaxie/beego/orm"
)

type IServiceCatalogRepository interface {
	GetAllServices() ([]query.ThirdPartyService, error)
}

type ServiceCatalogRepository struct {
	IServiceCatalogRepository
}

func (ServiceCatalogRepository) GetAllServices() ([]query.ThirdPartyService, error) {
	o := orm.NewOrm()
	var services []query.ThirdPartyService

	_, err := o.QueryTable(new(query.ThirdPartyService)).All(&services)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return services, nil
}

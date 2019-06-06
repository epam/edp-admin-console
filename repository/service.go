package repository

import (
	"edp-admin-console/models/query"
	"github.com/astaxie/beego/orm"
)

type IServiceCatalogRepository interface {
	GetAllServices() ([]query.Service, error)
}

type ServiceCatalogRepository struct {
	IServiceCatalogRepository
}

func (ServiceCatalogRepository) GetAllServices() ([]query.Service, error) {
	o := orm.NewOrm()
	var services []query.Service

	_, err := o.QueryTable(new(query.Service)).All(&services)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return services, nil
}

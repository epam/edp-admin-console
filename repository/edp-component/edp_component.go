package edp_component

import (
	"edp-admin-console/models/query"
	dberror "edp-admin-console/util/error/db-errors"
	"github.com/astaxie/beego/orm"
)

type IEDPComponentRepository interface {
	GetEDPComponent(componentType string) (*query.EDPComponent, error)
	GetEDPComponents() ([]*query.EDPComponent, error)
}

type EDPComponent struct {
}

func (EDPComponent) GetEDPComponent(componentType string) (*query.EDPComponent, error) {
	o := orm.NewOrm()
	c := query.EDPComponent{}

	err := o.QueryTable(new(query.EDPComponent)).
		Filter("type", componentType).
		One(&c)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, dberror.StatusError{
				Status: dberror.StatusReasonNotFound,
			}
		}
		return nil, err
	}

	return &c, nil
}

func (EDPComponent) GetEDPComponents() ([]*query.EDPComponent, error) {
	o := orm.NewOrm()
	var c []*query.EDPComponent

	_, err := o.QueryTable(new(query.EDPComponent)).
		OrderBy("type").
		All(&c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

package repository

import (
	"edp-admin-console/models"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
)

type EDPTenantRep struct {
}

func (this EDPTenantRep) GetAllEDPTenantsByNames(adminClients []string) ([]*models.EDPTenant, error) {
	var edpTenants []*models.EDPTenant
	o := orm.NewOrm()
	_, err := o.QueryTable(new(models.EDPTenant)).Filter("name__in", adminClients).All(&edpTenants, "name", "version")
	if err != nil {
		return nil, err
	}
	return edpTenants, nil
}

func (this EDPTenantRep) GetEdpVersionByName(edpName string) (string, error) {
	var edp models.EDPTenant
	o := orm.NewOrm()
	err := o.QueryTable(new(models.EDPTenant)).Filter("name__istartswith", edpName).One(&edp, "name", "version")
	if err == orm.ErrMultiRows {
		return "", errors.New("The problem has been detected due to db contains more than one value.")
	}
	if err == orm.ErrNoRows {
		return "", errors.New(fmt.Sprintf("Couldn't find EDP version for %s.", edpName))
	}
	return edp.Version, nil
}

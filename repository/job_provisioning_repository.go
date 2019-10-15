package repository

import (
	"edp-admin-console/models/query"
	"github.com/astaxie/beego/orm"
)

type IJobProvisioningRepository interface {
	GetAllJobProvisioners() ([]*query.JobProvisioning, error)
}

type JobProvisioning struct {
}

func (JobProvisioning) GetAllJobProvisioners() ([]*query.JobProvisioning, error) {
	o := orm.NewOrm()
	var jobsProvisioning []*query.JobProvisioning

	_, err := o.QueryTable(new(query.JobProvisioning)).
		OrderBy("name").
		All(&jobsProvisioning)
	if err != nil {
		return nil, err
	}

	return jobsProvisioning, nil
}

package query

type JobProvisioning struct {
	Id   int    `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
}

func (c *JobProvisioning) TableName() string {
	return "job_provisioning"
}

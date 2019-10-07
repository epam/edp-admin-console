package query

type JenkinsSlave struct {
	Id   int    `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
}

func (c *JenkinsSlave) TableName() string {
	return "jenkins_slave"
}

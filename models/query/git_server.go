package query

type GitServer struct {
	Id     int    `json:"id" orm:"column(id)"`
	Name   string `json:"name" orm:"column(name)"`
	Status Status `json:"status" orm:"column(status)"`
}

func (c *GitServer) TableName() string {
	return "git_server"
}

type GitServerCriteria struct {
	Status Status
}

package query

type GitServer struct {
	Id        int    `json:"id" orm:"column(id)"`
	Name      string `json:"name" orm:"column(name)"`
	Hostname  string `json:"hostname" orm:"column(hostname)"`
	Available bool   `json:"available" orm:"column(available)"`
}

func (c *GitServer) TableName() string {
	return "git_server"
}

type GitServerCriteria struct {
	Available bool
}

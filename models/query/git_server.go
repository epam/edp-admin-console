package query

type GitServer struct {
	Id        int    `json:"id" orm:"column(id)"`
	Name      string `json:"name" orm:"column(name)"`
	Available bool   `json:"available" orm:"column(available)"`
}

func (c *GitServer) TableName() string {
	return "git_server"
}

type GitServerCriteria struct {
	Available bool
}

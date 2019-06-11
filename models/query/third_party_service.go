package query

type ThirdPartyService struct {
	Id          int    `json:"id" orm:"column(id)"`
	Name        string `json:"name" orm:"column(name)"`
	Description string `json:"description" orm:"column(description)"`
	Version     string `json:"version" orm:"column(version)"`
}

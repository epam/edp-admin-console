package query

type ThirdPartyService struct {
	Id          int    `json:"id" orm:"column(id)"`
	Name        string `json:"name" orm:"column(name)"`
	Description string `json:"description" orm:"column(description)"`
	Version     string `json:"version" orm:"column(version)"`
	Url         string `json:"url" orm:"column(url)"`
	Icon        string `json:"-" orm:"column(icon)"`
}

func (cb *ThirdPartyService) TableName() string {
	return "third_party_service"
}

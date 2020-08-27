package query

type EDPComponent struct {
	Id      int    `json:"id" orm:"column(id)"`
	Type    string `json:"type" orm:"column(type)"`
	Url     string `json:"url" orm:"column(url)"`
	Icon    string `json:"icon" orm:"column(icon)"`
	Visible bool   `json:"visible" orm:"column(visible)"`
}

func (c *EDPComponent) TableName() string {
	return "edp_component"
}

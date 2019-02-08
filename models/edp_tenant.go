package models

type EDPTenant struct {
	Name    string `json:"name" orm:"pk"`
	Version string `json:"version"`
}

func (u *EDPTenant) TableName() string {
	return "edp_specification"
}

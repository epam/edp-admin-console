package models

type AutotestView struct {
	Name      string `json:"name"`
	Language  string `json:"language"`
	BuildTool string `json:"build_tool"`
	Status    string `json:"status"`
}

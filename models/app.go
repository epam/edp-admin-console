package models

type App struct {
	Name      string   `json:"name" valid:"Required"`
	Strategy  string   `json:"strategy"`
	Lang      string   `json:"lang" valid:"Required"`
	BuildTool string   `json:"buildTool" valid:"Required"`
	Framework string   `json:"framework" valid:"Required"`
	Git       string   `json:"git" valid:"Required; MinSize(1)"`
	Route     Route    `json:"route"`
	Database  Database `json:"database"`
}

type Route struct {
	Site string `json:"site"`
	Path string `json:"path"`
}

type Database struct {
	Kind     string `json:"kind"`
	Version  string `json:"version"`
	Capacity string `json:"capacity"`
	Storage  string `json:"storage"`
}

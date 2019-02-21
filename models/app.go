package models

type App struct {
	Name      string    `json:"name" valid:"Required"`
	Strategy  string    `json:"strategy" valid:"Required"`
	Lang      string    `json:"lang" valid:"Required"`
	BuildTool string    `json:"buildTool" valid:"Required"`
	Framework string    `json:"framework" valid:"Required"`
	Git       *Git       `json:"git,omitempty"`
	Route     *Route    `json:"route,omitempty"`
	Database  *Database `json:"database,omitempty"`
}

type Git struct {
	Url      string `json:"url,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type Route struct {
	Site string `json:"site,omitempty"`
	Path string `json:"path,omitempty"`
}

type Database struct {
	Kind     string `json:"kind,omitempty"`
	Version  string `json:"version,omitempty"`
	Capacity string `json:"capacity,omitempty"`
	Storage  string `json:"storage,omitempty"`
}

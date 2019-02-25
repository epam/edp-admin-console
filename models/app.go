package models

type App struct {
	Name       string      `json:"name" valid:"Required;Match(/^[a-z][a-z0-9-.]+[a-z]$/)"`
	Strategy   string      `json:"strategy" valid:"Required"`
	Lang       string      `json:"lang" valid:"Required"`
	BuildTool  string      `json:"buildTool" valid:"Required"`
	Framework  string      `json:"framework" valid:"Required"`
	Repository *Repository `json:"repository,omitempty"`
	Route      *Route      `json:"route,omitempty"`
	Database   *Database   `json:"database,omitempty"`
	Vcs        *Vcs        `json:"vcs,omitempty"`
}

type Repository struct {
	Url      string `json:"url,omitempty" valid:"Match(/(?:^git|^ssh|^https?|^git@[-\\w.]+):(\\/\\/)?(.*?)(\\.git)(\\/?|\\#[-\\d\\w._]+?)$/)"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type Vcs struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type Route struct {
	Site string `json:"site,omitempty"`
	Path string `json:"path,omitempty" valid:"Match(/^(?:http(s)?:\\/\\/)?[\\w.-]+(?:\\.[\\w\\.-]+)+[\\w\\-\\._~:/?#[\\]@!\\$&'\\(\\)\\*\\+,;=.]+$/)"`
}

type Database struct {
	Kind     string `json:"kind,omitempty"`
	Version  string `json:"version,omitempty"`
	Capacity string `json:"capacity,omitempty"`
	Storage  string `json:"storage,omitempty"`
}

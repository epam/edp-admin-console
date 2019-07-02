package query

type Codebase struct {
	Id                  int               `json:"id" orm:"column(id)"`
	Name                string            `json:"name" orm:"column(name)"`
	Language            string            `json:"language" orm:"column(language)"`
	BuildTool           string            `json:"build_tool" orm:"column(build_tool)"`
	Framework           string            `json:"framework" orm:"column(framework)"`
	Strategy            string            `json:"strategy" orm:"column(strategy)"`
	GitUrl              string            `json:"git_url" orm:"column(repository_url)"`
	RouteSite           string            `json:"route_site" orm:"column(route_site)"`
	RoutePath           string            `json:"route_path" orm:"column(route_path)"`
	DbKind              string            `json:"db_kind" orm:"column(database_kind)"`
	DbVersion           string            `json:"db_version" orm:"column(database_version)"`
	DbCapacity          string            `json:"db_capacity" orm:"column(database_capacity)"`
	DbStorage           string            `json:"db_storage" orm:"column(database_storage)"`
	Type                CodebaseType      `json:"type" orm:"column(type)"`
	Status              Status            `json:"status" orm:"column(status)"`
	TestReportFramework string            `json:"testReportFramework" orm:"column(test_report_framework)"`
	Description         string            `json:"description" orm:"column(description)"`
	CodebaseBranch      []*CodebaseBranch `json:"codebase_branch" orm:"reverse(many)"`
	ActionLog           []*ActionLog      `json:"action_log" orm:"rel(m2m);rel_table(codebase_action_log)"`
}

func (c *Codebase) TableName() string {
	return "codebase"
}

type CodebaseCriteria struct {
	BranchStatus Status
	Status       Status
	Type         CodebaseType
}

type CodebaseType string

const (
	App       CodebaseType = "application"
	Autotests CodebaseType = "autotests"
	Library   CodebaseType = "library"
)

var CodebaseTypes = map[string]CodebaseType{
	"application": App,
	"autotests":   Autotests,
	"library":     Library,
}

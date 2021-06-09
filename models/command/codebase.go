package command

type CreateCodebase struct {
	Name                     string      `json:"name" valid:"Required;Match(/^[a-z][a-z0-9-]*[a-z0-9]$/)"`
	DefaultBranch            string      `json:"defaultBranch" valid:"Required"`
	Strategy                 string      `json:"strategy"`
	Lang                     string      `json:"lang" valid:"Required"`
	Framework                *string     `json:"framework,omitempty"`
	BuildTool                string      `json:"buildTool" valid:"Required"`
	TestReportFramework      *string     `json:"testReportFramework"`
	MultiModule              bool        `json:"multiModule,omitempty"`
	Type                     string      `json:"type,omitempty" valid:"Required"`
	Repository               *Repository `json:"repository,omitempty"`
	Vcs                      *Vcs        `json:"vcs,omitempty"`
	Description              *string     `json:"description,omitempty"`
	Username                 string      `json:"username"`
	GitServer                string      `json:"gitServer"`
	Versioning               Versioning  `json:"versioning"`
	GitUrlPath               *string     `json:"gitUrlPath"`
	JenkinsSlave             *string     `json:"jenkinsSlave,omitempty"`
	JobProvisioning          *string     `json:"jobProvisioning,omitempty"`
	DeploymentScript         string      `json:"deploymentScript"`
	JiraServer               *string     `json:"jiraServer,omitempty"`
	CommitMessageRegex       *string     `json:"commitMessagePattern"`
	TicketNameRegex          *string     `json:"ticketNamePattern"`
	CiTool                   string      `json:"ciTool" valid:"Required"`
	Perf                     *Perf       `json:"perf,omitempty"`
	JiraIssueMetadataPayload *string     `json:"jiraIssueMetadataPayload"`
	EmptyProject             bool        `json:"emptyProject" valid:"Required"`
}

type Perf struct {
	Name        string   `json:"name" valid:"Required"`
	DataSources []string `json:"dataSources" valid:"Required"`
}

type Repository struct {
	Url      string `json:"url,omitempty" valid:"Required;Match(/(?:^git|^ssh|^https?|^git@[-\\w.]+):(\\/\\/)?(.*?)(\\.git)(\\/?|\\#[-\\d\\w._]+?)$/)"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type Vcs struct {
	Login    string `json:"login,omitempty" valid:"Required"`
	Password string `json:"password,omitempty" valid:"Required"`
}

type DeleteCodebaseCommand struct {
	Name string `json:"name"`
}

type Versioning struct {
	Type      string  `json:"type" valid:"Required"`
	StartFrom *string `json:"startFrom,omitempty"`
}

type UpdateCodebaseCommand struct {
	Name                     string
	CommitMessageRegex       string
	TicketNameRegex          string
	JiraIssueMetadataPayload *string
	JiraServer               string
}

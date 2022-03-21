package consts

const (
	//EDP Components
	Jenkins        = "jenkins"
	Openshift      = "openshift"
	Kubernetes     = "kubernetes"
	DockerRegistry = "docker-registry"
	Gerrit         = "gerrit"

	//Codebase types
	Application = "application"
	Autotest    = "autotests"
	Library     = "library"

	//Kinds
	CodebasePlural       = "codebases"
	CodebaseBranchPlural = "codebasebranches"
	StagePlural          = "stages"
	CDPipelinePlural     = "cdpipelines"
	CodebaseKind         = "Codebase"
	Overview             = "overview"

	ImportStrategy        = "import"
	CloneStrategy         = "clone"
	LanguageJava          = "Java"
	DefaultVersioningType = "default"
	JenkinsCITool         = "Jenkins"

	InitializedStatus            = "initialized"
	CdPipelineRegistrationAction = "cd_pipeline_registration"
	SuccessResult                = "success"
	InactiveValue                = "inactive"
	ActiveValue                  = "active"

	IssuesLinksKey = "issuesLinks"
	CDMenuItem     = "delivery"
)

var DefaultBuildNumber = "0"

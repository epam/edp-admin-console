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

	//Kinds
	CodebasePlural       = "codebases"
	CodebaseBranchPlural = "codebasebranches"
	StagePlural          = "stages"
	CDPipelinePlural     = "cdpipelines"
	CodebaseKind         = "Codebase"

	ImportStrategy        = "import"
	LanguageJava          = "Java"
	DefaultVersioningType = "default"
	JenkinsCITool         = "Jenkins"

	InitializedStatus            = "initialized"
	CdPipelineRegistrationAction = "cd_pipeline_registration"
	SuccessResult                = "success"
	InactiveValue                = "inactive"
)

var DefaultBuildNumber = "0"

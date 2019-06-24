package models

type StageView struct {
	Name            string             `json:"name"`
	CDPipeline      string             `json:"cdPipeline"`
	Description     string             `json:"description"`
	QualityGate     string             `json:"qualityGate"`
	TriggerType     string             `json:"triggerType"`
	Order           string             `json:"order"`
	JenkinsStepName string             `json:"jenkinsStepName"`
	Applications    []ApplicationStage `json:"applications"`
	Autotests       []Autotests        `json:"autotests"`
}

type ApplicationStage struct {
	Name       string `json:"name"`
	BranchName string `json:"branchName"`
	InputIs    string `json:"inputIs"`
	OutputIs   string `json:"outputIs"`
}

type CDPipelineStageView struct {
	Name                 string `json:"name"`
	Description          string `json:"description"`
	TriggerType          string `json:"triggerType"`
	QualityGate          string `json:"qualityGate"`
	JenkinsStepName      string `json:"jenkinsStepName"`
	Order                int    `json:"order"`
	OpenshiftProjectLink string `json:"openshiftProjectLink"`
}

type Autotests struct {
	Name       string `json:"name"`
}

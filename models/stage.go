package models

type StageView struct {
	Name         string             `json:"name"`
	CDPipeline   string             `json:"cdPipeline"`
	Description  string             `json:"description"`
	QualityGate  string             `json:"qualityGate"`
	TriggerType  string             `json:"triggerType"`
	Applications []ApplicationStage `json:"applications"`
}

type ApplicationStage struct {
	Name     string `json:"name"`
	InputIs  string `json:"inputIs"`
	OutputIs string `json:"outputIs"`
}

type CDPipelineStageView struct {
	Name                 string `json:"name"`
	Description          string `json:"description"`
	TriggerType          string `json:"triggerType"`
	QualityGate          string `json:"qualityGate"`
	JenkinsStepName      string `json:"jenkinsStepName"`
	OpenshiftProjectLink string `json:"openshiftProjectLink"`
}

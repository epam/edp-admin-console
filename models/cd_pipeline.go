package models

type CDPipelineDTO struct {
	Name             string              `json:"name"`
	Status           string              `json:"status"`
	CodebaseBranches []CodebaseBranchDTO `json:"codebaseBranches"`
}

type CDPipelineView struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	JenkinsLink string
}

type StageCreate struct {
	Name            string `json:"name" valid:"Required;Match(/^[a-z0-9]([-a-z0-9]*[a-z0-9])$/)"`
	Description     string `json:"description" valid:"Required"`
	StepName        string `json:"stepName" valid:"Required;Match(/^[a-z0-9]([-a-z0-9]*[a-z0-9])$/)"`
	QualityGateType string `json:"qualityGateType" valid:"Required"`
	TriggerType     string `json:"triggerType" valid:"Required"`
	Order           int    `json:"order"`
}

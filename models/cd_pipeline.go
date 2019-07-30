package models

import "edp-admin-console/models/dto"

type CDPipelineDTO struct {
	Name                string                                                  `json:"name"`
	Status              string                                                  `json:"status"`
	JenkinsLink         string                                                  `json:"jenkinsLink"`
	CodebaseBranches    []dto.CodebaseBranchDTO                                 `json:"codebaseBranches"`
	Stages              []CDPipelineStageView                                   `json:"stages"`
	CodebaseStageMatrix map[CDCodebaseStageMatrixKey]CDCodebaseStageMatrixValue `json:"codebaseStageMatrix"`
}

type CDCodebaseStageMatrixKey struct {
	CodebaseBranch dto.CodebaseBranchDTO `json:"codebaseBranch"`
	Stage          CDPipelineStageView   `json:"stage"`
}

type CDCodebaseStageMatrixValue struct {
	DockerVersion string `json:"dockerVersion"`
}

func (c *CDPipelineDTO) getCDCodebaseStageMatrixValue(codebase dto.CodebaseBranchDTO, stage CDPipelineStageView) CDCodebaseStageMatrixValue {
	return c.CodebaseStageMatrix[CDCodebaseStageMatrixKey{
		CodebaseBranch: codebase,
		Stage:          stage,
	}]
}

type CDPipelineView struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	JenkinsLink string `json:"jenkinsLink"`
}

type StageCreate struct {
	Name         string        `json:"name" valid:"Required;Match(/^[a-z0-9]([-a-z0-9]*[a-z0-9])$/)"`
	Description  string        `json:"description" valid:"Required"`
	TriggerType  string        `json:"triggerType" valid:"Required"`
	Order        int           `json:"order" valid:"Match(/^[0-9]$/)"`
	QualityGates []QualityGate `json:"qualityGates" valid:"Required"`
}

type QualityGate struct {
	QualityGateType string  `json:"qualityGateType" valid:"Required"`
	StepName        string  `json:"stepName" valid:"Required;Match(/^[A-z0-9-._]/)"`
	AutotestName    *string `json:"autotestName"`
	BranchName      *string `json:"branchName"`
}

type Autotest struct {
	Name       string `json:"autotestName"`
	BranchName string `json:"branchName"`
}

type CDPipelineCommand struct {
	Name                 string                             `json:"name" valid:"Required;Match(/^[a-z0-9]([-a-z0-9]*[a-z0-9])$/)"`
	Applications         []ApplicationWithInputDockerStream `json:"applications" valid:"Required"`
	ThirdPartyServices   []string                           `json:"services"`
	Stages               []StageCreate                      `json:"stages"`
	ApplicationToApprove []string                           `json:"-"`
}

type ApplicationWithInputDockerStream struct {
	ApplicationName   string `json:"appName" valid:"Required;Match(/^[a-z][a-z0-9-]*[a-z0-9]$/)"`
	BranchName        string `json:"branchName" valid:"Required;Match(/^[a-z0-9][a-z0-9-._]*[a-z0-9]$/)"`
	InputDockerStream string `json:"inputDockerStream"`
}

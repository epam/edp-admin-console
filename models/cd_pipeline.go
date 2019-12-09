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

type CDPipelineApplicationCommand struct {
	ApplicationName   string `json:"appName" valid:"Required;Match(/^[a-z][a-z0-9-]*[a-z0-9]$/)"`
	InputDockerStream string `json:"inputDockerStream"`
}

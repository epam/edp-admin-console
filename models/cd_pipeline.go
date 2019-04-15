package models

type CDPipelineDTO struct {
	Name             string              `json:"name"`
	CodebaseBranches []CodebaseBranchDTO `json:"codebaseBranches"`
}

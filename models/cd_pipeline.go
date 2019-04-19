package models

type CDPipelineDTO struct {
	Name             string              `json:"name"`
	Status           string              `json:"status"`
	CodebaseBranches []CodebaseBranchDTO `json:"codebaseBranches"`
}

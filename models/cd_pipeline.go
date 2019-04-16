package models

type CDPipelineReadRestApi struct {
	Name             string                      `json:"name"`
	CodebaseBranches []CodebaseBranchReadRestApi `json:"codebaseBranches"`
}

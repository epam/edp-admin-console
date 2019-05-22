package models

type CodebaseCriteria struct {
	Status *string
	Type   *string
}

type BranchCriteria struct {
	Status *string
}

type CDPipelineCriteria struct {
	Status *string
}

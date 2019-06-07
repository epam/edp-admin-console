package command

type CreateCodebaseBranch struct {
	Name     string `json:"name" valid:"Required;Match(/^[a-z0-9][a-z0-9-._]*[a-z0-9]$/)"`
	Commit   string `json:"commit"`
	Username string `json:"username"`
}

type BranchCriteria struct {
	Status *string
}

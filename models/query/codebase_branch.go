package query

type CodebaseBranch struct {
	Id         int       `json:"id" orm:"column(id)"`
	Name       string    `json:"branchName" orm:"column(name)"`
	FromCommit string    `json:"from_commit" orm:"column(from_commit)"`
	Status     string    `json:"status" orm:"column(status)"`
	VCSLink    string    `json:"branchLink" orm:"-"`
	CICDLink   string    `json:"jenkinsLink" orm:"-"`
	AppName    string    `json:"appName" orm:"-"`
	Codebase   *Codebase `json:"-" orm:"rel(fk)"`
}

func (cb *CodebaseBranch) TableName() string {
	return "codebase_branch"
}

type CodebaseBranchCriteria struct {
	Status Status
}

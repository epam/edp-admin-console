package query

type CodebaseBranch struct {
	Id         int       `json:"id" orm:"column(id)"`
	Name       string    `json:"name" orm:"column(name)"`
	FromCommit string    `json:"from_commit" orm:"column(from_commit)"`
	Status     string    `json:"status" orm:"column(status)"`
	VCSLink    string    `json:"vcs_link" orm:"-"`
	CICDLink   string    `json:"cicd_link" orm:"-"`
	Codebase   *Codebase `orm:"rel(fk)"`
}

func (cb *CodebaseBranch) TableName() string {
	return "codebase_branch"
}

type CodebaseBranchCriteria struct {
	Status Status
}

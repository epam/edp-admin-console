package query

type CDPipeline struct {
	Id                int                  `json:"-" orm:"column(id)"`
	Name              string               `json:"name" orm:"column(name)"`
	Status            string               `json:"status" orm:"column(status)"`
	JenkinsLink       string               `json:"jenkinsLink" orm:"-"`
	CodebaseBranch    []*CodebaseBranch    `json:"codebaseBranches" orm:"rel(m2m);rel_table(cd_pipeline_codebase_branch)"`
	Stage             []*Stage             `json:"stages" orm:"reverse(many)"`
	ThirdPartyService []*ThirdPartyService `json:"services" orm:"rel(m2m);rel_table(cd_pipeline_third_party_service)"`
}

type CDPipelineCriteria struct {
	Status Status
}

func (cb *CDPipeline) TableName() string {
	return "cd_pipeline"
}

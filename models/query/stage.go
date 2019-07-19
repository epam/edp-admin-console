package query

type Stage struct {
	Id                   int           `json:"id" orm:"column(id)"`
	Name                 string        `json:"name" orm:"column(name)"`
	Description          string        `json:"description" orm:"column(description)"`
	TriggerType          string        `json:"triggerType" orm:"column(trigger_type)"`
	Order                int           `json:"order" orm:"column(order)"`
	OpenshiftProjectLink string        `json:"openshiftProjectLink" orm:"-"`
	OpenshiftProjectName string        `json:"openshiftProjectName" orm:"-"`
	CDPipeline           *CDPipeline   `json:"-" orm:"rel(fk);column(cd_pipeline_id)"`
	QualityGates         []QualityGate `json:"qualityGates" orm:"-"`
}

type QualityGate struct {
	Id               int             `json:"id" orm:"column(id)"`
	QualityGateType  string          `json:"qualityGateType" orm:"column(quality_gate)"`
	StepName         string          `json:"stepName" orm:"column(step_name)"`
	CdStageId        *int            `json:"cdStageId" orm:"column(cd_stage_id)"`
	CodebaseId       *int            `json:"-" orm:"column(codebase_id)"`
	CodebaseBranchId *int            `json:"-" orm:"column(codebase_branch_id)"`
	Autotest         *Codebase       `json:"autotest" orm:"-"`
	Branch           *CodebaseBranch `json:"codebaseBranch" orm:"-"`
}

func (cb *Stage) TableName() string {
	return "cd_stage"
}

func (qg *QualityGate) TableName() string {
	return "quality_gate_stage"
}

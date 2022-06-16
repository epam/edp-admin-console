package command

import (
	cdPipelineAPI "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
)

type CDStageCommand struct {
	Name            string                      `json:"name" valid:"Required;Match(/^[a-z0-9]([-a-z0-9]*[a-z0-9])$/)"`
	Description     string                      `json:"description" valid:"Required"`
	TriggerType     string                      `json:"triggerType" valid:"Required"`
	Order           int                         `json:"order" valid:"Match(/^[0-9]$/)"`
	Source          cdPipelineAPI.Source        `json:"source"`
	QualityGates    []cdPipelineAPI.QualityGate `json:"qualityGates" valid:"Required"`
	Username        string                      `json:"username"`
	JobProvisioning string                      `json:"jobProvisioning"`
}

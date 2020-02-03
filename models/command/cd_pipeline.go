package command

import "edp-admin-console/models"

type CDPipelineCommand struct {
	Name                 string                                `json:"name" valid:"Required;Match(/^[a-z0-9]([-a-z0-9]*[a-z0-9])$/)"`
	Applications         []models.CDPipelineApplicationCommand `json:"applications" valid:"Required"`
	ThirdPartyServices   []string                              `json:"services"`
	Stages               []CDStageCommand                      `json:"stages"`
	ApplicationToApprove []string                              `json:"-"`
	Username             string                                `json:"username"`
}

type DeleteStageCommand struct {
	Name           string `json:"name"`
	CDPipelineName string `json:"pipelineName"`
}

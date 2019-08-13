package command

import "edp-admin-console/models"

type CDStageCommand struct {
	Name         string               `json:"name" valid:"Required;Match(/^[a-z0-9]([-a-z0-9]*[a-z0-9])$/)"`
	Description  string               `json:"description" valid:"Required"`
	TriggerType  string               `json:"triggerType" valid:"Required"`
	Order        int                  `json:"order" valid:"Match(/^[0-9]$/)"`
	QualityGates []models.QualityGate `json:"qualityGates" valid:"Required"`
	Username     string               `json:"username"`
}

package models

import "time"

type ReleaseBranch struct {
	Name            string    `json:"name"`
	Event           string    `json:"event"`
	DetailedMessage string    `json:"detailed_message"`
	Username        string    `json:"username"`
	UpdatedAt       time.Time `json:"updated_at"`
	VCSLink         string
	CICDLink        string
}

type ReleaseBranchRequestData struct {
	Name   string `json:"name" valid:"Required;Match(/^[a-z][a-z0-9-.]*[a-z0-9]$/)"`
	Commit string `json:"commit"`
}

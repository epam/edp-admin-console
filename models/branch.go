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
	Name   string `json:"name" valid:"Required;Match(/^[a-z\\d](?:[a-z\\d]|-([a-z\\d])){0,38}$/)"`
	Commit string `json:"commit"`
}

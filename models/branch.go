package models

import "time"

type ReleaseBranchView struct {
	Name            string    `json:"name"`
	Event           string    `json:"event"`
	DetailedMessage string    `json:"detailed_message"`
	Username        string    `json:"username"`
	UpdatedAt       time.Time `json:"updated_at"`
	VCSLink         string
	CICDLink        string
}

type ReleaseBranchCreateCommand struct {
	Name   string `json:"name" valid:"Required;Match(/^[a-z0-9][a-z0-9-._]*[a-z0-9]$/)"`
	Commit string `json:"commit"`
}

type ReleaseBranchCreatePipelineCommand struct {
	AppName    string
	BranchName string
}

type CodebaseBranchDTO struct {
	AppName     string `json:"appName"`
	BranchName  string `json:"branchName"`
	BranchLink  string `json:"branchLink"`
	JenkinsLink string `json:"jenkinsLink"`
}

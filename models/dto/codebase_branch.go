package dto

type CodebaseBranchDTO struct {
	AppName     string `json:"appName"`
	BranchName  string `json:"branchName"`
	BranchLink  string `json:"branchLink"`
	JenkinsLink string `json:"jenkinsLink"`
}

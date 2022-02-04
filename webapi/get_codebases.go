package webapi

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type GetCodebasesResponse []GetCodebase

type GetCodebase struct {
	ID                   int                    `json:"id"` // legacy ?
	Name                 string                 `json:"name"`
	Language             string                 `json:"language"`
	BuildTool            string                 `json:"build_tool"`
	Framework            string                 `json:"framework"`
	Strategy             string                 `json:"strategy"`
	GitURL               string                 `json:"git_url"`
	Type                 string                 `json:"type"`
	Status               string                 `json:"status"`
	TestReportFramework  string                 `json:"testReportFramework"`
	Description          string                 `json:"description"`
	CodebaseBranch       []CodebaseBranch       `json:"codebase_branch"`
	GitServer            string                 `json:"gitServer"`
	GitProjectPath       *string                `json:"gitProjectPath"`
	JenkinsSlave         string                 `json:"jenkinsSlave"`
	JobProvisioning      string                 `json:"jobProvisioning"`
	DeploymentScript     string                 `json:"deploymentScript"`
	VersioningType       string                 `json:"versioningType"`
	StartFrom            *string                `json:"startFrom"`
	JiraServer           *string                `json:"jiraServer"`
	CommitMessagePattern string                 `json:"commitMessagePattern"`
	TicketNamePattern    string                 `json:"ticketNamePattern"`
	CiTool               string                 `json:"ciTool"`
	Perf                 *Perf                  `json:"perf"`
	DefaultBranch        string                 `json:"defaultBranch"`
	JiraIssueFields      map[string]interface{} `json:"jiraIssueFields"`
	EmptyProject         bool                   `json:"emptyProject"`
}

type CodebaseBranch struct {
	ID                   int                    `json:"id"`
	BranchName           string                 `json:"branchName"`
	FromCommit           string                 `json:"from_commit"`
	Status               string                 `json:"status"`
	Version              *string                `json:"version"`
	BuildNumber          *string                `json:"build_number"`
	LastSuccessBuild     *string                `json:"last_success_build"`
	BranchLink           string                 `json:"branchLink"`
	JenkinsLink          string                 `json:"jenkinsLink"`
	AppName              string                 `json:"appName"`
	Release              bool                   `json:"release"`
	CodebaseDockerStream []CodebaseDockerStream `json:"codebaseDockerStream"`
}

type CodebaseDockerStream struct {
	ID                int    `json:"id"`
	OcImageStreamName string `json:"ocImageStreamName"`
	ImageLink         string `json:"imageLink"`
	JenkinsLink       string `json:"jenkinsLink"`
	// CdPipelines       interface{} `json:"CdPipelines"` // looks suspicious
}

type Perf struct {
	Name        string   `json:"name"`
	DataSources []string `json:"dataSources"`
}

func (h *HandlerEnv) GetCodebases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := LoggerFromContext(ctx)
	urlCodebases := r.URL.Query().Get("codebases")
	if urlCodebases == "" {
		logger.Error("codebases not passed")
		BadRequestResponse(ctx, w, "codebases not passed")
		return
	}

	codebases := strings.Split(urlCodebases, ",")
	cleanCodebases := make([]string, 0)
	for _, codebase := range codebases {
		if codebase != "" {
			cleanCodebases = append(cleanCodebases, codebase)
		}
	}
	if len(cleanCodebases) == 0 {
		logger.Error("empty codebases")
		BadRequestResponse(ctx, w, "empty codebases")
		return
	}

	codebasesResponse := make([]GetCodebase, 0)
	for _, codebaseName := range cleanCodebases {
		crCodebase, err := h.NamespacedClient.GetCodebase(ctx, codebaseName)
		if err != nil {
			logger.Error("get codebase by name failed", zap.Error(err), zap.String("codebase_name", codebaseName))
			InternalErrorResponse(ctx, w, "get codebase by name failed")
			return
		}
		codebaseResponse := GetCodebase{
			Name:      crCodebase.Name,
			GitServer: crCodebase.Spec.GitServer,
		}
		codebasesResponse = append(codebasesResponse, codebaseResponse)
	}

	response := GetCodebasesResponse(codebasesResponse)
	OKJsonResponse(ctx, w, response)
}

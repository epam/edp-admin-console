package webapi

import (
	"net/http"
	"strings"

	"go.uber.org/zap"

	"edp-admin-console/internal/applog"
)

type GetCodebasesResponse []CodebaseForList

type CodebaseMini struct {
	Name             string `json:"name"`
	DeploymentScript string `json:"deploymentScript"`
	GitServer        string `json:"gitServer"`
	VersioningType   string `json:"versioningType"`
	Strategy         string `json:"strategy"`
}

type CodebaseForList struct {
	CodebaseMini
	Type           string  `json:"type"`
	GitProjectPath *string `json:"gitProjectPath"`
	JenkinsSlave   string  `json:"jenkinsSlave"`
	EmptyProject   bool    `json:"emptyProject"`
}

func (h *HandlerEnv) GetCodebases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := applog.LoggerFromContext(ctx)
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

	codebasesResponse := make([]CodebaseForList, 0)
	for _, codebaseName := range cleanCodebases {
		crCodebase, err := h.NamespacedClient.GetCodebase(ctx, codebaseName)
		if err != nil {
			logger.Error("get codebase by name failed", zap.Error(err), zap.String("codebase_name", codebaseName))
			InternalErrorResponse(ctx, w, "get codebase by name failed")
			return
		}
		codebaseResponse := CodebaseForList{
			CodebaseMini: CodebaseMini{
				Name:             crCodebase.Name,
				GitServer:        crCodebase.Spec.GitServer,
				Strategy:         string(crCodebase.Spec.Strategy),
				DeploymentScript: crCodebase.Spec.DeploymentScript,
				VersioningType:   string(crCodebase.Spec.Versioning.Type),
			},
			GitProjectPath: crCodebase.Spec.GitUrlPath,
		}
		codebasesResponse = append(codebasesResponse, codebaseResponse)
	}

	response := GetCodebasesResponse(codebasesResponse)
	OKJsonResponse(ctx, w, response)
}

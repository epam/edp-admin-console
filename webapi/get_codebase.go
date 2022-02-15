package webapi

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type GetCodebaseResponse GetCodebase

type GetCodebase struct {
	CodebaseMini
	BuildTool       string `json:"build_tool"`
	Type            string `json:"type"`
	EmptyProject    bool   `json:"emptyProject"`
	JenkinsSlave    string `json:"jenkinsSlave"`
	Language        string `json:"language"`
	Framework       string `json:"framework"`
	JobProvisioning string `json:"jobProvisioning"`
}

func (h *HandlerEnv) GetCodebase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := LoggerFromContext(ctx)
	codebaseName := chi.URLParam(r, "codebaseName")
	if codebaseName == "" {
		logger.Error("codebase not passed")
		BadRequestResponse(ctx, w, "codebase not passed")
		return
	}

	crCodebase, err := h.NamespacedClient.GetCodebase(ctx, codebaseName)
	if err != nil {
		logger.Error("get codebase by name failed", zap.Error(err), zap.String("codebase_name", codebaseName))
		InternalErrorResponse(ctx, w, "get codebase by name failed")
		return
	}
	codebaseResponse := GetCodebase{
		CodebaseMini: CodebaseMini{
			Name:             crCodebase.Name,
			DeploymentScript: crCodebase.Spec.DeploymentScript,
			GitServer:        crCodebase.Spec.GitServer,
			Strategy:         string(crCodebase.Spec.Strategy),
			VersioningType:   string(crCodebase.Spec.Versioning.Type),
		},
		BuildTool:       strings.ToLower(crCodebase.Spec.BuildTool),
		Type:            crCodebase.Spec.Type,
		EmptyProject:    crCodebase.Spec.EmptyProject,
		JenkinsSlave:    pointerToStr(crCodebase.Spec.JenkinsSlave),
		Language:        crCodebase.Spec.Lang,
		Framework:       pointerToStr(crCodebase.Spec.Framework),
		JobProvisioning: pointerToStr(crCodebase.Spec.JobProvisioning),
	}

	response := GetCodebaseResponse(codebaseResponse)
	OKJsonResponse(ctx, w, response)
}

func pointerToStr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

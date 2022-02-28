package webapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
)

type GetCodebaseResponse GetCodebase

type GetCodebase struct {
	CodebaseMini
	BuildTool            string           `json:"build_tool"`
	Type                 string           `json:"type"`
	EmptyProject         bool             `json:"emptyProject"`
	JenkinsSlave         string           `json:"jenkinsSlave"`
	Language             string           `json:"language"`
	Framework            string           `json:"framework"`
	JobProvisioning      string           `json:"jobProvisioning"`
	CodebaseBranch       []CodebaseBranch `json:"codebase_branch"`
	CommitMessagePattern string           `json:"commitMessagePattern"`
	TicketNamePattern    string           `json:"ticketNamePattern"`
	DefaultBranch        string           `json:"defaultBranch"`
	TestReportFramework  string           `json:"testReportFramework"`
}

type CodebaseBranch struct {
	BranchName  string  `json:"branchName"`
	BuildNumber string  `json:"build_number"` //number, but string (sick)
	Version     *string `json:"version"`
	Release     bool    `json:"release"`
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
		var statusErr *k8sErrors.StatusError
		if errors.As(err, &statusErr) {
			if statusErr.ErrStatus.Code == http.StatusNotFound {
				logger.Error("codebase not found", zap.String("codebase_name", codebaseName))
				NotFoundResponse(ctx, w, "codebase not found")
				return
			}
		}
		logger.Error("get codebase by name failed", zap.Error(err), zap.String("codebase_name", codebaseName))
		InternalErrorResponse(ctx, w, "get codebase by name failed")
		return
	}
	crCbBranchList, err := h.NamespacedClient.CodebaseBranchesListByCodebaseName(ctx, crCodebase.Name)
	if err != nil {
		logger.Error("get codebase_branch list by name failed", zap.Error(err), zap.String("codebase_name", codebaseName))
		InternalErrorResponse(ctx, w, "get codebase branch list by name failed")
		return
	}
	codebaseBranches := make([]CodebaseBranch, len(crCbBranchList))
	for i := range crCbBranchList {
		crCbBranch := crCbBranchList[i]
		codebaseBranches[i] = CodebaseBranch{
			BranchName:  crCbBranch.Spec.BranchName,
			BuildNumber: strPointerValueOrDefault(crCbBranch.Status.Build, ""),
			Version:     crCbBranch.Spec.Version,
			Release:     crCbBranch.Spec.Release,
		}
	}

	codebaseResponse := GetCodebase{
		CodebaseMini: CodebaseMini{
			Name:             crCodebase.Name,
			DeploymentScript: crCodebase.Spec.DeploymentScript,
			GitServer:        crCodebase.Spec.GitServer,
			Strategy:         string(crCodebase.Spec.Strategy),
			VersioningType:   string(crCodebase.Spec.Versioning.Type),
		},
		BuildTool:            strings.ToLower(crCodebase.Spec.BuildTool),
		Type:                 crCodebase.Spec.Type,
		EmptyProject:         crCodebase.Spec.EmptyProject,
		JenkinsSlave:         pointerToStr(crCodebase.Spec.JenkinsSlave),
		Language:             crCodebase.Spec.Lang,
		Framework:            pointerToStr(crCodebase.Spec.Framework),
		JobProvisioning:      pointerToStr(crCodebase.Spec.JobProvisioning),
		CodebaseBranch:       codebaseBranches,
		CommitMessagePattern: pointerToStr(crCodebase.Spec.CommitMessagePattern),
		TicketNamePattern:    pointerToStr(crCodebase.Spec.TicketNamePattern),
		DefaultBranch:        crCodebase.Spec.DefaultBranch,
		TestReportFramework:  pointerToStr(crCodebase.Spec.TestReportFramework),
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

func strPointerValueOrDefault(str *string, defaultValue string) string {
	if str == nil {
		return defaultValue
	}
	return *str
}

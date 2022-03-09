package webapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"edp-admin-console/internal/applog"
)

type cdPipelineResponse struct {
	CodebaseBranch        []*codebaseBranchResponse `json:"codebaseBranches"`
	ApplicationsToPromote []string                  `json:"applicationsToPromote"`
}

type codebaseBranchResponse struct {
	AppName string `json:"appName"`
}

func (h *HandlerEnv) GetPipeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := applog.LoggerFromContext(ctx)
	cdPipelineName := chi.URLParam(r, "pipelineName")
	if cdPipelineName == "" {
		logger.Error("pipeline not passed")
		NotFoundResponse(ctx, w, "pipeline not passed")
		return
	}

	cdPipeline, err := h.NamespacedClient.GetCDPipeline(ctx, cdPipelineName)
	if err != nil {
		logger.Error("get CDPipeline by name failed", zap.Error(err), zap.String("cdPipelineName", cdPipelineName))
		InternalErrorResponse(ctx, w, "get CDPipeline by name failed")
		return
	}

	applications := cdPipeline.Spec.ApplicationsToPromote

	if len(applications) == 0 {
		logger.Error("empty applications to promote")
		InternalErrorResponse(ctx, w, "empty applications to promote")
		return
	}

	codebaseBranches := make([]*codebaseBranchResponse, 0)
	for _, application := range applications {
		if application != "" {
			codebaseBranches = append(codebaseBranches, &codebaseBranchResponse{
				AppName: application,
			})
		}
	}

	response := &cdPipelineResponse{
		CodebaseBranch:        codebaseBranches,
		ApplicationsToPromote: applications,
	}

	OKJsonResponse(ctx, w, response)
}

package webapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type StagePipelineResponse struct {
	StageName      string `json:"name"`
	CDPipelineName string `json:"cdPipeline"`
}

func (h *HandlerEnv) GetStagePipeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := LoggerFromContext(ctx)
	logger.Debug("in handler")
	stageName := chi.URLParam(r, "stageName")
	cdPipelineName := chi.URLParam(r, "pipelineName")
	response := &StagePipelineResponse{
		StageName:      stageName,
		CDPipelineName: cdPipelineName,
	}
	OKJsonResponse(ctx, w, response)
}

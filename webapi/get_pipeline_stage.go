package webapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"edp-admin-console/internal/applog"
	"edp-admin-console/internal/imagestream"
	pipelinestage "edp-admin-console/internal/pipeline-stage"
)

type Response struct {
	Name            string                           `json:"name"`
	CDPipeline      string                           `json:"cdPipeline"`
	Description     string                           `json:"description"`
	TriggerType     string                           `json:"triggerType"`
	Order           string                           `json:"order"`
	Applications    []pipelinestage.ApplicationStage `json:"applications"`
	QualityGates    []pipelinestage.QualityGate      `json:"qualityGates"`
	JobProvisioning string                           `json:"jobProvisioning"`
}

func (h *HandlerEnv) GetStagePipeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := applog.LoggerFromContext(ctx)
	logger.Debug("in handler")

	stageName := chi.URLParam(r, "stageName")
	cdPipelineName := chi.URLParam(r, "pipelineName")
	stageCRName := cdPipelineName + "-" + stageName

	partialStageView, err := pipelinestage.StageViewByCRName(ctx, h.NamespacedClient, stageCRName)
	if err != nil {
		logger.Error(err.Error())
		NotFoundResponse(ctx, w, "stage not found")
		return
	}

	cdPipelineAppNames, err := pipelinestage.CdPipelineAppNamesByCRName(ctx, h.NamespacedClient, cdPipelineName)
	if err != nil {
		logger.Error(err.Error())
		NotFoundResponse(ctx, w, "cdPipeline not found")
		return
	}

	inputIS, err := imagestream.GetInputISForStage(ctx, h.NamespacedClient, stageName, cdPipelineName)
	if err != nil {
		logger.Error("get input image stream failed", zap.Error(err),
			zap.String("stage_name", stageName), zap.String("cd_pipeline_name", cdPipelineName))
		NotFoundResponse(ctx, w, "input IS not found")
		return
	}

	outputIS, err := imagestream.GetOutputISForStage(ctx, h.NamespacedClient, cdPipelineName, stageName)
	if err != nil {
		logger.Error(err.Error())
		NotFoundResponse(ctx, w, "output IS not found")
		return
	}

	applications, err := pipelinestage.BuildApplicationStages(ctx, h.NamespacedClient, inputIS, outputIS, cdPipelineAppNames)
	if err != nil {
		logger.Error(err.Error())
		NotFoundResponse(ctx, w, "inputIS and outputIS not the same size")
		return
	}

	response := &Response{
		Name:            stageName,
		CDPipeline:      cdPipelineName,
		TriggerType:     partialStageView.TriggerType,
		Order:           partialStageView.Order,
		JobProvisioning: partialStageView.JobProvisioning,
		Description:     partialStageView.Description,
		Applications:    applications,
		QualityGates:    partialStageView.QualityGates,
	}
	OKJsonResponse(ctx, w, response)
}

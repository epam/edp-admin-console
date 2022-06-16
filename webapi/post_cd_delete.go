package webapi

import (
	"fmt"
	"net/http"
	"sort"

	"edp-admin-console/internal/cdpipelines"

	cdPipelineAPI "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	"go.uber.org/zap"

	"edp-admin-console/internal/applog"
	pipelinestage "edp-admin-console/internal/pipeline-stage"
)

func (h *HandlerEnv) DeleteCD(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := applog.LoggerFromContext(ctx)

	err := r.ParseForm()
	if err != nil {
		logger.Error("parse form failed", zap.Error(err))
		InternalErrorResponse(ctx, w, "parse form failed")
		return
	}

	cdName := r.FormValue("name")
	cdPipelineCR, err := cdpipelines.ByNameIFExists(ctx, h.NamespacedClient, cdName)
	if err != nil {
		logger.Error("get cd pipeline failed", zap.Error(err), zap.String("cd_pipeline_name", cdName))
		http.Redirect(w, r, fmt.Sprintf("%s/v2/admin/edp/cd-pipeline/overview?name=%s#cdPipelineDeletedErrorModal", h.Config.BasePath, cdName), http.StatusFound)
		return
	}

	if cdPipelineCR == nil {
		logger.Error("cd pipeline not found", zap.String("cd_pipeline_name", cdName))
		http.Redirect(w, r, fmt.Sprintf("%s/v2/admin/edp/cd-pipeline/overview#cdPipelineDeletedErrorModal", h.Config.BasePath), http.StatusFound)
		return
	}

	stageCRList, err := pipelinestage.StageListByPipelineName(ctx, h.NamespacedClient, cdName)
	if err != nil {
		logger.Error("get cd pipeline stage list failed", zap.Error(err), zap.String("cd_pipeline_name", cdName))
		http.Redirect(w, r, fmt.Sprintf("%s/v2/admin/edp/cd-pipeline/overview?name=%s#cdPipelineDeletedErrorModal", h.Config.BasePath, cdName), http.StatusFound)
		return
	}

	orderedStageCRList := stageListFromLastToFirst(stageCRList)
	for _, stageCR := range orderedStageCRList {
		deleteErr := h.NamespacedClient.DeleteStage(ctx, &stageCR)
		if deleteErr != nil {
			logger.Error("delete stage CR failed", zap.Error(deleteErr), zap.String("stage_cr_name", stageCR.Name))
			http.Redirect(w, r, fmt.Sprintf("%s/v2/admin/edp/cd-pipeline/overview?name=%s#cdPipelineDeletedErrorModal", h.Config.BasePath, cdName), http.StatusFound)
		}
	}

	err = h.NamespacedClient.DeleteCDPipeline(ctx, cdPipelineCR)
	if err != nil {
		logger.Error("delete cd pipeline failed", zap.Error(err), zap.String("cd_pipeline_cr_name", cdPipelineCR.Name))
		http.Redirect(w, r, fmt.Sprintf("%s/v2/admin/edp/cd-pipeline/overview?name=%s#cdPipelineDeletedErrorModal", h.Config.BasePath, cdName), http.StatusFound)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/v2/admin/edp/cd-pipeline/overview?name=%s#cdPipelineDeletedSuccessModal", h.Config.BasePath, cdName), http.StatusFound)
}

func stageListFromLastToFirst(stageList []cdPipelineAPI.Stage) []cdPipelineAPI.Stage {
	sort.Slice(stageList, func(i, j int) bool {
		return stageList[i].Spec.Order > stageList[j].Spec.Order
	})
	return stageList
}

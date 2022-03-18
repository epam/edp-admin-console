package webapi

import (
	"fmt"
	"net/http"
	"strings"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"edp-admin-console/internal/applog"
)

func (h *HandlerEnv) UpdateCDPipeline(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := applog.LoggerFromContext(ctx)
	logger.Debug("in handler")

	err := request.ParseForm()
	if err != nil {
		logger.Error("cant parse form", zap.Error(err))
		InternalErrorResponse(ctx, writer, "cant parse form")
		return
	}

	cdPipeName := chi.URLParam(request, "pipelineName")

	apps := request.Form["app"]
	var appsInfo []appInfoPipelineCreation
	for i := range apps {
		promote := false
		if strings.ToLower(request.Form.Get(strings.Join([]string{apps[i], "promote"}, "-"))) == "true" {
			promote = true
		}
		appInfo := appInfoPipelineCreation{
			Name:      apps[i],
			InputIS:   request.Form.Get(apps[i]),
			IsPromote: promote,
		}
		appsInfo = append(appsInfo, appInfo)
	}

	stages := request.Form["stageName"]
	var stagesInfo []stageInfoPipelineCreation
	for i := range stages {
		steps := request.Form[stages[i]+"-stageStepName"]
		var stepsInfo []stepInfoPipelineCreation
		for _, step := range steps {
			stepInfo := stepInfoPipelineCreation{
				Name:            step,
				Autotest:        strToPtr(request.Form.Get(strings.Join([]string{stages[i], step, "stageAutotests"}, "-"))),
				StageBranch:     strToPtr(request.Form.Get(strings.Join([]string{stages[i], step, "stageBranch"}, "-"))),
				QualityGateType: request.Form.Get(strings.Join([]string{stages[i], step, "stageQualityGateType"}, "-")),
			}
			stepsInfo = append(stepsInfo, stepInfo)
		}

		stageInfo := stageInfoPipelineCreation{
			Name:                  stages[i],
			Steps:                 stepsInfo,
			Description:           request.Form.Get(strings.Join([]string{stages[i], "stageDesc"}, "-")),
			TriggerType:           request.Form.Get(strings.Join([]string{stages[i], "triggerType"}, "-")),
			PipelineLibraryName:   request.Form.Get(strings.Join([]string{stages[i], "pipelineLibraryName"}, "-")),
			JobProvisioner:        request.Form.Get(strings.Join([]string{stages[i], "jobProvisioning"}, "-")),
			PipelineLibraryBranch: request.Form.Get(strings.Join([]string{stages[i], "pipelineLibraryBranch"}, "-")),
		}
		stagesInfo = append(stagesInfo, stageInfo)
	}

	var appsToPromote []string
	var inputDockerStreams []string
	for i := range appsInfo {
		if appsInfo[i].IsPromote {
			appsToPromote = append(appsToPromote, appsInfo[i].Name)
		}
		inputDockerStreams = append(inputDockerStreams, appsInfo[i].InputIS)
	}

	existedCDPipeCR, err := h.NamespacedClient.GetCDPipeline(ctx, cdPipeName)
	if err != nil {
		http.Redirect(writer, request, fmt.Sprintf("%s/admin/edp/cd-pipeline/overview?name=%s#cdPipelineEditErrorModal'", h.Config.BasePath, cdPipeName), http.StatusFound)
		logger.Error("cant get cd pipeline", zap.Error(err))

		return
	}
	existedCDPipeCR.Spec.InputDockerStreams = inputDockerStreams
	existedCDPipeCR.Spec.ApplicationsToPromote = appsToPromote
	err = h.NamespacedClient.Update(ctx, existedCDPipeCR)
	if err != nil {
		http.Redirect(writer, request, fmt.Sprintf("%s/admin/edp/cd-pipeline/overview?name=%s#cdPipelineEditErrorModal'", h.Config.BasePath, cdPipeName), http.StatusFound)
		logger.Error("cant get cd pipeline", zap.Error(err))
		return
	}

	for i, stage := range stagesInfo {
		var qualityGates []cdPipeApi.QualityGate
		for _, step := range stage.Steps {
			qualityGate := cdPipeApi.QualityGate{
				QualityGateType: step.QualityGateType,
				StepName:        step.Name,
				AutotestName:    step.Autotest,
				BranchName:      step.StageBranch,
			}
			qualityGates = append(qualityGates, qualityGate)
		}
		var stageLibrary cdPipeApi.Library
		stageLibraryType := "default"
		if stage.PipelineLibraryName != "default" {
			stageLibraryType = "library"
			stageLibrary = cdPipeApi.Library{
				Name:   stage.PipelineLibraryName,
				Branch: stage.PipelineLibraryBranch,
			}
		}

		stageSpec := cdPipeApi.StageSpec{
			Name:         stage.Name,
			CdPipeline:   cdPipeName,
			Description:  stage.Description,
			TriggerType:  stage.TriggerType,
			Order:        i,
			QualityGates: qualityGates,
			Source: cdPipeApi.Source{
				Type:    stageLibraryType,
				Library: stageLibrary,
			},
			JobProvisioning: stage.JobProvisioner,
		}
		stageName := strings.Join([]string{cdPipeName, stage.Name}, "-")
		errStage := h.NamespacedClient.CreateCDStageBySpec(ctx, stageName, stageSpec)
		if errStage != nil {
			logger.Error("cant create stage CR", zap.Error(errStage))
			http.Redirect(writer, request, fmt.Sprintf("%s/admin/edp/cd-pipeline/overview?name=%s#stageCreateErrorModal'", h.Config.BasePath, stageName), http.StatusFound)
			return
		}
	}

	http.Redirect(writer, request, fmt.Sprintf("%s/admin/edp/cd-pipeline/overview#cdPipelineEditSuccessModal", h.Config.BasePath), http.StatusFound)
}

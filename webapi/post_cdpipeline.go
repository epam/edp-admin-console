package webapi

import (
	"fmt"
	"net/http"
	"strings"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"edp-admin-console/internal/applog"
)

const paramWaitingForCdPipeline = "waitingforcdpipeline"

type appInfoPipelineCreation struct {
	Name      string
	InputIS   string
	IsPromote bool
}

type stepInfoPipelineCreation struct {
	Name            string
	Autotest        *string
	StageBranch     *string
	QualityGateType string
}

type stageInfoPipelineCreation struct {
	Name                  string
	Steps                 []stepInfoPipelineCreation
	Description           string
	TriggerType           string
	PipelineLibraryName   string
	JobProvisioner        string
	PipelineLibraryBranch string
}

func (h *HandlerEnv) CreateCDPipeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := applog.LoggerFromContext(ctx)
	logger.Debug("in handler")

	err := r.ParseForm()
	if err != nil {
		logger.Error("cant parse form", zap.Error(err))
		InternalErrorResponse(ctx, w, "cant parse form")
		return
	}

	cdPipeName := r.Form.Get("pipelineName")
	deploymentType := r.Form.Get("deploymentType")

	apps := r.Form["app"]
	var appsInfo []appInfoPipelineCreation
	for i := range apps {
		promote := false
		if strings.ToLower(r.Form.Get(strings.Join([]string{apps[i], "promote"}, "-"))) == "true" {
			promote = true
		}
		appInfo := appInfoPipelineCreation{
			Name:      apps[i],
			InputIS:   r.Form.Get(apps[i]),
			IsPromote: promote,
		}
		appsInfo = append(appsInfo, appInfo)
	}

	stages := r.Form["stageName"]
	var stagesInfo []stageInfoPipelineCreation
	for i := range stages {
		steps := r.Form[stages[i]+"-stageStepName"]
		var stepsInfo []stepInfoPipelineCreation
		for _, step := range steps {
			stepInfo := stepInfoPipelineCreation{
				Name:            step,
				Autotest:        strToPtr(r.Form.Get(strings.Join([]string{stages[i], step, "stageAutotests"}, "-"))),
				StageBranch:     strToPtr(r.Form.Get(strings.Join([]string{stages[i], step, "stageBranch"}, "-"))),
				QualityGateType: r.Form.Get(strings.Join([]string{stages[i], step, "stageQualityGateType"}, "-")),
			}
			stepsInfo = append(stepsInfo, stepInfo)
		}

		stageInfo := stageInfoPipelineCreation{
			Name:                  stages[i],
			Steps:                 stepsInfo,
			Description:           r.Form.Get(strings.Join([]string{stages[i], "stageDesc"}, "-")),
			TriggerType:           r.Form.Get(strings.Join([]string{stages[i], "triggerType"}, "-")),
			PipelineLibraryName:   r.Form.Get(strings.Join([]string{stages[i], "pipelineLibraryName"}, "-")),
			JobProvisioner:        r.Form.Get(strings.Join([]string{stages[i], "jobProvisioning"}, "-")),
			PipelineLibraryBranch: r.Form.Get(strings.Join([]string{stages[i], "pipelineLibraryBranch"}, "-")),
		}
		stagesInfo = append(stagesInfo, stageInfo)
	}

	var appsToPromote []string
	pipelineApps := make([]string, len(appsInfo), len(appsInfo))
	inputDockerStreams := make([]string, len(appsInfo), len(appsInfo))
	for i := range appsInfo {
		if appsInfo[i].IsPromote {
			appsToPromote = append(appsToPromote, appsInfo[i].Name)
		}
		pipelineApps[i] = appsInfo[i].Name
		inputDockerStreams[i] = appsInfo[i].InputIS
	}

	cdPipelineSpec := cdPipeApi.CDPipelineSpec{
		Name:                  cdPipeName,
		DeploymentType:        deploymentType,
		InputDockerStreams:    inputDockerStreams,
		Applications:          pipelineApps,
		ApplicationsToPromote: appsToPromote,
	}
	err = h.NamespacedClient.CreateCDPipelineBySpec(ctx, cdPipeName, cdPipelineSpec)
	if err != nil {
		logger.Error("cant create cdpipe CR", zap.Error(err))
		http.Redirect(w, r, fmt.Sprintf("%s/v2/admin/edp/cd-pipeline/overview?name=%s#cdPipelineCreateErrorModal'", h.Config.BasePath, cdPipeName), http.StatusFound)
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

		stageName := strings.Join([]string{cdPipeName, stage.Name}, "-")

		stageCR := &cdPipeApi.Stage{
			ObjectMeta: metav1.ObjectMeta{
				Name:      stageName,
				Namespace: h.NamespacedClient.Namespace,
			},
			Spec: cdPipeApi.StageSpec{
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
			},
		}

		errStage := h.NamespacedClient.Create(ctx, stageCR)
		if errStage != nil {
			logger.Error("cant create stage CR", zap.Error(errStage))
			http.Redirect(w, r, fmt.Sprintf("%s/v2/admin/edp/cd-pipeline/overview?name=%s#stageCreateErrorModal'", h.Config.BasePath, stageName), http.StatusFound)
			return
		}
	}
	http.Redirect(w, r, fmt.Sprintf("%s/v2/admin/edp/cd-pipeline/overview?%s=%s#cdPipelineSuccessModal", h.Config.BasePath, paramWaitingForCdPipeline, cdPipeName), http.StatusFound)
}

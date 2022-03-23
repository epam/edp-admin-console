package webapi

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"go.uber.org/zap"

	"edp-admin-console/internal"
	"edp-admin-console/internal/applog"
	"edp-admin-console/service/platform"
	"edp-admin-console/util/consts"
)

type cdPipelineOverviewTmpl struct {
	Username           string
	Xsrfdata           string
	Type               string
	DiagramPageEnabled bool
	BasePath           string
	EDPVersion         string
	IsAdmin            bool
	IsOpenshift        bool
	CDPipeline         cdPipeline
	Error              string
}

type library struct {
	Name   string
	Branch string
}

type source struct {
	Type    string
	Library library
}

type codebaseDockerStream struct {
	OcImageStreamName string
	CICDLink          string
	ImageLink         string
	CodebaseName      string
}

type autotest struct {
	Name string
}

type branch struct {
	Name string
}

type qualityGate struct {
	QualityGateType string
	Autotest        autotest
	Branch          branch
	StepName        string
}

type jobProvisioning struct {
	Name string
}

type cdStage struct {
	Name                string
	Order               int
	Description         string
	TriggerType         string
	PlatformProjectLink string
	QualityGates        []qualityGate
	Source              source
	JobProvisioning     jobProvisioning
}

type actionLog struct {
	LastTimeUpdate time.Time
	UserName       string
	Message        string
	Action         string
	Result         string
}

type cdPipeline struct {
	Name                  string
	JenkinsLink           string
	DeploymentType        string
	CodebaseDockerStream  []codebaseDockerStream
	Stage                 []cdStage
	CodebaseStageMatrix   bool        // IDK. it's needed for original html template. should be false
	ActionLog             []actionLog //todo: add action log
	ApplicationsToPromote []string
}

func (h *HandlerEnv) GetPipelineOverviewPage(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := applog.LoggerFromContext(ctx)
	user := UserFromContext(ctx)
	xsrf := csrf.Token(request)
	pipelineName := chi.URLParam(request, "pipelineName")
	cdPipelineCR, err := h.NamespacedClient.GetCDPipeline(ctx, pipelineName)
	if err != nil {
		logger.Error("cant get cdPipeline CR", zap.Error(err))
		InternalErrorResponse(ctx, writer, "cant get cdPipeline CR")
		return
	}
	jenkinsName := "jenkins"
	jenkinsEdp, err := h.NamespacedClient.EDPComponentByCRName(ctx, jenkinsName)
	if err != nil {
		logger.Error("cant get jenkins edp component", zap.Error(err))
		InternalErrorResponse(ctx, writer, "cant get jenkins edp component")
		return
	}
	jenkinsUrl := jenkinsEdp.Spec.Url

	var codebaseDockerStreams []codebaseDockerStream
	for _, cis := range cdPipelineCR.Spec.InputDockerStreams {
		cisCR, errCIS := h.NamespacedClient.GetCodebaseImageStream(ctx, cis)
		if errCIS != nil {
			logger.Error("cant get cdPipeline image streams", zap.Error(errCIS))
			InternalErrorResponse(ctx, writer, "cant get cdPipeline image streams")
			return
		}

		codebaseDockerStream := codebaseDockerStream{
			CodebaseName:      cisCR.Spec.Codebase,
			OcImageStreamName: cisCR.Name,
			CICDLink:          fmt.Sprintf("%s/job/%s", jenkinsUrl, cisCR.Spec.Codebase),
			ImageLink:         internal.AddSchemeIfNeeded(cisCR.Spec.ImageName),
		}
		codebaseDockerStreams = append(codebaseDockerStreams, codebaseDockerStream)
	}

	var stages []cdStage

	stageCRs, err := h.NamespacedClient.StageList(ctx)
	if err != nil {
		logger.Error("cant get cdPipeline image streams", zap.Error(err))
		InternalErrorResponse(ctx, writer, "cant get stage list")
		return
	}
	for _, stageCR := range stageCRs {
		if stageCR.Spec.CdPipeline == pipelineName {
			var qualityGates []qualityGate
			for _, qg := range stageCR.Spec.QualityGates {
				qualityGate := qualityGate{
					QualityGateType: qg.QualityGateType,
					Autotest:        autotest{Name: pointerToStr(qg.AutotestName)},
					Branch:          branch{Name: pointerToStr(qg.BranchName)},
					StepName:        qg.StepName,
				}
				qualityGates = append(qualityGates, qualityGate)
			}

			stage := cdStage{
				Name:                stageCR.Name,
				Order:               stageCR.Spec.Order,
				Description:         stageCR.Spec.Description,
				TriggerType:         stageCR.Spec.TriggerType,
				PlatformProjectLink: fmt.Sprintf("%s/console/project/%s-%s/overview", h.ClusterConfig.Host, h.NamespacedClient.Namespace, stageCR.Name),
				QualityGates:        qualityGates,
				Source: source{
					Type:    stageCR.Spec.Source.Type,
					Library: library(stageCR.Spec.Source.Library),
				},
				JobProvisioning: jobProvisioning{
					Name: stageCR.Spec.JobProvisioning,
				},
			}
			stages = append(stages, stage)
		}
	}

	var tmplData = cdPipelineOverviewTmpl{
		Username:           user.UserName(),
		Xsrfdata:           xsrf,
		Type:               consts.CDMenuItem,
		DiagramPageEnabled: h.Config.DiagramPageEnabled,
		BasePath:           h.Config.BasePath,
		EDPVersion:         h.Config.EDPVersion,
		IsAdmin:            user.IsAdmin(),
		IsOpenshift:        platform.IsOpenshift(),
		CDPipeline: cdPipeline{
			Name:                  cdPipelineCR.Name,
			DeploymentType:        cdPipelineCR.Spec.DeploymentType,
			ApplicationsToPromote: cdPipelineCR.Spec.ApplicationsToPromote,
			CodebaseDockerStream:  codebaseDockerStreams,
			Stage:                 stages,
		},
		Error: "",
	}

	templatePaths := []string{
		path.Join(h.WorkingDir, "/viewsV2/cd_pipeline_overview.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/footer_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/header_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/navbar_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/delete_confirmation_template.html"),
	}

	template := &Template{
		Data:             tmplData,
		FuncMap:          h.FuncMap,
		MainTemplateName: "cd_pipeline_overview.html",
		TemplatePaths:    templatePaths,
	}

	OkHTMLResponse(ctx, writer, template)
}

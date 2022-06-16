package webapi

import (
	"net/http"
	"path"

	codebaseAPI "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"go.uber.org/zap"

	"edp-admin-console/internal/applog"
	"edp-admin-console/models/query"
	"edp-admin-console/util/consts"
)

type CdPipelineForUpdate struct {
	Name                  string
	CodebaseBranch        []CodebaseBranchForUpdate
	ApplicationsToPromote []string
}

type CodebaseBranchForUpdate struct {
	AppName              string
	CodebaseDockerStream []CodebaseDockerStreamForUpdate
}

type CodebaseDockerStreamForUpdate struct {
	OcImageStreamName string
}

type pipelineUpdateData struct {
	CDPipeline         CdPipelineForUpdate
	Type               query.CodebaseType
	BasePath           string
	EDPVersion         string
	Username           string
	Xsrfdata           string
	Error              string
	DiagramPageEnabled bool
	Success            bool
	Name               string
	Apps               []Application
	Autotests          []Autotest
	GroovyLibs         []GroovyLib
	JobProvisioners    []JobProvisioner
}

func (h *HandlerEnv) GetPipelineUpdatePage(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := applog.LoggerFromContext(ctx)
	user := UserFromContext(ctx)
	k8sClient := h.NamespacedClient

	cbList, err := k8sClient.GetCodebaseList(ctx)
	if err != nil {
		logger.Error("get codebases list failed", zap.Error(err))
		InternalErrorResponse(ctx, writer, "get codebases list failed")
		return
	}

	apps, err := fetchApplications(ctx, k8sClient, cbList.Items)
	if err != nil {
		InternalErrorResponse(ctx, writer, err.Error())
		return
	}

	autotests, err := fetchAutotests(ctx, k8sClient, cbList.Items)
	if err != nil {
		InternalErrorResponse(ctx, writer, err.Error())
		return
	}

	libs, err := fetchLibs(ctx, k8sClient, cbList.Items)
	if err != nil {
		InternalErrorResponse(ctx, writer, err.Error())
		return
	}

	jobProvisioners, err := fetchJobProvisioners(ctx, k8sClient)
	if err != nil {
		InternalErrorResponse(ctx, writer, err.Error())
		return
	}

	csrfToken := csrf.Token(request)

	cdPipelineName := chi.URLParam(request, "pipelineName")
	cdPipelineCR, err := h.NamespacedClient.GetCDPipeline(ctx, cdPipelineName)
	if err != nil {
		InternalErrorResponse(ctx, writer, err.Error())
		return
	}

	codebaseBranch := make([]CodebaseBranchForUpdate, 0)

	for _, imageStream := range cdPipelineCR.Spec.InputDockerStreams {
		codebaseImageStream := &codebaseAPI.CodebaseImageStream{}
		codebaseImageStream, err = h.NamespacedClient.GetCodebaseImageStream(ctx, imageStream)
		codebaseBranch = append(codebaseBranch, CodebaseBranchForUpdate{
			AppName:              codebaseImageStream.Spec.Codebase,
			CodebaseDockerStream: []CodebaseDockerStreamForUpdate{{OcImageStreamName: imageStream}}, // we use an array here due to current html template implementation
			// CodebaseDockerStreamForUpdate should contain only one image stream for a chosen codebase.
		})
	}

	cdPipelineForUpdate := CdPipelineForUpdate{
		Name:                  cdPipelineCR.Name,
		CodebaseBranch:        codebaseBranch,
		ApplicationsToPromote: cdPipelineCR.Spec.ApplicationsToPromote,
	}

	tmplData := pipelineUpdateData{
		CDPipeline:         cdPipelineForUpdate,
		Type:               consts.CDMenuItem,
		BasePath:           h.Config.BasePath,
		EDPVersion:         h.Config.EDPVersion,
		Username:           user.UserName(),
		Xsrfdata:           csrfToken,
		DiagramPageEnabled: h.Config.DiagramPageEnabled,
		Success:            true,
		Name:               cdPipelineName,
		Apps:               apps,
		Autotests:          autotests,
		GroovyLibs:         libs,
		JobProvisioners:    jobProvisioners,
	}

	templatePaths := []string{path.Join(h.WorkingDir, "/viewsV2/edit_cd_pipeline.html"), path.Join(h.WorkingDir, "/views/template/footer_template.html"),
		path.Join(h.WorkingDir, "/views/template/header_template.html"), path.Join(h.WorkingDir, "/views/template/navbar_template.html"),
		path.Join(h.WorkingDir, "/views/template/modal_success_template.html")}

	template := Template{
		Data:             tmplData,
		FuncMap:          h.FuncMap,
		MainTemplateName: "edit_cd_pipeline.html",
		TemplatePaths:    templatePaths,
	}

	OkHTMLResponse(ctx, writer, &template)
}

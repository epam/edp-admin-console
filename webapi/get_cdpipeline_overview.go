package webapi

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/csrf"
	"go.uber.org/zap"

	"edp-admin-console/internal/applog"
	"edp-admin-console/util"
	"edp-admin-console/util/consts"
)

type cdPipelineForOverview struct {
	Name        string
	Status      string
	JenkinsLink string
}

type cdPipelinePageOverviewData struct {
	Username                      string
	IsAdmin                       bool
	ActiveApplicationsAndBranches bool
	Applications                  bool
	CDPipelines                   []cdPipelineForOverview
	Xsrfdata                      string
	Error                         string
	BasePath                      string
	Type                          string
	DiagramPageEnabled            bool
	EDPVersion                    string
	Success                       bool
}

func (h *HandlerEnv) CDPipelineOverview(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := applog.LoggerFromContext(ctx)
	user := UserFromContext(ctx)

	cdPipelineList, err := h.NamespacedClient.GetCDPipelineList(ctx)
	if err != nil {
		logger.Error("cant get cdpipeline list", zap.Error(err))
		InternalErrorResponse(ctx, writer, "cant get cdpipeline list")
		return
	}

	codebaseList, err := h.NamespacedClient.GetCodebaseList(ctx)
	if err != nil {
		logger.Error("cant get codebase list", zap.Error(err))
		InternalErrorResponse(ctx, writer, "cant get codebase list")
		return
	}

	appExist := false
	branchExist := false
	for _, codebase := range codebaseList.Items {
		if strings.ToLower(codebase.Spec.Type) == application {
			appExist = true
			branches, errBranch := h.NamespacedClient.CodebaseBranchesListByCodebaseName(ctx, codebase.Name)
			if errBranch != nil {
				logger.Error("cant get codebase branches", zap.Error(errBranch))
				InternalErrorResponse(ctx, writer, "cant get codebase branches")
				return
			}
			for _, branch := range branches {
				if branch.Status.Value == "active" {
					branchExist = true
				}
			}
		}
	}

	ciTool := "jenkins" // TODO need to investigate different option
	jenkins, err := h.NamespacedClient.GetJenkins(ctx, ciTool)
	if err != nil {
		logger.Error("cant get jenkins branches", zap.Error(err))
		InternalErrorResponse(ctx, writer, "cant get codebase branches")
		return
	}
	jenkinsURL := fmt.Sprintf("https://%s-%s.%s", ciTool, h.NamespacedClient.Namespace, jenkins.Spec.EdpSpec.DnsWildcard)

	var cdPipelines []cdPipelineForOverview
	for _, cdPipeCR := range cdPipelineList.Items {
		cdPipeline := cdPipelineForOverview{
			Name:        cdPipeCR.Name,
			Status:      cdPipeCR.Status.Value,
			JenkinsLink: util.CreateCICDPipelineLink(jenkinsURL, cdPipeCR.Name),
		}
		cdPipelines = append(cdPipelines, cdPipeline)
	}

	csrfToken := csrf.Token(request)
	tmplData := cdPipelinePageOverviewData{
		Username:                      user.UserName(),
		IsAdmin:                       user.IsAdmin(),
		Applications:                  appExist,
		ActiveApplicationsAndBranches: appExist && branchExist,
		CDPipelines:                   cdPipelines,
		Error:                         "",
		Xsrfdata:                      csrfToken,
		BasePath:                      h.Config.BasePath,
		Type:                          consts.CDMenuItem,
		DiagramPageEnabled:            h.Config.DiagramPageEnabled,
		EDPVersion:                    h.Config.EDPVersion,
		Success:                       true,
	}

	templatePaths := []string{path.Join(h.WorkingDir, "/viewsV2/continuous_delivery.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/header_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/navbar_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/footer_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/modal_success_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/delete_confirmation_template.html"),
	}

	template := Template{
		Data:             tmplData,
		FuncMap:          h.FuncMap,
		MainTemplateName: "continuous_delivery.html",
		TemplatePaths:    templatePaths,
	}

	OkHTMLResponse(ctx, writer, &template)

}

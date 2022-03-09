package webapi

import (
	"context"
	"net/http"
	"path"

	"edp-admin-console/internal/applog"
	"edp-admin-console/models/query"
)

const (
	application = "application"
)

type Codebase struct {
	Name      string
	Status    string
	Language  string
	BuildTool string
	CiTool    string
}

type createApplicationOverviewData struct {
	Type               query.CodebaseType
	Codebases          []Codebase
	BasePath           string
	EDPVersion         string
	Username           string
	Xsrfdata           string
	Error              string
	DeletionError      string
	IsAdmin            bool
	DiagramPageEnabled bool
	JiraEnabled        bool
	Success            bool
}

func (h *HandlerEnv) ApplicationOverview(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := applog.LoggerFromContext(ctx)
	user := UserFromContext(ctx)

	applications := make([]Codebase, 0)

	codebaseList, err := h.NamespacedClient.GetCodebaseList(context.Background())
	if err != nil {
		logger.Error(err.Error())
		InternalErrorResponse(ctx, writer, "getCodebaseList failed")
		return
	}

	for _, codebase := range codebaseList.Items {
		applications = append(applications, Codebase{
			Name:      codebase.Name,
			Status:    codebase.Status.Value,
			Language:  codebase.Spec.Lang,
			BuildTool: codebase.Spec.BuildTool,
			CiTool:    codebase.Spec.CiTool,
		})
	}

	xsrfData, _ := GetCookieByName(request, "_xsrf")

	tmplData := createApplicationOverviewData{
		BasePath:           h.Config.BasePath,
		EDPVersion:         h.Config.EDPVersion,
		Codebases:          applications,
		Type:               application,
		Username:           user.UserName(),
		IsAdmin:            user.IsAdmin(),
		DiagramPageEnabled: h.Config.DiagramPageEnabled,
		JiraEnabled:        true,
		Xsrfdata:           xsrfData,
		// Now read from Beego.flashData
		Error:         "",
		DeletionError: "",
		Success:       true,
	}

	templatePaths := []string{path.Join(h.WorkingDir, "/viewsV2/codebase.html"), path.Join(h.WorkingDir, "/views/template/footer_template.html"),
		path.Join(h.WorkingDir, "/views/template/header_template.html"), path.Join(h.WorkingDir, "/views/template/navbar_template.html"),
		path.Join(h.WorkingDir, "/views/template/modal_success_template.html"), path.Join(h.WorkingDir, "/views/template/delete_confirmation_template.html")}

	template := Template{
		Data:             tmplData,
		FuncMap:          h.FuncMap,
		MainTemplateName: "codebase.html",
		TemplatePaths:    templatePaths,
	}

	OkHTMLResponse(request.Context(), writer, &template)
}

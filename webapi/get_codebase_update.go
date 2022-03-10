package webapi

import (
	"net/http"
	"path"

	"edp-admin-console/internal/applog"
	"edp-admin-console/util/consts"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"go.uber.org/zap"
)

type EditCodebase struct {
	Name                     string
	Type                     string
	JiraServer               string
	CommitMessagePattern     string
	TicketNamePattern        string
	JiraIssueFields          bool
	JiraIssueMetadataPayload string
}

type JiraServer struct {
	Name string
}

type EditCodebaseTpl struct {
	BasePath            string
	Username            string
	Type                string // menu item
	Codebase            *EditCodebase
	CodebaseUpdateError bool
	JiraServer          []JiraServer
	Csrf                string
	// footer_template
	EDPVersion string
	// modal_success_template
	Success bool
	// navbar_template
	DiagramPageEnabled bool
}

func (h *HandlerEnv) GetCodebaseUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := applog.LoggerFromContext(ctx)

	user := UserFromContext(ctx)

	crCodebaseName := chi.URLParam(r, "codebaseName")
	crCodebase, err := h.NamespacedClient.GetCodebase(ctx, crCodebaseName)
	if err != nil {
		logger.Error("get cr codebase failed", zap.Error(err), zap.String("codebase_name", crCodebaseName))
		InternalErrorResponse(ctx, w, "get codebase failed")
		return
	}

	metaDataPayload := pointerToStr(crCodebase.Spec.JiraIssueMetadataPayload)
	jiraIssueFieldsPayload, err := PayloadToFieldMapWithExclude(metaDataPayload, []string{consts.IssuesLinksKey})
	if err != nil {
		logger.Error("convert jiraIssueMetadata to map failed", zap.Error(err), zap.String("payload", metaDataPayload))
		InternalErrorResponse(ctx, w, "get codebase branch list by name failed")
		return
	}
	jiraIssueFields := false
	if len(jiraIssueFieldsPayload) != 0 {
		jiraIssueFields = true
	}
	editCodebase := &EditCodebase{
		Name:                     crCodebase.GetName(),
		Type:                     crCodebase.Spec.Type,
		JiraServer:               pointerToStr(crCodebase.Spec.JiraServer),
		CommitMessagePattern:     pointerToStr(crCodebase.Spec.CommitMessagePattern),
		TicketNamePattern:        pointerToStr(crCodebase.Spec.TicketNamePattern),
		JiraIssueFields:          jiraIssueFields,
		JiraIssueMetadataPayload: metaDataPayload,
	}

	crJiraServers, err := h.NamespacedClient.JiraServersList(ctx)
	if err != nil {
		logger.Error("get jira servers list failed", zap.Error(err))
		InternalErrorResponse(ctx, w, "get jira servers list failed")
		return
	}

	jiraServers := make([]JiraServer, len(crJiraServers))
	for i, crJiraServer := range crJiraServers {
		jiraServers[i] = JiraServer{
			Name: crJiraServer.GetName(),
		}
	}

	csrfToken := csrf.Token(r)
	tplData := &EditCodebaseTpl{
		BasePath:            h.Config.BasePath,
		Username:            user.UserName(),
		Type:                application,
		Codebase:            editCodebase,
		CodebaseUpdateError: false, // TODO: investigate usage
		JiraServer:          jiraServers,
		Csrf:                csrfToken,
		EDPVersion:          h.Config.EDPVersion,
		Success:             true, // TODO: investigate usage
		DiagramPageEnabled:  h.Config.DiagramPageEnabled,
	}

	response := &Template{
		Data:             tplData,
		FuncMap:          h.FuncMap,
		MainTemplateName: "edit_codebase.html",
		TemplatePaths: []string{
			path.Join(h.WorkingDir, "/viewsV2/edit_codebase.html"),
			path.Join(h.WorkingDir, "/views/template/header_template.html"),
			path.Join(h.WorkingDir, "/views/template/navbar_template.html"),
			path.Join(h.WorkingDir, "/views/template/jira_advance_mapping_help_template.html"),
			path.Join(h.WorkingDir, "/views/template/footer_template.html"),
			path.Join(h.WorkingDir, "/views/template/modal_success_template.html"),
			path.Join(h.WorkingDir, "/views/template/jira_issue_metadata_template.html"),
		},
	}
	OkHTMLResponse(ctx, w, response)
}

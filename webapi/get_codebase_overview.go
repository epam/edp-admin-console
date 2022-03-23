package webapi

import (
	"encoding/json"
	"net/http"
	"path"
	"time"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"go.uber.org/zap"

	"edp-admin-console/internal/applog"
	"edp-admin-console/internal/edpcomponent"
	"edp-admin-console/util"
	"edp-admin-console/util/consts"
)

type codebaseBranch struct {
	Name             string
	Status           string
	BasePath         string
	VCSLink          string //idk
	Version          *string
	LastSuccessBuild *string
	CICDLink         string //idk
	Build            *string
}

type ActionLog struct {
	LastTimeUpdate time.Time
	UserName       string
	Message        string
	Action         string
	Result         string
}

type codebase struct {
	Name                 string
	CiTool               string
	Language             string
	EmptyProject         bool
	BuildTool            string
	Framework            *string
	Strategy             string
	DefaultBranch        string
	GitProjectPath       *string //idk
	TestReportFramework  *string
	GitUrl               string
	Description          string
	JobProvisioning      *string
	JenkinsSlave         *string
	DeploymentScript     *string
	VersioningType       string
	StartVersioningFrom  *string
	JiraServer           *string
	CommitMessagePattern *string
	TicketNamePattern    *string
	JiraIssueFields      map[string]interface{}
	Perf                 *codeBaseApi.Perf
	Status               string
	Type                 string
	CodebaseBranch       []codebaseBranch
	ActionLog            []ActionLog
}

type codebaseOverviewTpl struct {
	BasePath           string
	EDPVersion         string
	Username           string
	Codebase           codebase
	Xsrfdata           string
	Type               string
	DiagramPageEnabled bool
	TypeCaption        string
	JiraEnabled        bool
	TypeSingular       string
	Success            bool
	ErrorBranch        string
	IsAdmin            bool
}

func (h *HandlerEnv) GetCodebaseOverview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := applog.LoggerFromContext(ctx)
	user := UserFromContext(ctx)
	codebaseName := chi.URLParam(r, "codebaseName")
	codebaseCR, err := h.NamespacedClient.GetCodebase(ctx, codebaseName)
	if err != nil {
		logger.Error("find codebase CR failed", zap.String("codebase_name", codebaseName), zap.Error(err))
		BadRequestResponse(ctx, w, "cant find CR")
		return
	}

	jiraServersList, err := h.NamespacedClient.JiraServersList(ctx)
	if err != nil {
		logger.Error("get jira servers list failed", zap.Error(err))
		InternalErrorResponse(ctx, w, "get jira servers list failed")
		return
	}
	jiraEnabled := jiraServersList != nil

	gerritURL := ""
	gerritComponent := consts.Gerrit
	gerritEDPComponent, err := edpcomponent.ByNameIFExists(ctx, h.NamespacedClient, gerritComponent)
	if err != nil {
		logger.Error("get edp component failed", zap.String("edp_component", gerritComponent), zap.Error(err))
		InternalErrorResponse(ctx, w, "get edp component failed")
		return
	}
	if gerritEDPComponent != nil {
		gerritURL = gerritEDPComponent.Spec.Url
	}

	jenkinsURL := ""
	jenkinsComponent := consts.Jenkins
	jenkinsEDPComponent, err := edpcomponent.ByNameIFExists(ctx, h.NamespacedClient, consts.Jenkins)
	if err != nil {
		logger.Error("get edp component failed", zap.String("edp_component", jenkinsComponent), zap.Error(err))
		InternalErrorResponse(ctx, w, "get edp component failed")
		return
	}
	if jenkinsEDPComponent != nil {
		jenkinsURL = jenkinsEDPComponent.Spec.Url
	}

	crCbBranchList, err := h.NamespacedClient.CodebaseBranchesListByCodebaseName(ctx, codebaseCR.Name)
	if err != nil {
		logger.Error("get codebase_branch list by name failed", zap.Error(err), zap.String("codebase_name", codebaseName))
		InternalErrorResponse(ctx, w, "get codebase branch list by name failed")
		return
	}
	codebaseProjectURL := ""
	if codebaseCR.Spec.GitUrlPath != nil {
		codebaseProjectURL = *codebaseCR.Spec.GitUrlPath
	}
	codebaseBranches := make([]codebaseBranch, len(crCbBranchList))
	for i := range crCbBranchList {
		crCbBranch := crCbBranchList[i]
		cicdLink := getCiLink(codebaseCR, jenkinsURL, crCbBranch.Spec.BranchName, gerritURL)

		vcsLink := ""
		switch codebaseCR.Spec.GitServer {
		case consts.Gerrit:
			vcsLink = util.CreateGerritLink(gerritURL, codebaseCR.Name, crCbBranch.Spec.BranchName)
		default:
			vcsLink = util.CreateGitLink(gerritURL, codebaseProjectURL, crCbBranch.Name) // git URL should be used here, but i don't know where it is located
		}
		codebaseBranches[i] = codebaseBranch{
			Name:             crCbBranch.Spec.BranchName,
			Status:           crCbBranch.Status.Value,
			Version:          crCbBranch.Spec.Version,
			LastSuccessBuild: crCbBranch.Status.LastSuccessfulBuild,
			CICDLink:         cicdLink,
			VCSLink:          vcsLink,
			Build:            crCbBranch.Status.Build,
		}
	}

	metaDataPayload := pointerToStr(codebaseCR.Spec.JiraIssueMetadataPayload)
	jiraIssueFields, err := PayloadToFieldMapWithExclude(metaDataPayload, []string{consts.IssuesLinksKey})
	if err != nil {
		logger.Error("convert jiraIssueMetadata to map failed", zap.Error(err), zap.String("payload", metaDataPayload))
		InternalErrorResponse(ctx, w, "get codebase branch list by name failed")
		return
	}
	xsrf := csrf.Token(r)
	var tmplData = codebaseOverviewTpl{
		BasePath:           h.Config.BasePath,
		EDPVersion:         h.Config.BasePath,
		Username:           user.UserName(),
		Type:               codebaseCR.Spec.Type,
		TypeCaption:        codebaseCR.Spec.Type,
		DiagramPageEnabled: h.Config.DiagramPageEnabled,
		JiraEnabled:        jiraEnabled,
		TypeSingular:       codebaseCR.Spec.Type,
		Success:            true, //idk
		ErrorBranch:        "",
		IsAdmin:            user.IsAdmin(),
		Xsrfdata:           xsrf,
		Codebase: codebase{
			Name:                 codebaseCR.Name,
			CiTool:               codebaseCR.Spec.CiTool,
			Language:             codebaseCR.Spec.Lang,
			EmptyProject:         codebaseCR.Spec.EmptyProject,
			BuildTool:            codebaseCR.Spec.BuildTool,
			Framework:            codebaseCR.Spec.Framework,
			Strategy:             string(codebaseCR.Spec.Strategy),
			DefaultBranch:        codebaseCR.Spec.DefaultBranch,
			TestReportFramework:  codebaseCR.Spec.TestReportFramework,
			GitUrl:               pointerToStr(codebaseCR.Spec.GitUrlPath),
			Description:          pointerToStr(codebaseCR.Spec.Description),
			JobProvisioning:      codebaseCR.Spec.JobProvisioning,
			JenkinsSlave:         codebaseCR.Spec.JenkinsSlave,
			DeploymentScript:     &codebaseCR.Spec.DeploymentScript,
			VersioningType:       string(codebaseCR.Spec.Versioning.Type),
			StartVersioningFrom:  codebaseCR.Spec.Versioning.StartFrom,
			JiraServer:           codebaseCR.Spec.JiraServer,
			CommitMessagePattern: codebaseCR.Spec.CommitMessagePattern,
			TicketNamePattern:    codebaseCR.Spec.TicketNamePattern,
			JiraIssueFields:      jiraIssueFields,
			Perf:                 codebaseCR.Spec.Perf,
			Status:               codebaseCR.Status.Value,
			Type:                 codebaseCR.Spec.Type,
			CodebaseBranch:       codebaseBranches,
		},
	}
	templatePaths := []string{
		path.Join(h.WorkingDir, "/viewsV2/codebase_overview.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/footer_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/header_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/modal_success_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/navbar_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/delete_confirmation_template.html"),
	}

	template := &Template{
		Data:             tmplData,
		FuncMap:          h.FuncMap,
		MainTemplateName: "codebase_overview.html",
		TemplatePaths:    templatePaths,
	}

	OkHTMLResponse(ctx, w, template)
}

func getCiLink(codebase *codeBaseApi.Codebase, jenkinsHost, branch, gitHost string) string {
	if consts.JenkinsCITool == codebase.Spec.CiTool {
		return util.CreateCICDApplicationLink(jenkinsHost, codebase.Name, util.ProcessNameToKubernetesConvention(branch))
	}
	return util.CreateGitlabCILink(gitHost, *codebase.Spec.GitUrlPath)
}

func PayloadToFieldMapWithExclude(payload string, keysToDelete []string) (map[string]interface{}, error) {
	requestPayload := make(map[string]interface{})
	if len(payload) == 0 {
		return make(map[string]interface{}), nil
	}
	if err := json.Unmarshal([]byte(payload), &requestPayload); err != nil {
		return nil, err
	}
	for k := range requestPayload {
		if keysToDelete != nil && util.Contains(keysToDelete, k) {
			delete(requestPayload, k)
		}
	}
	return requestPayload, nil
}

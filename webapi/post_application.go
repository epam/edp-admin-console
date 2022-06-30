package webapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	"go.uber.org/zap"

	"edp-admin-console/internal/applog"
	"edp-admin-console/util"
)

func (h *HandlerEnv) CreateApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := applog.LoggerFromContext(ctx)
	logger.Debug("in handler")

	err := r.ParseForm()
	if err != nil {
		logger.Error("cant parse form", zap.Error(err))
		InternalErrorResponse(ctx, w, "cant parse form")
		return
	}
	lang := r.Form.Get("appLang")
	nameCR := r.Form.Get("appName")
	buildTool := r.Form.Get("buildTool")
	ciTool := r.Form.Get("ciTool")
	defaultBranch := r.Form.Get("defaultBranchName")
	deploymentScript := r.Form.Get("deploymentScript")

	commitMsgPattern := r.Form.Get("commitMessagePattern")
	framework := r.Form.Get("framework")
	jenkinsSlave := r.Form.Get("jenkinsSlave")
	jobProvisioning := r.Form.Get("jobProvisioning")
	ticketNamePattern := r.Form.Get("ticketNamePattern")
	versioningType := r.Form.Get("versioningType")

	startVersioningFrom := r.Form.Get("startVersioningFrom")
	snapshotField := r.Form.Get("snapshotStaticField")
	startFrom := util.GetVersionOrNil(startVersioningFrom, snapshotField)

	strategy := strings.ToLower(r.Form.Get("strategy"))
	gitServer := "gerrit"
	gitRepoPath := ""
	if strategy == "import" {
		gitServer = r.Form.Get("gitServer")
		gitRepoPath = r.Form.Get("gitRelativePath")
	}

	repository := r.Form.Get("gitRepoUrl")

	isMultiModule := strings.ToLower(r.Form.Get("isMultiModule"))
	if isMultiModule == "true" {
		framework = fmt.Sprintf("%s-multimodule", framework)
	}

	jiraServer := r.Form.Get("jiraServer")
	var juraIssueMetadataPayload *string

	if len(jiraServer) > 0 {
		payload, errExtract := extractJsonJiraIssueMetadataPayload(
			ParseStrSlice(r.Form, "jiraFieldName"), ParseStrSlice(r.Form, "jiraPattern"))
		if errExtract != nil {
			logger.Error("cant marshal juraIssueMetadataPayload", zap.Error(errExtract))
			InternalErrorResponse(ctx, w, "cant marshal juraIssueMetadataPayload")
			return
		}
		juraIssueMetadataPayload = payload

	}

	emptyProjectStr := strings.ToLower(r.Form.Get("isEmpty"))
	emptyProject := false
	if emptyProjectStr == "true" {
		emptyProject = true
	}

	var perf *codeBaseApi.Perf
	if perfServer := r.Form.Get("perfServer"); len(perfServer) > 0 {
		perf.Name = perfServer
		perf.DataSources = ParseStrSlice(r.Form, "dataSource")
	}

	codebaseSpec := codeBaseApi.CodebaseSpec{
		Lang:             lang,
		Framework:        strToPtr(framework),
		BuildTool:        buildTool,
		Strategy:         codeBaseApi.Strategy(strategy),
		Repository:       &codeBaseApi.Repository{Url: repository},
		Type:             "application",
		GitServer:        gitServer,
		GitUrlPath:       strToPtr(gitRepoPath),
		JenkinsSlave:     strToPtr(jenkinsSlave),
		JobProvisioning:  strToPtr(jobProvisioning),
		DeploymentScript: deploymentScript,
		Versioning: codeBaseApi.Versioning{
			Type:      codeBaseApi.VersioningType(versioningType),
			StartFrom: startFrom,
		},
		JiraServer:               strToPtr(jiraServer),
		CommitMessagePattern:     strToPtr(commitMsgPattern),
		TicketNamePattern:        strToPtr(ticketNamePattern),
		CiTool:                   ciTool,
		Perf:                     perf,
		DefaultBranch:            defaultBranch,
		JiraIssueMetadataPayload: juraIssueMetadataPayload,
		EmptyProject:             emptyProject,
	}

	err = h.NamespacedClient.CreateCodebaseByCustomFields(ctx, nameCR, codebaseSpec)
	if err != nil {
		logger.Error("cant create codebase", zap.Error(err))
		InternalErrorResponse(ctx, w, "cant create codebase CR")
		return
	}

	codeBaseBranchSpec := codeBaseApi.CodebaseBranchSpec{
		CodebaseName: nameCR,
		BranchName:   defaultBranch,
	}
	err = h.NamespacedClient.CreateCBBranchByCustomFields(ctx, fmt.Sprintf("%s-%s", nameCR, defaultBranch), codeBaseBranchSpec)
	if err != nil {
		logger.Error("cant create codebase branch", zap.Error(err))
		InternalErrorResponse(ctx, w, "cant create codebase branch")
		return
	}
	http.Redirect(w, r, fmt.Sprintf("%s/v2/admin/edp/application/overview?%s=%s#codebaseSuccessModal", h.Config.BasePath, "waitingforcodebase", nameCR), http.StatusFound)
}

func strToPtr(s string) *string {
	if len(s) == 0 || s == "null" {
		return nil
	}
	return &s
}

func extractJsonJiraIssueMetadataPayload(jiraFieldNames, jiraPatterns []string) (*string, error) {
	if jiraFieldNames == nil && jiraPatterns == nil {
		return nil, nil
	}
	if len(jiraFieldNames) != len(jiraPatterns) {
		return nil, fmt.Errorf("jiraFieldNames (len = %v) and jiraPatterns (len = %v) are not the same size",
			len(jiraFieldNames), len(jiraPatterns))
	}
	payload := make(map[string]string, len(jiraFieldNames))
	for i, name := range jiraFieldNames {
		payload[name] = jiraPatterns[i]
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	jsonStr := string(jsonPayload)
	return &jsonStr, nil
}

func ParseStrSlice(form url.Values, key string) []string {
	i := 0
	var s []string

	for form.Get(fmt.Sprintf("%s.%v", key, i)) != "" {
		s = append(s, form.Get(fmt.Sprintf("%s.%v", key, i)))
		i++
	}
	return s
}

package webapi

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/prometheus/common/log"
	"go.uber.org/zap"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	jenkinsApi "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1"

	"edp-admin-console/internal/applog"
	"edp-admin-console/k8s"
	"edp-admin-console/models/query"
)

const (
	scope                           = "ci"
	CodeBaseIntegrationStrategyIsOn = true
	edpConfigMapName                = "edp-config"
	perfIntegrationEnabledKey       = "perf_integration_enabled"
)

type gitServer struct {
	Name      string
	Hostname  string
	Available bool
}

type jenkinsSlave struct {
	Name string
}

type perfServer struct {
	Name string
}

type jiraServer struct {
	Name      string
	Available bool
}

type jobProvisioner struct {
	Name  string
	Scope string
}

type createApplicationData struct {
	Type                        query.CodebaseType
	BasePath                    string
	EDPVersion                  string
	Username                    string
	Xsrfdata                    string
	Error                       string
	DeletionError               string
	IsOpenshift                 bool
	DiagramPageEnabled          bool
	Success                     bool
	IsPerfEnabled               bool
	IsVcsEnabled                bool
	CodeBaseIntegrationStrategy bool
	GitServers                  []*gitServer
	JenkinsSlaves               []*jenkinsSlave
	JobProvisioners             []*jobProvisioner
	JiraServer                  []*jiraServer
	PerfServer                  []*perfServer
	IntegrationStrategies       []string
	BuildTools                  []string
	VersioningTypes             []string
	DeploymentScripts           []string
	PerfDataSources             []string
	CiTools                     []string
}

func (h *HandlerEnv) CreateApplicationPage(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := applog.LoggerFromContext(ctx)
	user := UserFromContext(ctx)

	gitServersArray, err := getGitServers(ctx, h.NamespacedClient)
	if err != nil {
		log.Error("getGitServers failed", zap.Error(err))
		InternalErrorResponse(ctx, writer, "getGitServers failed")
		return
	}

	jenkinsSlavesArray, err := getJenkinsSlaves(ctx, h.NamespacedClient)
	if err != nil {
		logger.Error("getJenkinsSlaves failed", zap.Error(err))
		InternalErrorResponse(ctx, writer, "getJenkinsSlaves failed")
		return
	}

	jiraServersArray, err := getJiraServers(ctx, h.NamespacedClient)
	if err != nil {
		logger.Error("getJiraServers failed", zap.Error(err))
		InternalErrorResponse(ctx, writer, "getJiraServers failed")
		return
	}

	perfServersArray, err := getPerfServers(ctx, h.NamespacedClient)
	if err != nil {
		logger.Error("getPerfServers failed", zap.Error(err))
		InternalErrorResponse(ctx, writer, "getPerfServers failed")
		return
	}

	ciJobProvisionersArray, err := getJobProvisionsWithScope(ctx, h.NamespacedClient, scope)
	if err != nil {
		logger.Error("getJobProvisionsWithScope failed", zap.Error(err))
		InternalErrorResponse(ctx, writer, "getJobProvisionsWithScope failed")
		return
	}

	perfEnabled, err := isPerfEnabled(ctx, h.NamespacedClient, h.NamespacedClient.Namespace)
	if err != nil {
		logger.Error("MESSAGE", zap.Error(err))
		InternalErrorResponse(ctx, writer, "isPerfEnabled failed")
		return
	}

	xsrfData := csrf.Token(request)
	tmplData := createApplicationData{
		Type:                        application,
		BasePath:                    h.Config.BasePath,
		EDPVersion:                  h.Config.EDPVersion,
		Username:                    user.UserName(),
		Xsrfdata:                    xsrfData,
		Error:                       "", // flashData
		DeletionError:               "", // flashData
		IsOpenshift:                 h.Config.IsOpenshift,
		DiagramPageEnabled:          h.Config.DiagramPageEnabled,
		Success:                     true, // flashData
		IsPerfEnabled:               perfEnabled,
		IsVcsEnabled:                h.Config.IsVcsIntegrationEnabled,
		CodeBaseIntegrationStrategy: CodeBaseIntegrationStrategyIsOn,
		GitServers:                  gitServersArray,
		JenkinsSlaves:               jenkinsSlavesArray,
		JobProvisioners:             ciJobProvisionersArray,
		JiraServer:                  jiraServersArray,
		PerfServer:                  perfServersArray,
		IntegrationStrategies:       h.Config.Reference.IntegrationStrategies,
		BuildTools:                  h.Config.Reference.BuildTools,
		VersioningTypes:             h.Config.Reference.VersioningTypes,
		DeploymentScripts:           h.Config.Reference.DeploymentScript,
		PerfDataSources:             h.Config.Reference.PerfDataSources,
		CiTools:                     h.Config.Reference.CiTools,
	}

	templatePaths := []string{path.Join(h.WorkingDir, "/viewsV2/create_application.html"), path.Join(h.WorkingDir, "/views/template/footer_template.html"),
		path.Join(h.WorkingDir, "/views/template/header_template.html"), path.Join(h.WorkingDir, "/views/template/navbar_template.html"),
		path.Join(h.WorkingDir, "/views/template/accordion_codebase_template.html"), path.Join(h.WorkingDir, "/views/template/default_branch_template.html"),
		path.Join(h.WorkingDir, "/views/template/empty_project_template.html"), path.Join(h.WorkingDir, "/views/template/language_template.html"),
		path.Join(h.WorkingDir, "/views/template/java_framework_template.html"), path.Join(h.WorkingDir, "/views/template/java_script_framework_template.html"),
		path.Join(h.WorkingDir, "/views/template/dotnet_framework_template.html"), path.Join(h.WorkingDir, "/views/template/go_framework_template.html"),
		path.Join(h.WorkingDir, "/views/template/python_framework_template.html"), path.Join(h.WorkingDir, "/views/template/terraform_framework_template.html"),
		path.Join(h.WorkingDir, "/views/template/build_tool_template.html"), path.Join(h.WorkingDir, "/viewsV2/advanced_settings_block_template.html"),
		path.Join(h.WorkingDir, "/views/template/data_source_block_template.html"), path.Join(h.WorkingDir, "/views/template/accordion_vcs_template.html"),
		path.Join(h.WorkingDir, "/views/template/confirmation_popup_template.html"), path.Join(h.WorkingDir, "/views/template/jira_issue_metadata_template.html"),
		path.Join(h.WorkingDir, "/views/template/jira_advance_mapping_help_template.html"), path.Join(h.WorkingDir, "/views/template/perf_template.html")}

	template := Template{
		Data:             tmplData,
		FuncMap:          h.FuncMap,
		MainTemplateName: "create_application.html",
		TemplatePaths:    templatePaths,
	}

	OkHTMLResponse(ctx, writer, &template)
}

func getGitServers(ctx context.Context, client *k8s.RuntimeNamespacedClient) ([]*gitServer, error) {
	gitServers, err := client.GetGitServerList(ctx)
	if err != nil {
		return nil, err
	}

	gitServersArray := make([]*gitServer, 0)
	for _, gitServerCR := range gitServers.Items {
		gitServersArray = append(gitServersArray, &gitServer{
			Name:      gitServerCR.Name,
			Hostname:  gitServerCR.Spec.GitHost,
			Available: gitServerCR.Status.Available,
		})
	}

	return gitServersArray, nil
}

func getJenkinsSlaves(ctx context.Context, client *k8s.RuntimeNamespacedClient) ([]*jenkinsSlave, error) {
	jenkinsCRList, err := client.GetJenkinsList(ctx)
	if err != nil {
		return nil, err
	}

	jenkinsSlavesArray := make([]jenkinsApi.Slave, 0)
	for _, jenkinsCR := range jenkinsCRList.Items {
		jenkinsSlavesArray = append(jenkinsSlavesArray, jenkinsCR.Status.Slaves...)
	}

	jenkinsSlaves := make([]*jenkinsSlave, 0)
	for _, jenkinsCRSlave := range jenkinsSlavesArray {
		jenkinsSlaves = append(jenkinsSlaves, &jenkinsSlave{Name: jenkinsCRSlave.Name})
	}

	return jenkinsSlaves, err
}

func getJiraServers(ctx context.Context, client *k8s.RuntimeNamespacedClient) ([]*jiraServer, error) {
	gitServers, err := client.JiraServersList(ctx)
	if err != nil {
		return nil, err
	}

	jiraServersArray := make([]*jiraServer, 0)
	for _, jiraServerCR := range gitServers {
		jiraServersArray = append(jiraServersArray, &jiraServer{
			Name:      jiraServerCR.Name,
			Available: jiraServerCR.Status.Available,
		})
	}

	return jiraServersArray, nil
}

func getPerfServers(ctx context.Context, client *k8s.RuntimeNamespacedClient) ([]*perfServer, error) {
	perfServerList, err := client.GetPerfServerList(ctx)
	if err != nil {
		return nil, err
	}

	perfServersArray := make([]*perfServer, 0)
	for _, perfServerCR := range perfServerList.Items {
		perfServersArray = append(perfServersArray, &perfServer{
			Name: perfServerCR.Name,
		})
	}

	return perfServersArray, nil
}

func getJobProvisionsWithScope(ctx context.Context, client *k8s.RuntimeNamespacedClient, scope string) ([]*jobProvisioner, error) {
	jenkinsCRList, err := client.GetJenkinsList(ctx)
	if err != nil {
		return nil, err
	}

	jenkinsJobProvisioners := make([]jenkinsApi.JobProvision, 0)
	for _, jenkinsCR := range jenkinsCRList.Items {
		jenkinsJobProvisioners = append(jenkinsJobProvisioners, jenkinsCR.Status.JobProvisions...)
	}

	jobProvisionersArray := make([]*jobProvisioner, 0)
	for _, jenkinsCRJobProvisioners := range jenkinsJobProvisioners {
		if jenkinsCRJobProvisioners.Scope == scope {
			jobProvisionersArray = append(jobProvisionersArray, &jobProvisioner{
				Name:  jenkinsCRJobProvisioners.Name,
				Scope: jenkinsCRJobProvisioners.Scope,
			})
		}
	}

	return jobProvisionersArray, err
}

func isPerfEnabled(ctx context.Context, client *k8s.RuntimeNamespacedClient, namespace string) (bool, error) {
	configMap := coreV1.ConfigMap{}
	err := client.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      edpConfigMapName,
	}, &configMap)

	isPerfIntegrationEnabled := configMap.Data[perfIntegrationEnabledKey]
	if len(isPerfIntegrationEnabled) == 0 {
		return false, fmt.Errorf("there is no key %s in Config Map edp-config", perfIntegrationEnabledKey)
	}
	isPerfEnabled, err := strconv.ParseBool(isPerfIntegrationEnabled)
	if err != nil {
		return false, err
	}
	return isPerfEnabled, nil
}

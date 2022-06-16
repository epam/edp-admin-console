package webapi

import (
	"context"
	"net/http"
	"path"

	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	"github.com/gorilla/csrf"
	"go.uber.org/zap"

	"edp-admin-console/internal/applications"
	"edp-admin-console/internal/applog"
	"edp-admin-console/internal/imagestream"
	"edp-admin-console/internal/jobprovisioners"
	"edp-admin-console/k8s"
	"edp-admin-console/util/consts"
)

type CreateCDPage struct {
	// footer_template
	EDPVersion string
	// modal_success_template
	Success bool
	// navbar_template
	DiagramPageEnabled bool
	BasePath           string
	Username           string
	Type               string // menu item
	Error              string
	Apps               []Application
	Csrf               string
	Autotests          []Autotest
	GroovyLibs         []GroovyLib
	JobProvisioners    []JobProvisioner
}

type Application struct {
	CodebaseBranch []CodebaseBranchApp
	Name           string
}

type CodebaseBranchApp struct {
	CodebaseDockerStream []CodebaseDockerStream
}

type CodebaseDockerStream struct {
	OcImageStreamName string
}

type Autotest struct {
	Name           string
	CodebaseBranch []CodebaseBranchAutotest
}

type CodebaseBranchAutotest struct {
	Name string
}

type GroovyLib struct {
	Name           string
	CodebaseBranch []CodebaseBranchGroovyLib
}

type CodebaseBranchGroovyLib struct {
	Name string
}

type JobProvisioner struct {
	Name string
}

type activeCodebase struct { // is a helper struct, used only in this handler
	app         codeBaseApi.Codebase
	appBranches []*codeBaseApi.CodebaseBranch
}

func (h *HandlerEnv) GetCDCreatePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := applog.LoggerFromContext(ctx)
	user := UserFromContext(ctx)
	k8sClient := h.NamespacedClient

	cbList, err := k8sClient.GetCodebaseList(ctx)
	if err != nil {
		logger.Error("get codebases list failed", zap.Error(err))
		InternalErrorResponse(ctx, w, "get codebases list failed")
		return
	}

	apps, err := fetchApplications(ctx, k8sClient, cbList.Items)
	if err != nil {
		InternalErrorResponse(ctx, w, err.Error())
		return
	}

	autotests, err := fetchAutotests(ctx, k8sClient, cbList.Items)
	if err != nil {
		InternalErrorResponse(ctx, w, err.Error())
		return
	}

	libs, err := fetchLibs(ctx, k8sClient, cbList.Items)
	if err != nil {
		InternalErrorResponse(ctx, w, err.Error())
		return
	}

	jobProvisioners, err := fetchJobProvisioners(ctx, k8sClient)
	if err != nil {
		InternalErrorResponse(ctx, w, err.Error())
		return
	}

	csrfToken := csrf.Token(r)
	tplData := &CreateCDPage{
		Username:           user.UserName(),
		DiagramPageEnabled: h.Config.DiagramPageEnabled,
		Type:               consts.CDMenuItem,
		EDPVersion:         h.Config.EDPVersion,
		Apps:               apps,
		Csrf:               csrfToken,
		Autotests:          autotests,
		GroovyLibs:         libs,
		JobProvisioners:    jobProvisioners,
	}

	templatePaths := []string{
		path.Join(h.WorkingDir, "/viewsV2/create_cd_pipeline.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/footer_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/header_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/modal_success_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/navbar_template.html"),
		path.Join(h.WorkingDir, "/viewsV2/template/delete_confirmation_template.html"),
	}

	template := &Template{
		Data:             tplData,
		FuncMap:          h.FuncMap,
		MainTemplateName: "create_cd_pipeline.html",
		TemplatePaths:    templatePaths,
	}

	OkHTMLResponse(ctx, w, template)
}

func fetchApplications(ctx context.Context, k8sClient *k8s.RuntimeNamespacedClient, cbList []codeBaseApi.Codebase) ([]Application, error) {
	logger := applog.LoggerFromContext(ctx)

	activeAppsCR, err := applications.ActiveApplications(cbList)
	if err != nil {
		logger.Error("get active applications failed", zap.Error(err))
		return nil, err
	}

	appsActiveBranches := make([]*codeBaseApi.CodebaseBranch, 0)
	activeApps := make([]activeCodebase, 0)
	activeAppNames := make([]string, 0)
	for i := range activeAppsCR {
		codebaseName := activeAppsCR[i].Name
		activeBranches, branchErr := applications.ActiveCodebaseBranches(ctx, k8sClient, codebaseName)
		if branchErr != nil {
			logger.Error("get applications active branches failed", zap.Error(err), zap.String("codebase_name", codebaseName))
			return nil, err
		}
		appsActiveBranches = append(appsActiveBranches, activeBranches...)
		activeApps = append(activeApps, activeCodebase{
			app:         activeAppsCR[i],
			appBranches: activeBranches,
		})
		activeAppNames = append(activeAppNames, codebaseName)
	}

	cbStreams := make(map[string][]codeBaseApi.CodebaseImageStream)
	if len(activeAppNames) > 0 {
		var streamsErr error
		cbStreams, streamsErr = imagestream.OutputCBStreamsForCodebaseNames(ctx, k8sClient, activeAppNames)
		if streamsErr != nil {
			logger.Error("get codebase image streams for the active apps failed", zap.Error(err), zap.Strings("codebase_names", activeAppNames))
			return nil, err
		}
	}

	apps := make([]Application, 0)
	for _, activeApplication := range activeApps {
		codebaseBranches := make([]CodebaseBranchApp, 0)
		outputCBStreams := make([]CodebaseDockerStream, 0)
		if appCBStreams, ok := cbStreams[activeApplication.app.Name]; ok {
			for _, appCBStreamCR := range appCBStreams {
				outputCBStreams = append(outputCBStreams, CodebaseDockerStream{
					OcImageStreamName: appCBStreamCR.Name,
				})
			}
		}
		codebaseBranches = append(codebaseBranches, CodebaseBranchApp{
			CodebaseDockerStream: outputCBStreams,
		})
		tplApp := Application{
			CodebaseBranch: codebaseBranches,
			Name:           activeApplication.app.Name,
		}
		apps = append(apps, tplApp)
	}
	return apps, err
}

func fetchAutotests(ctx context.Context, k8sClient *k8s.RuntimeNamespacedClient, cbList []codeBaseApi.Codebase) ([]Autotest, error) {
	logger := applog.LoggerFromContext(ctx)

	activeAutotestsCR, err := applications.ActiveAutotests(cbList)
	if err != nil {
		logger.Error("get active autotests failed", zap.Error(err))
		return nil, err
	}

	appsActiveBranches := make([]*codeBaseApi.CodebaseBranch, 0)
	activeApps := make([]activeCodebase, 0)
	for i := range activeAutotestsCR {
		codebaseName := activeAutotestsCR[i].Name
		activeBranches, branchErr := applications.ActiveCodebaseBranches(ctx, k8sClient, codebaseName)
		if branchErr != nil {
			logger.Error("get applications active branches failed", zap.Error(err), zap.String("codebase_name", codebaseName))
			return nil, err
		}
		appsActiveBranches = append(appsActiveBranches, activeBranches...)
		activeApps = append(activeApps, activeCodebase{
			app:         activeAutotestsCR[i],
			appBranches: activeBranches,
		})
	}

	autotests := make([]Autotest, 0)
	for _, activeApplication := range activeApps {
		codebaseBranches := make([]CodebaseBranchAutotest, 0)
		for _, cbBranches := range activeApplication.appBranches {
			codebaseBranches = append(codebaseBranches, CodebaseBranchAutotest{
				Name: cbBranches.Spec.BranchName,
			})
		}
		tplApp := Autotest{
			CodebaseBranch: codebaseBranches,
			Name:           activeApplication.app.Name,
		}
		autotests = append(autotests, tplApp)
	}
	return autotests, err
}

func fetchLibs(ctx context.Context, k8sClient *k8s.RuntimeNamespacedClient, cbList []codeBaseApi.Codebase) ([]GroovyLib, error) {
	logger := applog.LoggerFromContext(ctx)

	activeGroovyLibsCR, err := applications.ActiveGroovyLibs(cbList)
	if err != nil {
		logger.Error("get active groovy libs failed", zap.Error(err))
		return nil, err
	}

	appsActiveBranches := make([]*codeBaseApi.CodebaseBranch, 0)
	activeCodebases := make([]activeCodebase, 0)
	for i := range activeGroovyLibsCR {
		codebaseName := activeGroovyLibsCR[i].Name
		activeBranches, branchErr := applications.ActiveCodebaseBranches(ctx, k8sClient, codebaseName)
		if branchErr != nil {
			logger.Error("get applications active branches failed", zap.Error(err), zap.String("codebase_name", codebaseName))
			return nil, err
		}
		appsActiveBranches = append(appsActiveBranches, activeBranches...)
		activeCodebases = append(activeCodebases, activeCodebase{
			app:         activeGroovyLibsCR[i],
			appBranches: activeBranches,
		})
	}

	autotests := make([]GroovyLib, 0)
	for _, activeCb := range activeCodebases {
		codebaseBranches := make([]CodebaseBranchGroovyLib, 0)
		for _, cbBranches := range activeCb.appBranches {
			codebaseBranches = append(codebaseBranches, CodebaseBranchGroovyLib{
				Name: cbBranches.Spec.BranchName,
			})
		}
		tplApp := GroovyLib{
			CodebaseBranch: codebaseBranches,
			Name:           activeCb.app.Name,
		}
		autotests = append(autotests, tplApp)
	}
	return autotests, err
}

func fetchJobProvisioners(ctx context.Context, k8sClient *k8s.RuntimeNamespacedClient) ([]JobProvisioner, error) {
	jobProvisionerNames, err := jobprovisioners.ListNames(ctx, k8sClient)
	if err != nil {
		return nil, err
	}
	jobProvisioners := make([]JobProvisioner, 0)
	for _, name := range jobProvisionerNames {
		jobProvisioners = append(jobProvisioners, JobProvisioner{Name: name})
	}
	return jobProvisioners, nil
}

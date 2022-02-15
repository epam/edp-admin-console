/*
 * Copyright 2020 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package webapi

import (
	"fmt"
	"path"

	"github.com/astaxie/beego"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"edp-admin-console/context"
	"edp-admin-console/controllers"
	"edp-admin-console/controllers/auth"
	cdPipeController "edp-admin-console/controllers/cd-pipeline"
	"edp-admin-console/controllers/stage"
	"edp-admin-console/filters"
	"edp-admin-console/k8s"
	"edp-admin-console/repository"
	edpComponentRepo "edp-admin-console/repository/edp-component"
	jirarepo "edp-admin-console/repository/jira-server"
	perfRepo "edp-admin-console/repository/perfboard"
	"edp-admin-console/service"
	"edp-admin-console/service/cd_pipeline"
	cbs "edp-admin-console/service/codebasebranch"
	edpComponentService "edp-admin-console/service/edp-component"
	jiraservice "edp-admin-console/service/jira-server"
	"edp-admin-console/service/logger"
	"edp-admin-console/service/perfboard"
	"edp-admin-console/util"
)

var zaplog = logger.GetLogger()

const (
	integrationStrategies = "integrationStrategies"
	buildTools            = "buildTools"
	versioningTypes       = "versioningTypes"
	testReportTools       = "testReportTools"
	deploymentScript      = "deploymentScript"
	ciTools               = "ciTools"
	perfDataSources       = "perfDataSources"

	CreateStrategy = "Create"
	apiV2Scope     = "/api/v2"
	edpScope       = "/edp"
	edpScopeV2     = "/v2/admin/edp"
)

func SetupRouter(namespacedClient *k8s.RuntimeNamespacedClient, workingDir string) {
	zaplog.Info("Start application...",
		zap.String("mode", beego.AppConfig.String("runmode")),
		zap.String("edp version", context.EDPVersion))
	authEnabled, err := beego.AppConfig.Bool("keycloakAuthEnabled")
	if err != nil {
		zaplog.Error("Cannot read property keycloakAuthEnabled. Set default: false", zap.Error(err))
		authEnabled = false
	}

	permissions := filters.PermissionsMap()
	accessHandlerEnv := &filters.AccessControlEnv{
		Permissions: permissions,
	}

	if authEnabled {
		context.InitAuth()
		beego.Router(fmt.Sprintf("%s/auth/callback", context.BasePath), &auth.AuthController{}, "get:Callback")
		beego.InsertFilter(fmt.Sprintf("%s/admin/*", context.BasePath), beego.BeforeRouter, filters.AuthFilter)
		beego.InsertFilter(fmt.Sprintf("%s/api/v1/edp/*", context.BasePath), beego.BeforeRouter, filters.AuthRestFilter)
		beego.InsertFilter(fmt.Sprintf("%s/admin/edp/*", context.BasePath), beego.BeforeRouter, accessHandlerEnv.RoleAccessControlFilter)
		beego.InsertFilter(fmt.Sprintf("%s/api/v1/edp/*", context.BasePath), beego.BeforeRouter, accessHandlerEnv.RoleAccessControlRestFilter)
		// auth and role access for v2 api
		beego.InsertFilter(fmt.Sprintf("%s%s%s/*", context.BasePath, apiV2Scope, edpScope), beego.BeforeRouter, filters.AuthRestFilter)
		beego.InsertFilter(fmt.Sprintf("%s%s%s/*", context.BasePath, apiV2Scope, edpScope), beego.BeforeRouter, accessHandlerEnv.RoleAccessControlRestFilter)
	} else {
		beego.InsertFilter(fmt.Sprintf("%s/*", context.BasePath), beego.BeforeRouter, filters.StubAuthFilter)
	}

	dbEnable, err := beego.AppConfig.Bool("dbEnabled")
	if err != nil {
		zaplog.Error("Cannot read property dbEnabled. Set default: true", zap.Error(err))
		dbEnable = true
	}

	if dbEnable {
		context.InitDb()
	}

	k8sClientConf := k8s.ClientConfig()
	clients := k8s.CreateOpenShiftClients(k8sClientConf)
	codebaseRepository := repository.CodebaseRepository{}
	branchRepository := repository.CodebaseBranchRepository{}
	pipelineRepository := repository.CDPipelineRepository{}
	gitServerRepository := repository.GitServerRepository{}
	sr := repository.SlaveRepository{}
	pr := repository.JobProvisioning{}
	ecr := edpComponentRepo.EDPComponent{}
	jsr := jirarepo.JiraServer{}
	psr := perfRepo.PerfServer{}

	gitServerService := service.GitServerService{IGitServerRepository: gitServerRepository}
	ss := service.SlaveService{ISlaveRepository: sr}
	ps := service.JobProvisioning{IJobProvisioningRepository: pr}
	js := jiraservice.JiraServer{IJiraServer: jsr}
	pbs := perfboard.PerfBoard{
		PerfRepo:   psr,
		CoreClient: clients.CoreClient,
	}
	ecs := edpComponentService.EDPComponentService{IEDPComponent: ecr}
	edpService := service.EDPTenantService{Clients: clients}
	clusterService := service.ClusterService{Clients: clients}
	branchService := cbs.CodebaseBranchService{
		Clients:                  clients,
		IReleaseBranchRepository: branchRepository,
		ICDPipelineRepository:    pipelineRepository,
		ICodebaseRepository:      codebaseRepository,
		CodebaseBranchValidation: map[string]func(string, string) ([]string, error){
			"application": pipelineRepository.GetCDPipelinesUsingApplicationAndBranch,
			"autotests":   pipelineRepository.GetCDPipelinesUsingAutotestAndBranch,
			"library":     pipelineRepository.GetCDPipelinesUsingLibraryAndBranch,
		},
	}
	codebaseService := service.CodebaseService{
		Clients:               clients,
		ICodebaseRepository:   codebaseRepository,
		ICDPipelineRepository: pipelineRepository,
		BranchService:         branchService,
		PerfService:           pbs,
	}
	pipelineService := cd_pipeline.CDPipelineService{
		Clients:               clients,
		ICDPipelineRepository: pipelineRepository,
		CodebaseService:       codebaseService,
		BranchService:         branchService,
		EDPComponent:          ecs,
	}

	beego.ErrorController(&controllers.ErrorController{})
	beego.Router(fmt.Sprintf("%s/", context.BasePath), &controllers.MainController{EDPTenantService: edpService}, "get:Index")
	beego.SetStaticPath(fmt.Sprintf("%s/static", context.BasePath), "static")

	integrationStrategies := util.GetValuesFromConfig(integrationStrategies)
	if integrationStrategies == nil {
		zaplog.Fatal("integrationStrategies config variable is empty.")
	}

	buildTools := util.GetValuesFromConfig(buildTools)
	if buildTools == nil {
		zaplog.Fatal("buildTools config variable is empty.")
	}

	vt := util.GetValuesFromConfig(versioningTypes)
	if vt == nil {
		zaplog.Fatal("versioningTypes config variable is empty.")
	}

	testReportTools := util.GetValuesFromConfig(testReportTools)
	if testReportTools == nil {
		zaplog.Fatal("testReportTools config variable is empty.")
	}

	ds := util.GetValuesFromConfig(deploymentScript)
	if ds == nil {
		zaplog.Fatal("deploymentScript config variable is empty.")
	}

	ciTools := util.GetValuesFromConfig(ciTools)
	if ciTools == nil {
		zaplog.Fatal("ciTools config variable is empty.")
	}

	perfDataSources := util.GetValuesFromConfig(perfDataSources)
	if perfDataSources == nil {
		zaplog.Fatal("perfDataSources config variable is empty.")
	}

	is := make([]string, len(integrationStrategies))
	copy(is, integrationStrategies)

	appc := controllers.ApplicationController{
		CodebaseService:  codebaseService,
		EDPTenantService: edpService,
		BranchService:    branchService,
		GitServerService: gitServerService,
		SlaveService:     ss,
		JobProvisioning:  ps,
		JiraServer:       js,
		PerfService:      pbs,

		IntegrationStrategies: is,
		BuildTools:            buildTools,
		VersioningTypes:       vt,
		DeploymentScript:      ds,
		CiTools:               ciTools,
		PerfDataSources:       perfDataSources,
	}

	autis := make([]string, len(integrationStrategies))
	copy(autis, integrationStrategies)

	autc := controllers.AutotestsController{
		EDPTenantService: edpService,
		CodebaseService:  codebaseService,
		BranchService:    branchService,
		GitServerService: gitServerService,
		SlaveService:     ss,
		JobProvisioning:  ps,
		JiraServer:       js,
		PerfService:      pbs,

		IntegrationStrategies: util.RemoveElByValue(autis, CreateStrategy),
		BuildTools:            buildTools,
		VersioningTypes:       vt,
		TestReportTools:       testReportTools,
		DeploymentScript:      ds,
		CiTools:               ciTools,
		PerfDataSources:       perfDataSources,
	}

	lc := controllers.LibraryController{
		EDPTenantService: edpService,
		CodebaseService:  codebaseService,
		GitServerService: gitServerService,
		SlaveService:     ss,
		JobProvisioning:  ps,
		JiraServer:       js,
		PerfService:      pbs,

		IntegrationStrategies: is,
		BuildTools:            buildTools,
		VersioningTypes:       vt,
		DeploymentScript:      ds,
		CiTools:               ciTools,
		PerfDataSources:       perfDataSources,
	}

	ec := controllers.EDPTenantController{
		EDPTenantService: edpService,
		EDPComponent:     ecs,
	}

	cc := controllers.CodebaseController{
		CodebaseService:  codebaseService,
		EDPTenantService: edpService,
		BranchService:    branchService,
		GitServerService: gitServerService,
		EDPComponent:     ecs,
		JiraServer:       js,
	}

	cpc := cdPipeController.CDPipelineController{
		CodebaseService:  codebaseService,
		PipelineService:  pipelineService,
		EDPTenantService: edpService,
		BranchService:    branchService,
		EDPComponent:     ecs,
		JobProvisioning:  ps,
	}

	cbc := controllers.BranchController{
		BranchService:   branchService,
		CodebaseService: codebaseService,
	}

	dc := controllers.DiagramController{
		CodebaseService: codebaseService,
		PipelineService: pipelineService,
	}

	csc := stage.InitStageController(pipelineService)

	adminEdpNamespace := beego.NewNamespace(fmt.Sprintf("%s/admin/edp", context.BasePath),
		beego.NSRouter("/overview", &ec, "get:GetEDPComponents"),
		beego.NSRouter("/application/overview", &appc, "get:GetApplicationsOverviewPage"),
		beego.NSRouter("/application/create", &appc, "get:GetCreateApplicationPage"),
		beego.NSRouter("/application", &appc, "post:CreateApplication"),

		beego.NSRouter("/cd-pipeline/overview", &cpc, "get:GetContinuousDeliveryPage"),
		beego.NSRouter("/cd-pipeline/create", &cpc, "get:GetCreateCDPipelinePage"),
		beego.NSRouter("/cd-pipeline/:name/update", &cpc, "get:GetEditCDPipelinePage"),
		beego.NSRouter("/cd-pipeline", &cpc, "post:CreateCDPipeline"),
		beego.NSRouter("/cd-pipeline/:name/update", &cpc, "post:UpdateCDPipeline"),
		beego.NSRouter("/cd-pipeline/:pipelineName/overview", &cpc, "get:GetCDPipelineOverviewPage"),
		beego.NSRouter("/cd-pipeline/:pipelineName/cd-stage/edit", &csc, "post:UpdateCDStage"),
		beego.NSRouter("/autotest/overview", &autc, "get:GetAutotestsOverviewPage"),
		beego.NSRouter("/autotest/create", &autc, "get:GetCreateAutotestsPage"),
		beego.NSRouter("/autotest", &autc, "post:CreateAutotests"),

		beego.NSRouter("/codebase/:codebaseName/overview", &cc, "get:GetCodebaseOverviewPage"),
		beego.NSRouter("/codebase", &cc, "post:Delete"),
		beego.NSRouter("/codebase/branch/delete", &cbc, "post:Delete"),
		beego.NSRouter("/codebase/:name/update", &cc, "get:GetEditCodebasePage"),
		beego.NSRouter("/codebase/:name/update", &cc, "post:Update"),
		beego.NSRouter("/stage", &cpc, "post:DeleteCDStage"),
		beego.NSRouter("/cd-pipeline/delete", &cpc, "post:DeleteCDPipeline"),
		beego.NSRouter("/codebase/:codebaseName/branch", &cbc, "post:CreateCodebaseBranch"),

		beego.NSRouter("/library/overview", &lc, "get:GetLibraryListPage"),
		beego.NSRouter("/library/create", &lc, "get:GetCreatePage"),
		beego.NSRouter("/library", &lc, "post:Create"),

		beego.NSRouter("/diagram/overview", &dc, "get:GetDiagramPage"),
	)
	beego.AddNamespace(adminEdpNamespace)

	apiV1EdpNamespace := beego.NewNamespace(fmt.Sprintf("%s/api/v1/edp", context.BasePath),
		beego.NSRouter("/codebase", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "post:CreateCodebase"),
		beego.NSRouter("/codebase", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "get:GetCodebases"),
		beego.NSRouter("/codebase/:codebaseName", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "get:GetCodebase"),
		beego.NSRouter("/vcs", &ec, "get:GetVcsIntegrationValue"),
		beego.NSRouter("/cd-pipeline/:name", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "get:GetCDPipelineByName"),
		beego.NSRouter("/cd-pipeline/:pipelineName/stage/:stageName", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "get:GetStage"),
		beego.NSRouter("/cd-pipeline", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "post:CreateCDPipeline"),
		beego.NSRouter("/cd-pipeline/:name", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "put:UpdateCDPipeline"),
		beego.NSRouter("/codebase", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "delete:Delete"),
		beego.NSRouter("/stage", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "delete:DeleteCDStage"),
	)
	beego.AddNamespace(apiV1EdpNamespace)

	apiV1Namespace := beego.NewNamespace(fmt.Sprintf("%s/api/v1", context.BasePath),
		beego.NSRouter("/storage-class", &controllers.OpenshiftRestController{ClusterService: clusterService}, "get:GetAllStorageClasses"),
		beego.NSRouter("/repository/available", &controllers.RepositoryRestController{}, "post:IsGitRepoAvailable"),
	)
	beego.AddNamespace(apiV1Namespace)

	v2APIHandler := NewHandlerEnv(WithClient(namespacedClient), WithWorkingDir(workingDir), WithFuncMapTemplate(CreateCommonFuncMap()))
	v2APIRouter := V2APIRouter(v2APIHandler, zaplog)
	// see https://github.com/beego/beedoc/blob/master/en-US/mvc/controller/router.md#handler-register
	// and isPrefix parameter
	beego.Handler(path.Join(context.BasePath, apiV2Scope, edpScope), v2APIRouter, true)
	beego.Handler(path.Join(context.BasePath, edpScopeV2), v2APIRouter, true)
}

func V2APIRouter(h *HandlerEnv, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(WithLoggerMw(logger))

	V2RoutePrefix := path.Join(context.BasePath, apiV2Scope)
	router.Route(V2RoutePrefix, func(v2APIRouter chi.Router) {
		v2APIRouter.Route(edpScope, func(edpScope chi.Router) {
			edpScope.Route("/cd-pipeline", func(pipelinesRoute chi.Router) {
				pipelinesRoute.Route("/{pipelineName}", func(pipelineRoute chi.Router) {
					pipelineRoute.Get("/", h.GetPipeline)
					pipelineRoute.Get("/stage/{stageName}", h.GetStagePipeline)
				})
			})
			edpScope.Route("/codebase", func(codebasesRoute chi.Router) {
				codebasesRoute.Get("/", h.GetCodebases)
				codebasesRoute.Get("/{codebaseName}", h.GetCodebase)
			})
		})

	})

	edpV2Prefix := path.Join(context.BasePath, edpScopeV2)
	router.Route(edpV2Prefix, func(edpScope chi.Router) {
		edpScope.Get("/index", h.Index)
	})

	return router
}

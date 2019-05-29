/*
 * Copyright 2019 EPAM Systems.
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

package routers

import (
	"edp-admin-console/context"
	"edp-admin-console/controllers"
	"edp-admin-console/filters"
	"edp-admin-console/k8s"
	"edp-admin-console/repository"
	"edp-admin-console/service"
	"github.com/astaxie/beego"
	"log"
)

func init() {
	log.Printf("Start application in %s mode with %s EDP version...", beego.AppConfig.String("runmode"), context.EDPVersion)
	context.InitDb()
	context.InitAuth()
	clients := k8s.CreateOpenShiftClients()
	codebaseRepository := repository.CodebaseEntityRepository{}
	branchRepository := repository.ReleaseBranchRepository{}
	pipelineRepository := repository.CDPipelineRepository{}
	edpService := service.EDPTenantService{Clients: clients}
	clusterService := service.ClusterService{Clients: clients}
	branchService := service.BranchService{Clients: clients, IReleaseBranchRepository: branchRepository}
	codebaseService := service.CodebaseService{Clients: clients, ICodebaseRepository: codebaseRepository, BranchService: branchService}
	pipelineService := service.CDPipelineService{Clients: clients, ICDPipelineRepository: pipelineRepository, CodebaseService: codebaseService, BranchService: branchService}

	beego.Router("/auth/callback", &controllers.AuthController{}, "get:Callback")
	beego.InsertFilter("/admin/*", beego.BeforeRouter, filters.AuthFilter)
	beego.InsertFilter("/api/v1/edp/*", beego.BeforeRouter, filters.AuthRestFilter)
	beego.InsertFilter("/admin/edp/*", beego.BeforeRouter, filters.RoleAccessControlFilter)
	beego.InsertFilter("/api/v1/edp/*", beego.BeforeRouter, filters.RoleAccessControlRestFilter)

	beego.ErrorController(&controllers.ErrorController{})
	beego.Router("/", &controllers.MainController{EDPTenantService: edpService}, "get:Index")
	beego.Router("/cockpit-redirect-confirmation", &controllers.MainController{}, "get:CockpitRedirectConfirmationPage")

	adminEdpNamespace := beego.NewNamespace("/admin/edp",
		beego.NSRouter("/overview", &controllers.EDPTenantController{EDPTenantService: edpService}, "get:GetEDPComponents"),
		beego.NSRouter("/application/overview", &controllers.ApplicationController{CodebaseService: codebaseService, EDPTenantService: edpService, BranchService: branchService}, "get:GetApplicationsOverviewPage"),
		beego.NSRouter("/application/create", &controllers.ApplicationController{CodebaseService: codebaseService, EDPTenantService: edpService, BranchService: branchService}, "get:GetCreateApplicationPage"),
		beego.NSRouter("/application", &controllers.ApplicationController{CodebaseService: codebaseService, EDPTenantService: edpService, BranchService: branchService}, "post:CreateApplication"),

		beego.NSRouter("/cd-pipeline/overview", &controllers.CDPipelineController{CodebaseService: codebaseService, PipelineService: pipelineService, EDPTenantService: edpService, BranchService: branchService}, "get:GetContinuousDeliveryPage"),
		beego.NSRouter("/cd-pipeline/create", &controllers.CDPipelineController{CodebaseService: codebaseService, PipelineService: pipelineService, EDPTenantService: edpService, BranchService: branchService}, "get:GetCreateCDPipelinePage"),
		beego.NSRouter("/cd-pipeline", &controllers.CDPipelineController{CodebaseService: codebaseService, PipelineService: pipelineService, EDPTenantService: edpService, BranchService: branchService}, "post:CreateCDPipeline"),
		beego.NSRouter("/cd-pipeline/:pipelineName/overview", &controllers.CDPipelineController{EDPTenantService: edpService, BranchService: branchService, PipelineService: pipelineService}, "get:GetCDPipelineOverviewPage"),
		beego.NSRouter("/autotest/overview", &controllers.AutotestController{EDPTenantService: edpService, CodebaseService: codebaseService, BranchService: branchService}, "get:GetAutotestsOverviewPage"),
		beego.NSRouter("/autotest/create", &controllers.AutotestController{EDPTenantService: edpService, CodebaseService: codebaseService, BranchService: branchService}, "get:GetCreateAutotestPage"),
		beego.NSRouter("/autotest/:testName/overview", &controllers.AutotestController{EDPTenantService: edpService, CodebaseService: codebaseService, BranchService: branchService}, "get:GetAutotestOverviewPage"),
		beego.NSRouter("/autotest", &controllers.AutotestController{EDPTenantService: edpService, CodebaseService: codebaseService, BranchService: branchService}, "post:CreateAutotest"),

		beego.NSRouter("/codebase/:codebaseName/overview", &controllers.CodebaseController{CodebaseService: codebaseService, EDPTenantService: edpService, BranchService: branchService}, "get:GetCodebaseOverviewPage"),
		beego.NSRouter("/codebase/:codebaseName/branch", &controllers.BranchController{BranchService: branchService}, "post:CreateReleaseBranch"),
	)
	beego.AddNamespace(adminEdpNamespace)

	apiV1EdpNamespace := beego.NewNamespace("/api/v1/edp",
		beego.NSRouter("/codebase", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "post:CreateCodebase"),
		beego.NSRouter("/codebase", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "get:GetCodebases"),
		beego.NSRouter("/codebase/:codebaseName", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "get:GetCodebase"),
		beego.NSRouter("/vcs", &controllers.EDPTenantController{EDPTenantService: edpService}, "get:GetVcsIntegrationValue"),
		beego.NSRouter("/cd-pipeline/:name", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "get:GetCDPipelineByName"),
		beego.NSRouter("/cd-pipeline/:pipelineName/stage/:stageName", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "get:GetStage"),
		beego.NSRouter("/cd-pipeline", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "post:CreateCDPipeline"),
	)
	beego.AddNamespace(apiV1EdpNamespace)

	apiV1Namespace := beego.NewNamespace("/api/v1",
		beego.NSRouter("/storage-class", &controllers.OpenshiftRestController{ClusterService: clusterService}, "get:GetAllStorageClasses"),
		beego.NSRouter("/repository/available", &controllers.RepositoryRestController{}, "post:IsGitRepoAvailable"),
	)
	beego.AddNamespace(apiV1Namespace)
}

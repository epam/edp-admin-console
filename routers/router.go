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
	"edp-admin-console/repository/sql_builder"
	"edp-admin-console/service"
	"github.com/astaxie/beego"
	"log"
)

func init() {
	log.Printf("Start application in %s mode...", beego.AppConfig.String("runmode"))
	context.InitDb()
	context.InitAuth()
	clients := k8s.CreateOpenShiftClients()
	appQueryManager := sql_builder.ApplicationQueryBuilder{}
	branchQueryManager := sql_builder.BranchQueryBuilder{}
	appRepository := repository.ApplicationEntityRepository{QueryManager: appQueryManager}
	branchRepository := repository.ReleaseBranchRepository{QueryManager: branchQueryManager}
	pipelineRepository := repository.CDPipelineRepository{}
	edpService := service.EDPTenantService{Clients: clients}
	clusterService := service.ClusterService{Clients: clients}
	branchService := service.BranchService{Clients: clients, IReleaseBranchRepository: branchRepository}
	appService := service.ApplicationService{Clients: clients, IApplicationRepository: appRepository, BranchService: branchService}
	pipelineService := service.CDPipelineService{Clients: clients, ICDPipelineRepository: pipelineRepository}

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
		beego.NSRouter("/application/overview", &controllers.ApplicationController{AppService: appService, EDPTenantService: edpService, BranchService: branchService}, "get:GetApplicationsOverviewPage"),
		beego.NSRouter("/application/:appName/overview", &controllers.ApplicationController{AppService: appService, EDPTenantService: edpService, BranchService: branchService}, "get:GetApplicationOverviewPage"),
		beego.NSRouter("/application/create", &controllers.ApplicationController{AppService: appService, EDPTenantService: edpService, BranchService: branchService}, "get:GetCreateApplicationPage"),
		beego.NSRouter("/application", &controllers.ApplicationController{AppService: appService, EDPTenantService: edpService, BranchService: branchService}, "post:CreateApplication"),
		beego.NSRouter("/application/:appName/branch", &controllers.BranchController{BranchService: branchService}, "post:CreateReleaseBranch"),

		beego.NSRouter("/cd-pipeline/overview", &controllers.CDPipelineController{AppService: appService, PipelineService: pipelineService, EDPTenantService: edpService, BranchService: branchService}, "get:GetContinuousDeliveryPage"),
		beego.NSRouter("/cd-pipeline/create", &controllers.CDPipelineController{AppService: appService, PipelineService: pipelineService, EDPTenantService: edpService, BranchService: branchService}, "get:GetCreateCDPipelinePage"),
		beego.NSRouter("/cd-pipeline", &controllers.CDPipelineController{AppService: appService, PipelineService: pipelineService, EDPTenantService: edpService, BranchService: branchService}, "post:CreateCDPipeline"),
	)
	beego.AddNamespace(adminEdpNamespace)

	apiV1EdpNamespace := beego.NewNamespace("/api/v1/edp",
		beego.NSRouter("/application", &controllers.ApplicationRestController{AppService: appService}, "post:CreateApplication"),
		beego.NSRouter("/application", &controllers.ApplicationRestController{AppService: appService}, "get:GetApplications"),
		beego.NSRouter("/application/:appName", &controllers.ApplicationRestController{AppService: appService}, "get:GetApplication"),
		beego.NSRouter("/vcs", &controllers.EDPTenantController{EDPTenantService: edpService}, "get:GetVcsIntegrationValue"),
		beego.NSRouter("/cd-pipeline/:name", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "get:GetCDPipelineByName"),
		beego.NSRouter("/cd-pipeline/:pipelineName/stage/:stageName", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "get:GetStage"),
	)
	beego.AddNamespace(apiV1EdpNamespace)

	apiV1Namespace := beego.NewNamespace("/api/v1",
		beego.NSRouter("/storage-class", &controllers.OpenshiftRestController{ClusterService: clusterService}, "get:GetAllStorageClasses"),
		beego.NSRouter("/repository/available", &controllers.RepositoryRestController{}, "post:IsGitRepoAvailable"),
	)
	beego.AddNamespace(apiV1Namespace)
}

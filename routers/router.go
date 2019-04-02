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
	log.Printf("Start application in %s mode...", beego.AppConfig.String("runmode"))
	context.InitDb()
	context.InitAuth()
	clients := k8s.CreateOpenShiftClients()
	edpRepository := repository.EDPTenantRepository{}
	appRepository := repository.ApplicationEntityRepository{}
	edpService := service.EDPTenantService{IEDPTenantRep: edpRepository, Clients: clients}
	appService := service.ApplicationService{Clients: clients, IApplicationRepository: appRepository}
	clusterService := service.ClusterService{Clients: clients}

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
		beego.NSRouter("/application/overview", &controllers.ApplicationController{AppService: appService, EDPTenantService: edpService}, "get:GetApplicationsOverviewPage"),
		beego.NSRouter("/application/:appName/overview", &controllers.ApplicationController{AppService: appService, EDPTenantService: edpService}, "get:GetApplicationOverviewPage"),
		beego.NSRouter("/application/create", &controllers.ApplicationController{AppService: appService, EDPTenantService: edpService}, "get:GetCreateApplicationPage"),
		beego.NSRouter("/application", &controllers.ApplicationController{AppService: appService, EDPTenantService: edpService}, "post:CreateApplication"),
	)
	beego.AddNamespace(adminEdpNamespace)

	apiV1EdpNamespace := beego.NewNamespace("/api/v1/edp",
		beego.NSRouter("/application", &controllers.ApplicationRestController{AppService: appService}, "post:CreateApplication"),
		beego.NSRouter("/application", &controllers.ApplicationRestController{AppService: appService}, "get:GetApplications"),
		beego.NSRouter("/application/:appName", &controllers.ApplicationRestController{AppService: appService}, "get:GetApplication"),
		beego.NSRouter("/vcs", &controllers.EDPTenantController{EDPTenantService: edpService}, "get:GetVcsIntegrationValue"),
	)
	beego.AddNamespace(apiV1EdpNamespace)

	apiV1Namespace := beego.NewNamespace("/api/v1",
		beego.NSRouter("/storage-class", &controllers.OpenshiftRestController{ClusterService: clusterService}, "get:GetAllStorageClasses"),
		beego.NSRouter("/repository/available", &controllers.RepositoryRestController{}, "post:IsGitRepoAvailable"),
	)
	beego.AddNamespace(apiV1Namespace)
}

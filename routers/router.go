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
	/*START init context required for app*/
	log.Printf("Start application in %s mode...", beego.AppConfig.String("runmode"))
	context.InitDb()
	context.InitAuth()

	clients := k8s.CreateOpenShiftClients()
	edpRepository := repository.EDPTenantRep{}
	edpService := service.EDPTenantService{EDPTenantRep: edpRepository, Clients: clients}
	appService := service.ApplicationService{Clients: clients}
	clusterService := service.ClusterService{Clients: clients}
	/*END*/

	/*START security routing*/
	beego.Router("/auth/callback", &controllers.AuthController{}, "get:Callback")
	beego.InsertFilter("/admin/*", beego.BeforeRouter, filters.AuthFilter)
	beego.InsertFilter("/api/v1/edp/*", beego.BeforeRouter, filters.AuthRestFilter)
	beego.InsertFilter("/admin/edp/:name/*", beego.BeforeRouter, filters.RoleAccessControlFilter)
	beego.InsertFilter("/api/v1/edp/:name/*", beego.BeforeRouter, filters.RoleAccessControlFilter)
	/*END*/

	/*START general routing*/
	beego.Router("/", &controllers.MainController{}, "get:Index")
	/*END*/

	/*START edp routing*/
	beego.Router("/admin/edp/overview", &controllers.EDPTenantController{EDPTenantService: edpService}, "get:GetEDPTenants")
	beego.Router("/admin/edp/:name/overview", &controllers.EDPTenantController{EDPTenantService: edpService}, "get:GetEDPComponents")
	/*END*/

	beego.Router("/admin/edp/:name/application/overview", &controllers.MainController{}, "get:GetApplicationPage")
	beego.Router("/admin/edp/:name/application/create", &controllers.MainController{}, "get:GetCreateApplicationPage")

	/*START rest api*/
	restrictedApi := beego.NewNamespace("/api/v1/edp",
		beego.NSRouter("/:name/vcs", &controllers.EDPTenantController{EDPTenantService: edpService}, "get:GetVcsIntegrationValue"),
		beego.NSRouter("/:name/application", &controllers.AppRestController{AppService: appService}, "post:CreateApplication"),
	)
	beego.AddNamespace(restrictedApi)

	notRestrictedApi := beego.NewNamespace("/api/v1",
		beego.NSRouter("/storage-class", &controllers.OpenshiftRestController{ClusterService: clusterService}, "get:GetAllStorageClasses"),
		beego.NSRouter("/repository/available", &controllers.RepositoryRestController{}, "post:IsGitRepoAvailable"),
	)
	beego.AddNamespace(notRestrictedApi)
	/*END*/

}

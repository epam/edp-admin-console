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
	beego.InsertFilter(`/admin/edp/:name`, beego.BeforeRouter, filters.CheckEDPTenantRole)
	/*END*/

	/*START general routing*/
	beego.Router("/", &controllers.MainController{}, "get:Index")
	/*END*/

	/*START edp routing*/
	beego.Router("/admin/edp", &controllers.EDPTenantController{EDPTenantService: edpService}, "get:GetEDPTenants")
	beego.Router(`/admin/edp/:name`, &controllers.EDPTenantController{EDPTenantService: edpService}, "get:GetEDPComponents")
	/*END*/

	beego.Router("/admin/:name/application", &controllers.MainController{}, "get:GetApplicationPage")
	beego.Router("/admin/:name/application/create", &controllers.MainController{}, "get:GetCreateApplicationPage")
	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/:name/vcs", &controllers.EDPTenantController{EDPTenantService: edpService}, "get:GetVcsIntegrationValue"),
		beego.NSRouter("/storage", &controllers.ClusterController{ClusterService: clusterService}, "get:GetAllStorageClasses"),
		beego.NSRouter("/repository", &controllers.RepositoryController{}, "post:IsGitRepoAvailable"),
		beego.NSNamespace("/:name/application",
			beego.NSRouter("/create", &controllers.AppController{AppService: appService}, "post:CreateApplication"),
		),
	)
	beego.AddNamespace(ns)

}
